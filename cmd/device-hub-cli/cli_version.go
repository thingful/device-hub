// Copyright © 2017 thingful

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	hub "github.com/thingful/device-hub"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Display version information.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(hub.ClientVersionString())
		return nil
	},
}
