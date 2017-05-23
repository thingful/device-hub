// Copyright © 2017 thingful

package engine

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/robertkrimen/otto"

	hub "github.com/thingful/device-hub"
)

var (
	ScriptTimedOutErr = errors.New("script timed out")
)

func New() engine {
	return engine{}
}

type engine struct{}

func (e engine) Execute(script Script, input hub.Message) (out hub.Message, err error) {

	var output otto.Value

	if script.Input == Raw {
		env := map[string]interface{}{
			"__input": input.Payload,
		}
		main := fmt.Sprintf("%s (__input);", script.Main)
		output, err = e.run(script.Contents, main, env, time.Second*1)

	}
	if script.Input == CSV {

		env, err := prepareCSV(input)

		if err != nil {
			return input, err
		}

		main := fmt.Sprintf("%s (__header, __lines);", script.Main)
		output, err = e.run(script.Contents, main, env, time.Second*1)

	}

	if script.Input == JSON {

		env := map[string]interface{}{
			"__input": string(input.Payload),
		}

		main := fmt.Sprintf("%s ( JSON.parse( __input ));", script.Main)
		output, err = e.run(script.Contents, main, env, time.Second*1)

	}

	if err != nil {
		return input, err
	}

	gov, err := output.Export()

	if err != nil {
		return input, err
	}

	input.Output = gov

	return input, nil
}

// run javascript engine in-process
// TODO : think about tolerance to malicious input, fault tolerance
func (engine) run(code, main string, env map[string]interface{}, timeout time.Duration) (val otto.Value, err error) {

	vm := otto.New()

	// load the environment
	for key, val := range env {
		vm.Set(key, val)
	}

	vm.Set("__log", func(call otto.FunctionCall) otto.Value {
		// TODO : allow logger to be defined
		log.Println(call.ArgumentList)
		return otto.UndefinedValue()
	})
	vm.Run("console.log = __log")

	defer func() {

		if caught := recover(); caught != nil {
			val = otto.Value{}
			if caught == ScriptTimedOutErr {
				err = ScriptTimedOutErr
			} else {
				err = fmt.Errorf("fatal error in function: %s", caught)
			}
		}
	}()

	vm.Interrupt = make(chan func(), 1)

	go func() {
		time.Sleep(timeout)
		vm.Interrupt <- func() {
			panic(ScriptTimedOutErr)
		}
	}()

	return vm.Run(fmt.Sprintf("%s;\n %s", code, main))
}

func prepareCSV(input hub.Message) (map[string]interface{}, error) {

	r := csv.NewReader(strings.NewReader(string(input.Payload)))

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// TODO : think about max records or map/reducing it
	return map[string]interface{}{
		"__header": records[0],
		"__lines":  records[1:],
	}, nil

}
