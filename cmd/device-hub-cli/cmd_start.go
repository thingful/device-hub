// Copyright © 2017 thingful

package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

func startCommand() *cobra.Command {

	type clientConfig struct {
		URI          string   `yaml:"uri"`
		Type         string   `yaml:"type"`
		EndpointUIDs []string `yaml:"endpoint-uids"`
		ListenerUID  string   `yaml:"listener-uid"`
		Tags         []string `yaml:"tags"`
	}

	request := proto.StartRequest{
		Endpoints: []string{},
		Tags:      map[string]string{},
	}

	tags := []string{}

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start processing messages on a uri",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				return errors.New("specify a profile")
			}

			request.Profile = args[0]

			if _config.RequestFile != "" {

				readFields := clientConfig{}

				content, err := ioutil.ReadFile(_config.RequestFile)
				if err != nil {
					return fmt.Errorf("failed to read file [%s]: %s", _config.RequestFile, err.Error())
				}
				err = yaml.Unmarshal(content, &readFields)
				if err != nil {
					return fmt.Errorf("error parsing file [%s]: %s", _config.RequestFile, err.Error())
				}
				// No sure about the name yet (process, pipe, etc.)
				if readFields.Type != "process" {
					return fmt.Errorf("file doesn't have the needed type")
				}
				request.Listener = readFields.ListenerUID
				request.Uri = readFields.URI
				request.Endpoints = readFields.EndpointUIDs
				tags = readFields.Tags
			}
			for _, m := range tags {

				bits := strings.Split(m, ":")

				if len(bits) != 2 {
					return fmt.Errorf("metadata not colon (:) separated : %s", m)
				}

				request.Tags[bits[0]] = bits[1]
			}

			err := roundTrip(request, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

				resp, err := cli.Start(context.Background(), &request)

				if err != nil {
					return err
				}

				return out.Encode(resp)

			})
			return err
		},
	}

	startCommand.Flags().StringVarP(&request.Listener, "listener", "l", request.Listener, "listener uid to accept messages on")
	startCommand.Flags().StringVarP(&request.Uri, "uri", "u", request.Uri, "uri to listen on")
	startCommand.Flags().StringSliceVarP(&request.Endpoints, "endpoint", "e", request.Endpoints, "endpoint uid to push messages to, may be specified multiple times")
	startCommand.Flags().StringSliceVarP(&tags, "tags", "t", tags, "colon separated (k:v) runtime tags to attach to requests, may be specified multiple times")

	return startCommand
}

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop processing messages on a uri",
	RunE: func(cmd *cobra.Command, args []string) error {

		v := proto.StopRequest{}

		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {
			if len(args) == 0 {
				return errors.New("specify a uri to stop")
			}

			v.Uri = strings.TrimSpace(args[0])

			resp, err := cli.Stop(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		return err
	},
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "List running pipes",
	RunE: func(cmd *cobra.Command, args []string) error {

		v := proto.StatusRequest{}

		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

			resp, err := cli.Status(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		return err
	},
}
