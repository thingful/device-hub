package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/config"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/pipe"
	"github.com/thingful/go/file"
)

var (
	SourceVersion string
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

	// TODO : ensure ths path is constrained to a few well known paths
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

		// TDDO : stop hardcoding the endpoint in
		//found, endpointConf := configuration.Endpoints.FindByUID(pipe.Endpoint)

		//if !found {
		//	exitWithError(fmt.Errorf("endpoint with UID %s not found", pipe.Endpoint))
		//}

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

		go StartPipe(ctx, listener, channel, profile)

	}

	<-ctx.Done()

}

func StartPipe(ctx context.Context, listener hub.Listener, channel hub.Channel, profile config.Profile) {

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
				continue
			}

			bytes, err := json.Marshal(output)
			if err != nil {
				exitWithError(err)
			}

			_, err = pipe.WriteToStdOut(bytes)

			if err != nil {
				log.Println(err)
			}
		}
	}

}

func StartListener(endpoint config.Endpoint, cancel context.CancelFunc) (hub.Listener, error) {

	if endpoint.Type == "std" {

		return pipe.NewStdInListener(cancel)
	}
	if endpoint.Type == "mqtt" {

		clientName := fmt.Sprintf("device-hub-%s", SourceVersion)

		brokerAddress := endpoint.Configuration.MString("MQTTBrokerAddress")

		options := pipe.DefaultMQTTOptions(brokerAddress, clientName)
		client := pipe.DefaultMQTTClient(options)

		// TODO : set sensible wait time
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			exitWithError(token.Error())
		}

		return pipe.NewMQTTListener(client)
	}

	if endpoint.Type == "http" {

		binding := endpoint.Configuration.MString("HTTPBindingAddress")
		return pipe.NewHTTPListener(binding)
	}

	return nil, fmt.Errorf("listener of type %s not found", endpoint.Type)

}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
