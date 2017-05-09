// Copyright © 2017 thingful

package main

import (
	"context"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
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

			resp, err := cli.Create(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		return err
	},
}
