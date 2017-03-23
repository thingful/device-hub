package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

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

	var broker pipe.Broker

	if in == "std" {

		broker = pipe.FromStdIn(cancel)

	}
	if in == "mqtt" {

		// TODO : pick up connection options from somewhere
		opts := mqtt.NewClientOptions()
		opts.AddBroker("tcp://0.0.0.0:1883")
		opts.SetClientID("device-hub")
		opts.SetKeepAlive(2 * time.Second)
		opts.SetPingTimeout(1 * time.Second)

		var err error
		broker, err = pipe.FromMQTT(opts, "#")

		if err != nil {
			exitWithError(err)
		}
	}

	if broker == nil {
		exitWithError(errors.New("unable to create broker"))
	}

	defer broker.Close()
	channel, err := broker.Channel()

	if err != nil {
		exitWithError(err)
	}

	errorChannel := channel.Errors()

	for {

		select {

		case <-ctx.Done():
			return

		case err := <-errorChannel:
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
