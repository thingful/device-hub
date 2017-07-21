// Copyright © 2017 thingful

package main

import (
	"github.com/spf13/cobra"
)

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create listener, endpoint and profile resources",
	RunE: func(cmd *cobra.Command, args []string) error {

		for _, r := range _resources.R {
			err := r.SendCreate()

			if err != nil {
				return err
			}

		}

		return nil
	},
}
