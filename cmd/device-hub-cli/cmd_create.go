// Copyright © 2017 thingful

package main

import (
	"context"
	"strings"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/endpoint"
	"github.com/thingful/device-hub/listener"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/registry"
)

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create listener, endpoint and profile resources",
	RunE: func(cmd *cobra.Command, args []string) error {

		/* TODO : add ability to generate examples */
		sample := proto.CreateRequest{
			Configuration: map[string]string{},
		}

		err := roundTrip(sample, func(cli proto.HubClient, in rawConf, out iocodec.Encoder) error {

			v := proto.CreateRequest{}

			err := in.Decode(&v)

			if err != nil {
				return err
			}

			// validate the policy file before sending it over the wire
			var params describe.Parameters

			register := registry.Default

			endpoint.Register(register)
			listener.Register(register)

			switch strings.ToLower(v.Type) {

			case "listener":
				params, err = register.DescribeListener(v.Kind)
			case "endpoint":
				params, err = register.DescribeEndpoint(v.Kind)
			case "process":
				request := proto.StartRequest{
					Endpoints: []string{},
					Tags:      map[string]string{},
				}
				tags := []string{}
				return startCall(args, request, tags, cli, in, out)
			}

			if err != nil {
				return err
			}

			_, err = describe.NewValues(v.Configuration, params)

			if err != nil {
				return err
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
