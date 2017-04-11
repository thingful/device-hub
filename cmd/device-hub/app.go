// Copyright © 2017 thingful

package main

import (
	"context"

	"github.com/thingful/device-hub/config"
	_ "github.com/thingful/device-hub/endpoint"
	_ "github.com/thingful/device-hub/listener"
	"github.com/thingful/device-hub/server"
)

type app struct {
	config  *config.Configuration
	options server.Options
}

func NewDeviceHub(options server.Options, config *config.Configuration) app {

	return app{
		config:  config,
		options: options,
	}
}

// Run attempts to start up the configuration by iterating thought the
// pipes section ensuring that the endpoints and listeners exist.
// Will return an error if the complete configuration is unable to start correctly.
func (a app) Run(ctx context.Context) error {

	manager, err := server.NewEndpointManager(ctx, a.config)

	if err != nil {
		return err
	}

	err = manager.Start()

	if err != nil {
		return err
	}

	return server.Serve(a.options, manager)
}
