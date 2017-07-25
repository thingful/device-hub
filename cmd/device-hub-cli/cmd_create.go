// Copyright © 2017 thingful

package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create listener, endpoint and profile resources",
	RunE: func(cmd *cobra.Command, args []string) error {

		conn, client, err := dial()
		if err != nil {
			return err
		}
		defer conn.Close()

		for _, r := range _resources.R {
			req := proto.CreateRequest{}
			err := r.Raw.Decode(&req)
			if err != nil {
				return err
			}
			if r.Data["type"] != "process" {
				resp, err := client.Create(context.Background(), &req)
				if err != nil {
					return err
				}
				_encoder.Encode(resp)
			} else {
				var conf processConf
				err := r.Raw.Decode(&conf)
				if err != nil {
					return err
				}
				startCall(conf, client)
			}

		}

		return nil
	},
}
