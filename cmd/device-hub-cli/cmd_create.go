// Copyright © 2017 thingful

package main

import (
	"context"
	"strings"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/proto"
)

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create listener, endpoint and profile resources",
	RunE: func(cmd *cobra.Command, args []string) error {

		/* TODO : add ability to generate examples */
		sample := proto.CreateRequest{
			Configuration: map[string]string{},
		}

		err := roundTrip(sample, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

			v := proto.CreateRequest{}

			err := in.Decode(&v)

			if err != nil {
				return err
			}
			switch strings.ToLower(v.Type) {

			case "listener":
				params, err := hub.DescribeListener(v.Kind)

				if err != nil {
					return err
				}
				_, err = describe.NewValues(v.Configuration, params)

				if err != nil {
					return err
				}

			case "endpoint":

				params, err := hub.DescribeEndpoint(v.Kind)

				if err != nil {
					return err
				}
				_, err = describe.NewValues(v.Configuration, params)

				if err != nil {
					return err
				}
			}

			resp, err := cli.Create(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		return err
	},
}
