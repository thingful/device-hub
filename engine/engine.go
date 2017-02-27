package engine

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"bitbucket.org/tsetsova/decode-prototype/hub/expando"
	"github.com/robertkrimen/otto"
)

var (
	ScriptTimedOutErr = errors.New("script timed out")
)

func New() engine {
	return engine{}
}

type engine struct{}

func (e engine) Execute(script expando.Script, input expando.Input) (interface{}, error) {

	var output otto.Value
	var err error

	if script.Input == expando.Raw {
		env := map[string]interface{}{
			"__input": input.Payload,
		}
		output, err = e.run(script.Contents, "decode(__input);", env, time.Second*1)

	}
	if script.Input == expando.CSV {

		env, err := prepareCSV(input)

		if err != nil {
			return nil, err
		}

		output, err = e.run(script.Contents, "decode(__header, __lines);", env, time.Second*1)

	}

	if script.Input == expando.JSON {

		env := map[string]interface{}{
			"__input": string(input.Payload),
		}
		output, err = e.run(script.Contents, "decode( JSON.parse( __input ));", env, time.Second*1)

	}

	if err != nil {
		return nil, err
	}

	gov, err := output.Export()

	if err != nil {
		return nil, err
	}

	return gov, nil
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

func prepareCSV(input expando.Input) (map[string]interface{}, error) {

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
