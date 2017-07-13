// Copyright © 2017 thingful

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"io/ioutil"
	"log"

	"os"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
	"gopkg.in/yaml.v2"
)

func startCommand() *cobra.Command {

	var configPath string
	var configFile bool
	type clientConfig struct {
		URI          string   `yaml:"uri"`
		Type         string   `yaml:"type"`
		EndpointUIDs string   `yaml:"endpoint-uid"`
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

	startCommand.Flags().BoolVarP(&configFile, "config-file", "c", configFile, "enable config file feature")
	startCommand.Flags().StringVar(&configPath, "config-path", configPath, "config file path with the required resources")
	startCommand.ParseFlags(os.Args)

	// fmt.Println("configFile:", configFile)
	// TODO check type
	if configFile {

		readFields := clientConfig{}

		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Failed to read config file [%s]: %s\n", configPath, err.Error())
		}
		err = yaml.Unmarshal(content, &readFields)
		if err != nil {
			log.Fatalf("Error parsing config file [%s]: %s\n", configPath, err.Error())
		}
		fmt.Println(readFields)
		request.Listener = readFields.ListenerUID
		request.Uri = readFields.URI
		// TODO Parse string slice
		//request.Endpoints = readFields.EndpointUIDs
		tags = readFields.Tags

	} else {
		startCommand.Flags().StringVarP(&request.Listener, "listener", "l", request.Listener, "listener uid to accept messages on")
		startCommand.Flags().StringVarP(&request.Uri, "uri", "u", request.Uri, "uri to listen on")
		startCommand.Flags().StringSliceVarP(&request.Endpoints, "endpoint", "e", request.Endpoints, "endpoint uid to push messages to, may be specified multiple times")
		startCommand.Flags().StringSliceVarP(&tags, "tags", "t", tags, "colon separated (k:v) runtime tags to attach to requests, may be specified multiple times")
	}
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
