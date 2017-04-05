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
	"github.com/thingful/device-hub/listener"
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

	ctx, cancel := context.WithCancel(context.Background())

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

			// TODO : cache existing endpoints
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

		// TODO : keep a list of successfully started endpoints
		listener, err := startListener(listenerConf, cancel)

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
			output.Metadata[hub.RUNTIME_VERSION_KEY] = SourceVersion

			for e := range endpoints {

				err = endpoints[e].Write(output)

				if err != nil {
					log.Println(err)
				}

			}
		}

	}
}

func startListener(endpoint config.Endpoint, cancel context.CancelFunc) (hub.Listener, error) {

	if endpoint.Type == "std" {

		return listener.NewStdInListener(cancel)
	}
	if endpoint.Type == "mqtt" {

		clientName := fmt.Sprintf("device-hub-%s", SourceVersion)

		brokerAddress := endpoint.Configuration.MString("MQTTBrokerAddress")

		options := listener.DefaultMQTTOptions(brokerAddress, clientName)
		client := listener.DefaultMQTTClient(options)

		// TODO : set sensible wait time
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			exitWithError(token.Error())
		}

		return listener.NewMQTTListener(client)
	}

	if endpoint.Type == "http" {

		binding := endpoint.Configuration.MString("HTTPBindingAddress")
		return listener.NewHTTPListener(binding)
	}

	return nil, fmt.Errorf("listener of type %s not found", endpoint.Type)

}
