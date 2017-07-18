// Copyright © 2017 thingful

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

type clientConfig struct {
	URI          string   `yaml:"uri"`
	Type         string   `yaml:"type"`
	EndpointUIDs []string `yaml:"endpoint-uids"`
	ListenerUID  string   `yaml:"listener-uid"`
	ProfileUID   string   `yaml:"profile-uid"`
	Tags         []string `yaml:"tags"`
}

func startCommand() *cobra.Command {

	request := proto.StartRequest{
		Endpoints: []string{},
		Tags:      map[string]string{},
	}

	tags := []string{}

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start processing messages on a uri",
		RunE: func(cmd *cobra.Command, args []string) error {
			// if no profile is provided as arg, then the profile field in the yaml file
			// will be loaded.
			if len(args) > 0 {
				request.Profile = args[0]
			}
			err := roundTrip(request, "start", func(cli proto.HubClient, in rawConf, out iocodec.Encoder) error {

				err := startCall(args, request, tags, cli, in, out)
				if err != nil {
					return err
				}
				return nil
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

func startCall(args []string, request proto.StartRequest, tags []string, cli proto.HubClient, in rawConf, out iocodec.Encoder) error {

	cfg := clientConfig{}

	if _config.RequestFile == "" {
		err := in.Decode(&cfg)
		if err != nil {
			return err
		}
	} else {
		err := yamlDecoder(_config.RequestFile, &cfg)
		if err != nil {
			return err
		}
	}

	// No sure about the name yet (process, pipe, etc.)
	if cfg.Type != "process" {
		return fmt.Errorf("file doesn't have the needed type [%v]", cfg.Type)
	}

	if request.Profile == "" {
		request.Profile = cfg.ProfileUID
	}

	request.Uri = cfg.URI
	request.Listener = cfg.ListenerUID
	request.Endpoints = cfg.EndpointUIDs

	if request.Profile == "" {
		return errors.New("no profile specified")
	}
	// review tags
	tags = cfg.Tags
	for _, m := range tags {

		bits := strings.Split(m, ":")

		if len(bits) != 2 {
			return fmt.Errorf("metadata not colon (:) separated : %s", m)
		}

		request.Tags[bits[0]] = bits[1]
	}

	resp, err := cli.Start(context.Background(), &request)

	if err != nil {
		return err
	}

	return out.Encode(resp)
}

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop processing messages on a uri",
	RunE: func(cmd *cobra.Command, args []string) error {

		v := proto.StopRequest{}

		err := roundTrip(v, "stop", func(cli proto.HubClient, in rawConf, out iocodec.Encoder) error {
			return stopCall(args, v, cli, in, out)
		})
		return err
	},
}

func stopCall(args []string, request proto.StopRequest, cli proto.HubClient, in rawConf, out iocodec.Encoder) error {

	err := in.Decode(&request)
	if err != nil {
		return err
	}
	if len(args) == 0 && request.Uri == "" {
		return errors.New("specify a uri to stop")
	}
	if len(args) > 0 {
		request.Uri = strings.TrimSpace(args[0])
	}

	resp, err := cli.Stop(context.Background(), &request)

	if err != nil {
		return err
	}

	return out.Encode(resp)
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "List running pipes",
	RunE: func(cmd *cobra.Command, args []string) error {

		v := proto.StatusRequest{}

		err := roundTrip(v, "status", func(cli proto.HubClient, in rawConf, out iocodec.Encoder) error {

			resp, err := cli.Status(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		return err
	},
}
