package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

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

	inputChannel := broker.Channel()
	errorChannel := broker.Errors()

	go broker.Next()

	for {

		select {

		case <-ctx.Done():
			return

		case err := <-errorChannel:
			panic(err)

		case input := <-inputChannel:

			output, err := scripter.Execute(script, input)

			if err != nil {
				panic(err)
			}

			bytes, err := json.Marshal(output)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(bytes))
			cancel()
		}
	}
}
