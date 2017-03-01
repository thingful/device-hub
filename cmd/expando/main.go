package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"bitbucket.org/tsetsova/decode-prototype/hub/expando"
	"bitbucket.org/tsetsova/decode-prototype/hub/expando/engine"
	"bitbucket.org/tsetsova/decode-prototype/hub/expando/pipe"
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

	ctx, cancel := context.WithCancel(context.Background())

	broker := pipe.FromStdIn(ctx)

	channel := broker.Channel()
	errorChannel := channel.Errors()

	go channel.Next()

	for {

		select {

		case <-ctx.Done():
			return

		case err := <-errorChannel:
			exitWithError(err)

		case input := <-channel.Out:

			output, err := scripter.Execute(script, input)

			if err != nil {
				exitWithError(err)
			}

			bytes, err := json.Marshal(output)
			if err != nil {
				exitWithError(err)
			}

			fmt.Println(string(bytes))
			cancel()
		}
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
