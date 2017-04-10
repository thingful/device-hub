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
	config *config.Configuration
}

func NewDeviceHub(config *config.Configuration) app {

	return app{
		config: config,
	}
}

// Run attempts to start up the configuration by iterating thought the
// pipes section ensuring that the endpoints and listeners exist.
// Will return an error if the complete configuration is unable to start correctly.
func (a app) Run() (context.Context, error) {

	ctx := context.Background()

	manager, err := server.NewEndpointManager(ctx, a.config)

	if err != nil {
		return nil, err
	}

	err = manager.Start()

	if err != nil {
		return nil, err
	}

	// launch RPC server
	// TODO : configure this properly
	go server.Serve(manager)

	return ctx, nil
}
