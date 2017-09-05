// Copyright © 2017 thingful

package main

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

var showCommand = &cobra.Command{
	Use:   "show",
	Short: "Display one or many resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := proto.ShowRequest{
			Filter: strings.Join(args, ","),
		}

		conn, client, err := dial()
		if err != nil {
			return err
		}
		defer conn.Close()

		resp, err := client.Show(context.Background(), &req)
		if err != nil {
			return err
		}

		return _encoder.Encode(resp)
	},
}
