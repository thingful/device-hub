// Copyright © 2017 thingful

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

func startCommand() *cobra.Command {

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start processing messages on a uri",
		RunE: func(cmd *cobra.Command, args []string) error {
			var profile string
			conn, client, err := dial()
			if err != nil {
				return err
			}
			defer conn.Close()

			if len(args) != 0 {
				profile = strings.TrimSpace(args[0])
			}
			if _config.RequestFile != "" {
				r := resource{}
				err = r.Load(_config.RequestFile)
				if err != nil {
					return err
				}
				r.Raw.Decode(&_config.ProcessFile)
			}
			err = startCall(processConf{ProfileUID: profile}, client)
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
		var uri string

		if len(args) != 0 {
			uri = strings.TrimSpace(args[0])
		}

		conn, client, err := dial()
		if err != nil {
			return err
		}
		defer conn.Close()

		if _config.RequestFile != "" {
			r := resource{}
			err = r.Load(_config.RequestFile)
			if err != nil {
				return err
			}
			r.Raw.Decode(&_config.ProcessFile)
			uri = _config.ProcessFile.URI
		}
		err = stopCall(uri, client)
		if err != nil {
			return err
		}
		return nil
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
	return _encoder.Encode(resp)
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "List running pipes",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := proto.StatusRequest{}

		conn, client, err := dial()
		if err != nil {
			return err
		}
		defer conn.Close()

		resp, err := client.Status(context.Background(), &req)

		if err != nil {
			return err
		}

		return _encoder.Encode(resp)
	},
}
