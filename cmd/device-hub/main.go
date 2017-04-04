// Copyright © 2017 thingful

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/config"
	"github.com/thingful/device-hub/endpoint"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/listener"
	"github.com/thingful/go/file"
)

var (
	SourceVersion = "DEVELOPMENT"
)

func main() {

	var showVersion bool
	var configurationPath string

	flag.StringVar(&configurationPath, "config", "./config.json", "path to a json configuration file.")
	flag.BoolVar(&showVersion, "version", false, "show version.")

	flag.Parse()

	if showVersion {
		fmt.Println(SourceVersion)
		return
	}

	// TODO : ensure this path is constrained to a few well known paths
	if !file.Exists(configurationPath) {
		exitWithError(fmt.Errorf("configuration at %s doesn't exist", configurationPath))
	}

	configuration, err := config.LoadProfile(configurationPath)

	if err != nil {
		exitWithError(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	for _, pipe := range configuration.Pipes {

		found, listenerConf := configuration.Listeners.FindByUID(pipe.Listener)

		if !found {
			exitWithError(fmt.Errorf("listener with UID %s not found", pipe.Listener))
		}

		endpoints := []hub.Endpoint{}

		for e := range pipe.Endpoints {

			found, endpointConf := configuration.Endpoints.FindByUID(pipe.Endpoints[e])

			if !found {
				exitWithError(fmt.Errorf("endpoint with UID %s not found", pipe.Endpoints[e]))
			}

			if endpointConf.Type == "stdout" {

				prettyPrint := endpointConf.Configuration.DBool("prettyPrint", false)

				endpoints = append(endpoints, endpoint.NewStdOutEndpoint(prettyPrint))
			}

		}

		found, profile := configuration.Profiles.FindByUID(pipe.Profile)

		if !found {
			exitWithError(fmt.Errorf("profile with UID %s not found", pipe.Profile))
		}

		listener, err := StartListener(listenerConf, cancel)

		if err != nil {
			exitWithError(err)
		}

		channel, err := listener.NewChannel(pipe.Uri)

		if err != nil {
			exitWithError(err)
		}

		go StartPipe(ctx, listener, channel, profile, endpoints)

	}

	<-ctx.Done()

}

func StartPipe(ctx context.Context, listener hub.Listener, channel hub.Channel, profile config.Profile, endpoints []hub.Endpoint) {

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

func StartListener(endpoint config.Endpoint, cancel context.CancelFunc) (hub.Listener, error) {

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

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
