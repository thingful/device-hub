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

		v := proto.ShowRequest{
			Filter: strings.Join(args, ","),
		}

		err := roundTrip(v, "show", func(cli proto.HubClient, in rawConf, out iocodec.Encoder) error {
			resp, err := cli.Show(context.Background(), &v)
			if err != nil {
				return err
			}

			return out.Encode(resp)

		})

		return err
	},
}
