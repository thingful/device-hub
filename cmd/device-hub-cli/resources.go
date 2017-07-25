// Copyright © 2017 thingful

package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"

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

// resource Represent a resource file, Data is basically used to order a resourceSlice
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

// resources store appended files
type resources struct {
	R []resource
}

func newResources() *resources {
	return &resources{}
}

func (r *resources) Append(e resource) {

	r.R = append(r.R, e)
}

func (r resources) Len() int {
	return len(r.R)
}

func (r resources) Less(i, j int) bool {
	if r.R[j].Data["type"] == "process" {
		return true
	}
	return false
}

func (r resources) Swap(i, j int) {
	r.R[i], r.R[j] = r.R[j], r.R[i]
}

// Sort is required to put processes at the end when executing create cmd
func (r resources) Sort() {
	sort.Sort(resources(r))
}

// Reverse is required to put processes at first to stop them before delete resources
func (r resources) Reverse() {
	sort.Sort(sort.Reverse(resources(r)))
}

func (r resources) Print() {
	for _, f := range r.R {
		fmt.Println(f.FileName)
	}
}

// SetResources get passed resources
func (r *resources) SetResources(cfg *config) error {
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

type processConf struct {
	URI          string   `yaml:"uri"`
	Type         string   `yaml:"type"`
	EndpointUIDs []string `yaml:"endpoint-uids"`
	ListenerUID  string   `yaml:"listener-uid"`
	ProfileUID   string   `yaml:"profile-uid"`
	Tags         []string `yaml:"tags"`
}
