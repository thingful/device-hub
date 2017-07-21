// Copyright © 2017 thingful

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/thingful/device-hub/proto"

	yaml "gopkg.in/yaml.v2"
)

type rawContent []byte

func (r rawContent) Decode(target interface{}) error {
	err := yaml.Unmarshal(r, target)
	if err != nil {
		return fmt.Errorf("error decoding data: %s", err.Error())
	}
	return nil
}

// Represent a resource file, Data is basically used to order a resourceSlice
// Raw contains the file content
type resource struct {
	FileName string
	Data     map[string]interface{}
	Raw      rawContent
}

// Load the configuration file to Data
func (r *resource) Load(filePath string) (err error) {
	r.Raw, err = ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file [%s]: %s", filePath, err.Error())
	}

	_, r.FileName = filepath.Split(filePath)

	r.Raw.Decode(&r.Data)
	if err != nil {
		return fmt.Errorf("error parsing file [%s]: %s", filePath, err.Error())
	}
	return nil
}

func (r *resource) sendCreateReq() error {
	err := roundTrip(func(client proto.HubClient, in rawContent, out iocodec.Encoder) error {
		req := proto.CreateRequest{}
		err := r.Raw.Decode(&req)
		if err != nil {
			return err
		}

		resp, err := client.Create(context.Background(), &req)
		if err != nil {
			return err
		}
		return out.Encode(resp)
	})
	return err
}

func (r *resource) sendStartReq(uri string) error {
	err := roundTrip(func(client proto.HubClient, in rawContent, out iocodec.Encoder) error {
		req := proto.StartRequest{
			Endpoints: []string{},
			Tags:      map[string]string{},
		}

		if _config.RequestFile == "" {
			err := r.Raw.Decode(&_config.ProcessConf)
			if err != nil {
				return err
			}
		} else {
			r.Raw.Decode(&_config.ProcessConf)
		}

		if len(uri) > 0 {
			req.Profile = uri
		} else {
			req.Profile = _config.ProcessConf.ProfileUID
		}

		req.Uri = _config.ProcessConf.URI
		req.Listener = _config.ProcessConf.ListenerUID
		req.Endpoints = _config.ProcessConf.EndpointUIDs

		// review tags
		for _, m := range _config.ProcessConf.Tags {
			bits := strings.Split(m, ":")
			if len(bits) != 2 {
				return fmt.Errorf("metadata not colon (:) separated : %s", m)
			}
			req.Tags[bits[0]] = bits[1]
		}

		resp, err := client.Start(context.Background(), &req)
		if err != nil {
			return err
		}
		return out.Encode(resp)
	})
	return err
}

func (r *resource) SendCreate(args ...string) error {
	var uri string
	if r.Data["type"] == "process" {
		if len(args) > 0 {
			uri = args[0]
		}
		return r.sendStartReq(uri)
	}
	return r.sendCreateReq()
}

func (r *resource) sendStopReq(uri string) error {
	err := roundTrip(func(client proto.HubClient, in rawContent, out iocodec.Encoder) error {
		if uri == "" {
			uri = r.Data["uri"].(string)
		}
		return stopCall(uri, client, out)
	})
	return err
}

func (r *resource) sendDeleteReq() error {
	err := roundTrip(func(client proto.HubClient, in rawContent, out iocodec.Encoder) error {
		req := proto.DeleteRequest{}
		err := r.Raw.Decode(&req)
		if err != nil {
			return err
		}

		resp, err := client.Delete(context.Background(), &req)
		if err != nil {
			return err
		}
		return out.Encode(resp)
	})
	return err
}

func (r *resource) SendDelete(uri string) error {
	if r.Data["type"] == "process" {
		return r.sendStopReq(uri)
	}
	return r.sendDeleteReq()
}

// resourceSlice contains configs and implements sort interface using
// "process" type file as the less weight
type resourceSlice struct {
	R []resource
}

func newResources() *resourceSlice {
	return &resourceSlice{}
}

func (r *resourceSlice) Append(e resource) {

	r.R = append(r.R, e)
}

func (r resourceSlice) Len() int {
	return len(r.R)
}

func (r resourceSlice) Less(i, j int) bool {
	if r.R[j].Data["type"] == "process" {
		return true
	}
	return false
}

func (r resourceSlice) Swap(i, j int) {
	r.R[i], r.R[j] = r.R[j], r.R[i]
}

// Sort is required to put processes at the end when executing create cmd
func (r resourceSlice) Sort() {
	sort.Sort(resourceSlice(r))
}

// Reverse is required to put processes at first to stop them before delete resources
func (r resourceSlice) Reverse() {
	sort.Sort(sort.Reverse(resourceSlice(r)))
}

func (r resourceSlice) Print() {
	for _, f := range r.R {
		fmt.Println(f.FileName)
	}
}

// GetCliConfig get config for the CLI app
func (r *resourceSlice) SetResources(cfg *config) error {
	var res resource
	if cfg.RequestFile != "" {
		err := res.Load(cfg.RequestFile)
		if err != nil {
			return err
		}
		r.Append(res)
		return nil

	} else if cfg.RequestDir != "" {
		listing, err := ioutil.ReadDir(cfg.RequestDir)
		if err != nil {
			return err
		}

		for _, f := range listing {
			folderPath := path.Join(cfg.RequestDir, f.Name())
			var _res resource // check if this could be outer scoped! (conf)
			err = _res.Load(folderPath)
			if err != nil {
				return err
			}
			r.Append(_res)
		}
	}
	// Sorted with process to the end by default
	r.Sort()
	// TODO validate?
	r.Print()
	return nil
}

// could be proto.StartRequest?
type processConf struct {
	URI          string   `yaml:"uri"`
	Type         string   `yaml:"type"`
	EndpointUIDs []string `yaml:"endpoint-uids"`
	ListenerUID  string   `yaml:"listener-uid"`
	ProfileUID   string   `yaml:"profile-uid"`
	Tags         []string `yaml:"tags"`
}
