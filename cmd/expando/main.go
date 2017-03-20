package main

import (
	"context"
	"encoding/json"
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
	var showVersion bool

	flag.StringVar(&scriptContents, "script", "function decode( input ){ return input }", "js to transform input")
	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Parse()

	if showVersion {
		fmt.Println(SourceVersion)
		return
	}

	scripter := engine.New()

	script := expando.Script{
		Main:     "decode",
		Runtime:  expando.Javascript,
		Input:    expando.JSON,
		Contents: scriptContents,
	}

	ctx, _ := context.WithCancel(context.Background())

	//	broker := pipe.FromStdIn(ctx)

	broker := pipe.FromMQTT(ctx)

	channel := broker.Channel()
	errorChannel := channel.Errors()

	go channel.Next()

	for {

		select {

		case <-ctx.Done():
			fmt.Println("context closed")
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
