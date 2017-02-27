package main

import (
	"bufio"
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

	in, err := getInputFromStdIn()

	if err != nil {
		panic(err)
	}

	scripter := engine.New()
	input := expando.Input{Payload: []byte(in)}

	script := expando.Script{
		Main:     "decode",
		Runtime:  expando.Javascript,
		Input:    expando.JSON,
		Contents: scriptContents,
	}

	output, err := scripter.Execute(script, input)

	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

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
