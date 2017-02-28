package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"bitbucket.org/tsetsova/decode-prototype/hub/expando"
	"bitbucket.org/tsetsova/decode-prototype/hub/expando/engine"
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

	broker := FromStdIn(ctx)

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

func FromStdIn(ctx context.Context) *stdin {

	channel := make(chan expando.Input)
	errors := make(chan error)

	return &stdin{
		ctx:     ctx,
		channel: channel,
		errors:  errors,
	}
}

type stdin struct {
	ctx     context.Context
	channel chan expando.Input
	errors  chan error
}

func (s *stdin) Channel() chan expando.Input {
	return s.channel
}

func (s *stdin) Errors() chan error {
	return s.errors
}

func (s *stdin) Next() {

	contents, err := getInputFromStdIn()

	if err != nil {
		s.errors <- err
	} else {
		s.channel <- expando.Input{Payload: []byte(contents)}
	}
}

// if we are being piped some input return it else error
func getInputFromStdIn() (string, error) {

	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return "", errors.New("input expected from stdin e.g. echo {} | ./expando")
	}

	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}
