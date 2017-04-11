// Copyright © 2017 thingful

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/config"
)

var RootCmd = &cobra.Command{
	Use: "device-hub",
}

func init() {

	var configurationPath string

	RootCmd.PersistentFlags().StringVarP(&configurationPath, "config", "c", "./config.json", "Path to configuration file.")

	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(hub.DaemonVersionString())
			return nil
		},
	}

	RootCmd.AddCommand(versionCommand)

	checkConfig := &cobra.Command{
		Use:   "check",
		Short: "Load,parse and check the configuration file.",
		RunE: func(cmd *cobra.Command, args []string) error {

			_, err := config.LoadFromFile(configurationPath)
			return err
		},
	}

	RootCmd.AddCommand(checkConfig)

	daemon := &cobra.Command{
		Use:   "start",
		Short: "Start device hub.",
		RunE: func(cmd *cobra.Command, args []string) error {

			conf, err := config.LoadFromFile(configurationPath)

			if err != nil {
				return err
			}
			app := NewDeviceHub(conf)

			ctx, err := app.Run()

			if err != nil {
				return err
			}

			<-ctx.Done()

			return nil
		},
	}
	RootCmd.AddCommand(daemon)

}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
