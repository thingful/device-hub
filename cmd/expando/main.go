package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/thingful/expando"
	"github.com/thingful/expando/engine"
	"github.com/thingful/expando/pipe"
	"github.com/yosssi/gmq/mqtt/client"
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

	script := expando.Script{
		Main:     "decode",
		Runtime:  expando.Javascript,
		Input:    expando.JSON,
		Contents: scriptContents,
	}

	ctx, cancel := context.WithCancel(context.Background())

	var broker pipe.Broker

	if in == "std" {

		broker = pipe.FromStdIn(cancel)

	}
	if in == "mqtt" {

		// TODO : pick up connection options from somewhere
		options := &client.ConnectOptions{
			Network:  "tcp",
			Address:  "0.0.0.0:1883",
			ClientID: []byte("expando-client"),
		}

		var err error
		broker, err = pipe.FromMQTT(options)

		if err != nil {
			exitWithError(err)
		}
	}

	if broker == nil {
		exitWithError(errors.New("unable to create broker"))
	}

	defer broker.Close()
	channel := broker.Channel()
	errorChannel := channel.Errors()

	for {

		select {

		case <-ctx.Done():
			return

		case err := <-errorChannel:
			exitWithError(err)

		case input := <-channel.Out():

			output, err := scripter.Execute(script, input)

			if err != nil {
				exitWithError(err)
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
