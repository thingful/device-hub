// Copyright © 2017 thingful

package main

import (
	"context"
	"errors"
	"strings"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

func startCommand() *cobra.Command {

	request := proto.StartRequest{
		Endpoints: []string{},
	}

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start processing messages on a uri",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				return errors.New("specify a profile")
			}

			request.Profile = args[0]

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
	startCommand.Flags().StringVarP(&request.Listener, "listener", "l", request.Listener, "listener to use")
	startCommand.Flags().StringVarP(&request.Uri, "uri", "u", request.Uri, "uri to listen on")
	startCommand.Flags().StringSliceVarP(&request.Endpoints, "endpoint", "e", request.Endpoints, "endpoint to use")

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

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List running pipes",
	RunE: func(cmd *cobra.Command, args []string) error {

		v := proto.ListRequest{}

		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

			resp, err := cli.List(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		return err
	},
}
