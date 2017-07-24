// Copyright © 2017 thingful

package main

import (
	"context"
	"errors"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

func startCommand() *cobra.Command {

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start processing messages on a uri",
		RunE: func(cmd *cobra.Command, args []string) error {
			var profile string

			if len(args) > 0 {
				profile = args[0]
			}
			if len(_resources.R) == 0 {
				_resources.R = append(_resources.R,
					resource{
						Data: map[string]interface{}{"type": "process"}})
			}
			err := _resources.R[0].SendCreate(profile)
			if err != nil {
				return err
			}

			return nil
		},
	}

	startCommand.Flags().StringVarP(&_config.ProcessConf.ListenerUID, "listener", "l", _config.ProcessConf.ListenerUID, "listener uid to accept messages on")
	startCommand.Flags().StringVarP(&_config.ProcessConf.URI, "uri", "u", _config.ProcessConf.URI, "uri to listen on")
	startCommand.Flags().StringSliceVarP(&_config.ProcessConf.EndpointUIDs, "endpoint", "e", _config.ProcessConf.EndpointUIDs, "endpoint uid to push messages to, may be specified multiple times")
	startCommand.Flags().StringSliceVarP(&_config.ProcessConf.Tags, "tags", "t", _config.ProcessConf.Tags, "colon separated (k:v) runtime tags to attach to requests, may be specified multiple times")

	return startCommand
}

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop processing messages on a uri",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("specify a uri to stop")
		}

		err := roundTrip(func(client proto.HubClient, in rawContent, out iocodec.Encoder) error {
			return stopCall(args[0], client, out)
		})

		return err
	},
}

func stopCall(uri string, client proto.HubClient, out iocodec.Encoder) error {
	req := proto.StopRequest{
		Uri: uri,
	}

	resp, err := client.Stop(context.Background(), &req)
	if err != nil {
		return err
	}
	return out.Encode(resp)

}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "List running pipes",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := roundTrip(func(client proto.HubClient, in rawContent, out iocodec.Encoder) error {
			req := proto.StatusRequest{}
			resp, err := client.Status(context.Background(), &req)

			if err != nil {
				return err
			}
			return out.Encode(resp)
		})
		return err
	},
}
