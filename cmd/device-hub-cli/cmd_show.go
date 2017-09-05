// Copyright © 2017 thingful

package main

import (
	"context"
	"strings"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

var showCommand = &cobra.Command{
	Use:   "show",
	Short: "Display one or many resources",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := roundTrip(func(cli proto.HubClient, in rawContent, out iocodec.Encoder) error {
			req := proto.ShowRequest{
				Filter: strings.Join(args, ","),
			}
			resp, err := cli.Show(context.Background(), &req)
			if err != nil {
				return err
			}
			return out.Encode(resp)
		})

		return err
	},
}
