// Copyright © 2017 thingful

package main

import (
	"context"
	"log"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

var listenerCommand = &cobra.Command{
	Use:   "listener",
	Short: "Add, Delete and List listeners.",
}

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create listeners and endpoints",
	Run: func(cmd *cobra.Command, args []string) {

		/* TODO : add ability to generate examplesi */
		v := proto.CreateRequest{
			Type: "listener",
			Kind: "http",
			Configuration: map[string]string{
				"HTTPAddress": "0.0.0.0:8085",
			},
		}

		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

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
		if err != nil {
			log.Fatal(err)
		}
	},
}
