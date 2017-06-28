// Copyright © 2017 thingful

package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/endpoint"
	"github.com/thingful/device-hub/listener"
	"github.com/thingful/device-hub/listener/grovepi"
	"github.com/thingful/device-hub/registry"
)

var describeCommand = &cobra.Command{
	Use:   "describe",
	Short: "Describe parameters for endpoint and listeners",
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 2 {
			return errors.New("describe [listener|endpoint] kind")
		}

		typez := strings.ToLower(args[0])
		kind := strings.ToLower(args[1])

		var params describe.Parameters
		var err error

		register := registry.Default

		endpoint.Register(register)
		listener.Register(register)
		grovepi.Register(register)

		switch typez {
		case "listener":
			params, err = register.DescribeListener(kind)

		case "endpoint":
			params, err = register.DescribeEndpoint(kind)

		default:
			return errors.New("describe [listener|endpoint] kind")
		}

		if err != nil {
			return err
		}
		for _, param := range params {
			fmt.Println(param.Describe())
		}

		return nil

	},
}
