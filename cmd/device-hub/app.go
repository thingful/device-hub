// Copyright © 2017 thingful

package main

import (
	"context"
	"fmt"
	"log"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/config"
	_ "github.com/thingful/device-hub/endpoint"
	"github.com/thingful/device-hub/engine"
	_ "github.com/thingful/device-hub/listener"
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
// pipes section ensuring that the endpoints and listneners exist.
// Will return an error if the complete configuration is unable to start correctly.
func (a app) Run() (context.Context, error) {

	ctx := context.Background()

	for _, pipe := range a.config.Pipes {

		found, listenerConf := a.config.Listeners.FindByUID(pipe.Listener)

		if !found {
			return nil, fmt.Errorf("listener with UID %s not found", pipe.Listener)
		}

		endpoints := []hub.Endpoint{}

		for e := range pipe.Endpoints {

			found, endpointConf := a.config.Endpoints.FindByUID(pipe.Endpoints[e])

			if !found {
				return nil, fmt.Errorf("endpoint with UID %s not found", pipe.Endpoints[e])
			}

			newendpoint, err := hub.EndpointByName(string(endpointConf.UID), endpointConf.Type, endpointConf.Configuration)

			if err != nil {
				return nil, err
			}

			endpoints = append(endpoints, newendpoint)

		}

		found, profile := a.config.Profiles.FindByUID(pipe.Profile)

		if !found {
			return nil, fmt.Errorf("profile with UID %s not found", pipe.Profile)
		}

		listener, err := hub.ListenerByName(string(listenerConf.UID), listenerConf.Type, listenerConf.Configuration)

		if err != nil {
			return nil, err
		}

		channel, err := listener.NewChannel(pipe.Uri)

		if err != nil {
			return nil, err
		}

		go startPipe(ctx, listener, channel, profile, endpoints)

	}

	return ctx, nil
}

func startPipe(ctx context.Context, listener hub.Listener, channel hub.Channel, profile config.Profile, endpoints []hub.Endpoint) {

	scripter := engine.New()

	for {

		select {

		case <-ctx.Done():
			return

		case err := <-channel.Errors():
			log.Println(err)

		case input := <-channel.Out():

			output, err := scripter.Execute(profile.Script, input)

			if err != nil {
				log.Println(err)
			}

			output.Metadata[hub.PROFILE_NAME_KEY] = profile.Name
			output.Metadata[hub.PROFILE_VERSION_KEY] = profile.Version
			output.Metadata[hub.RUNTIME_VERSION_KEY] = hub.SourceVersion

			for e := range endpoints {

				err = endpoints[e].Write(output)

				if err != nil {
					log.Println(err)
				}

			}
		}

	}
}
