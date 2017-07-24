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

func startCommand() *cobra.Command {

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start processing messages on a uri",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, client, err := dial()
			if err != nil {
				return err
			}
			defer conn.Close()

			if len(args) == 0 {
				return errors.New("specify a profile")
			}

			err = startCall(processConf{ProfileUID: args[0]}, client)
			if err != nil {
				return err
			}
			return nil
		},
	}

	startCommand.Flags().StringVarP(&_config.ProcessFile.ListenerUID, "listener", "l", _config.ProcessFile.ListenerUID, "listener uid to accept messages on")
	startCommand.Flags().StringVarP(&_config.ProcessFile.URI, "uri", "u", _config.ProcessFile.URI, "uri to listen on")
	startCommand.Flags().StringSliceVarP(&_config.ProcessFile.EndpointUIDs, "endpoint", "e", _config.ProcessFile.EndpointUIDs, "endpoint uid to push messages to, may be specified multiple times")
	startCommand.Flags().StringSliceVarP(&_config.ProcessFile.Tags, "tags", "t", _config.ProcessFile.Tags, "colon separated (k:v) runtime tags to attach to requests, may be specified multiple times")

	return startCommand
}

func startCall(conf processConf, client proto.HubClient) error {
	req := proto.StartRequest{
		Tags: map[string]string{},
	}

	req.Profile = conf.ProfileUID
	if req.Profile == "" {
		req.Profile = _config.ProcessFile.ProfileUID
	}
	req.Listener = conf.ListenerUID
	if req.Listener == "" {
		req.Listener = _config.ProcessFile.ListenerUID
	}
	req.Uri = conf.URI
	if req.Uri == "" {
		req.Uri = _config.ProcessFile.URI
	}
	req.Endpoints = conf.EndpointUIDs
	if len(req.Endpoints) == 0 {
		req.Endpoints = _config.ProcessFile.EndpointUIDs
	}

	var tags []string

	if len(conf.Tags) != 0 {
		tags = conf.Tags
	} else {
		tags = _config.ProcessFile.Tags
	}

	for _, m := range tags {
		bits := strings.Split(m, ":")
		if len(bits) != 2 {
			return fmt.Errorf("metadata not colon (:) separated : %s", m)
		}
		req.Tags[bits[0]] = bits[1]
	}
	fmt.Println(req)
	resp, err := client.Start(context.Background(), &req)
	if err != nil {
		return err
	}
	_encoder.Encode(resp)
	return nil
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

func stopCall(uri string, client proto.HubClient) error {
	req := proto.StopRequest{
		Uri: strings.TrimSpace(uri),
	}
	resp, err := client.Stop(context.Background(), &req)
	if err != nil {
		return err
	}
	_encoder.Encode(resp)
	return nil
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
