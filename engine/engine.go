// Copyright © 2017 thingful

package engine

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/robertkrimen/otto"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/utils"
)

var (
	// ScriptTimedOutErr returned when the script takes longer then the MaxScriptDuration
	ScriptTimedOutErr = errors.New("script timed out")

	// MaxScriptDuration defines the longest a script can run for
	MaxScriptDuration = time.Second * 1
)

// New returns a script runner
func New(logger utils.Logger) engine {
	return engine{
		logger:   logger,
		AuxFuncs: make(map[string]func(otto.FunctionCall) otto.Value),
		AuxObjs:  make(map[string]interface{}),
	}
}

type engine struct {
	logger   utils.Logger
	AuxFuncs map[string]func(otto.FunctionCall) otto.Value
	AuxObjs  map[string]interface{}
}

// Execute takes a script and a message - like a method and function arguments
func (e engine) Execute(script Script, input hub.Message) (out hub.Message, err error) {

	var output otto.Value

	if script.Input == Raw {
		env := map[string]interface{}{
			"__input": input.Payload,
		}
		main := fmt.Sprintf("%s (__input);", script.Main)
		output, err = e.run(script.Contents, main, env, MaxScriptDuration)

	}
	if script.Input == CSV {

		env, err := prepareCSV(input)

		if err != nil {
			return input, err
		}

		main := fmt.Sprintf("%s (__header, __lines);", script.Main)
		output, err = e.run(script.Contents, main, env, MaxScriptDuration)

	}

	if script.Input == JSON {

		env := map[string]interface{}{
			"__input": string(input.Payload),
		}

		main := fmt.Sprintf("%s ( JSON.parse( __input ));", script.Main)
		output, err = e.run(script.Contents, main, env, MaxScriptDuration)

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
func (e engine) run(code, main string, env map[string]interface{}, timeout time.Duration) (val otto.Value, err error) {

	vm := otto.New()

	// load the environment
	for key, val := range env {
		vm.Set(key, val)
	}

	vm.Set("__log", func(call otto.FunctionCall) otto.Value {
		e.logger.Info(call.ArgumentList)
		return otto.UndefinedValue()
	})
	vm.Run("console.log = __log")

	// Set aux objects
	for oName, o := range e.AuxObjs {
		if v, err := vm.ToValue(o); err == nil {
			err = vm.Set(oName, v)
		}
		if err != nil {
			return
		}
	}
	// Set aux functions
	for fname, f := range e.AuxFuncs {
		if err = vm.Set(fname, f); err != nil {
			return
		}
	}

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

func (e engine) SetGeoLocation(lat, lng float64) {
	pos := new(position)
	pos.setLocation(lat, lng)
	e.AuxObjs["geolocation"] = pos
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
