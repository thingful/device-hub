// Copyright © 2017 thingful

package main

import (
	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete listener, profile and endpoint resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		var uri string
		if len(args) > 0 {
			uri = args[0]
		}

		_resources.Reverse()

		for _, r := range _resources.R {
			err := r.SendDelete(uri)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
