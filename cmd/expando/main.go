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
)

var (
	SourceVersion string
)

func main() {

	var scriptContents string
	var in string
	var out string
	var showVersion bool

	flag.StringVar(&in, "in", "", "read from specified input e.g. std or mqtt")
	flag.StringVar(&out, "out", "std", "output to specified stream.")
	flag.StringVar(&scriptContents, "script", "function decode( input ){ return input }", "js to transform input")
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
		broker = pipe.FromMQTT(ctx)
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
			//broker.Close()
			//cancel()
		}
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
