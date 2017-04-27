package main

import (
	"context"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete listener and endpoint resources",
	RunE: func(cmd *cobra.Command, args []string) error {

		v := proto.DeleteRequest{}

		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

			err := in.Decode(&v)
			if err != nil {
				return err
			}

			resp, err := cli.Delete(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})

		return err

	},
}
