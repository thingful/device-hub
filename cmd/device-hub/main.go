package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/pipe"
)

var (
	SourceVersion string
)

func main() {

	var scriptContents string
	var in string
	var out string
	var showVersion bool

	flag.StringVar(&in, "in", "", "read from specified input. Known values are 'std' or 'mqtt'.")
	flag.StringVar(&out, "out", "std", "output to specified stream.")
	flag.StringVar(&scriptContents, "script", "function decode( input ){ return input }", "js to transform input.")
	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Parse()

	if showVersion {
		fmt.Println(SourceVersion)
		return
	}

	if in == "" {
		exitWithError(errors.New("must specify an -in flag"))
	}
	if out == "" {
		exitWithError(errors.New("must specify an -out flag"))
	}

	scripter := engine.New()

	script := hub.Script{
		Main:     "decode",
		Runtime:  hub.Javascript,
		Input:    hub.JSON,
		Contents: scriptContents,
	}

	ctx, cancel := context.WithCancel(context.Background())

	var channel pipe.Channel
	var err error

	router := pipe.DefaultRouter()

	if in == "std" {

		channel = pipe.NewStdInChannel(cancel)

	}
	if in == "mqtt" {

		// TODO : pick up connection options from somewhere
		clientName := fmt.Sprintf("device-hub-%s", SourceVersion)
		options := pipe.DefaultMQTTOptions("tcp://0.0.0.0:1883", clientName)
		client := pipe.DefaultClient(options)

		// TODO : set sensible wait time
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			exitWithError(token.Error())
		}

		// TODO : set a sensible timeout
		defer client.Disconnect(1000)

		channel, err = pipe.NewMQTTChannel("/xxx", client)

		if err != nil {
			exitWithError(err)
		}
	}

	if in == "http" {

		channel = pipe.NewHTTPChannel("/xxx", router)

		pipe.StartDefaultHTTPListener(ctx, router, ":8085")

	}

	if channel == nil {
		exitWithError(errors.New("unable to create channel"))
	}

	for {

		select {

		case <-ctx.Done():
			return

		case err := <-channel.Errors():
			log.Println(err)

		case input := <-channel.Out():

			output, err := scripter.Execute(script, input)

			if err != nil {
				log.Println(err)
				continue
			}

			bytes, err := json.Marshal(output)
			if err != nil {
				exitWithError(err)
			}

			fmt.Println(string(bytes))
		}
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
