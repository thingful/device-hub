// Copyright © 2017 thingful

package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/proto"
)

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete listener, profile and endpoint resources",
	RunE: func(cmd *cobra.Command, args []string) error {

		conn, client, err := dial()
		if err != nil {
			return err
		}
		defer conn.Close()

		_resources.Reverse()
		for _, r := range _resources.R {
			req := proto.DeleteRequest{}
			err := r.Raw.Decode(&req)
			if err != nil {
				return err
			}
			if r.Data["type"] != "process" {
				resp, err := client.Delete(context.Background(), &req)
				if err != nil {
					return err
				}
				return _encoder.Encode(resp)
			} else {
				var conf processConf
				err := r.Raw.Decode(&conf)
				if err != nil {
					return err
				}
				stopCall(conf.URI, client)
			}

		}

		return nil

	},
}
