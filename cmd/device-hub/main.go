// Copyright © 2017 thingful

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/config"
	"github.com/thingful/device-hub/server"
)

var RootCmd = &cobra.Command{
	Use: "device-hub",
}

func init() {

	var configurationPath string

	// Client can run either in insecure mode or provide details for mutual tls
	// The default is for secure connections to be used.
	var options server.Options

	RootCmd.PersistentFlags().StringVarP(&options.Binding, "binding", "b", ":50051", "RPC binding for the device-hub daemon.")
	RootCmd.PersistentFlags().BoolVar(&options.Insecure, "insecure", false, "Switch off Mutual TLS authentication.")
	RootCmd.PersistentFlags().StringVar(&options.CertFilePath, "cert-file", "", "Certificate used for SSL/TLS RPC connections to the device-hub daemon.")
	RootCmd.PersistentFlags().StringVar(&options.KeyFilePath, "key-file", "", "Key file for the certificate (--cert-file).")
	RootCmd.PersistentFlags().StringVar(&options.TrustedCAFilePath, "trusted-ca-file", "", "Trusted certificate authority.")

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

			app := NewDeviceHub(options, conf)

			ctx := context.Background()

			err = app.Run(ctx)

			if err != nil {
				return err
			}

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
