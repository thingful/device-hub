// Copyright © 2017 thingful

package engine

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/utils"
)

func TestRawDecodeValid(t *testing.T) {

	t.Parallel()

	script := Script{
		Main:    "xxx",
		Runtime: Javascript,
		Input:   Raw,
		Contents: `function xxx (input) {
				return {
						'value' : ((input[0] << 8) | input[1]) / 100,
					}
			}`,
	}

	// create the input payload as a byte array
	buf := &bytes.Buffer{}

	// multiple * 100 to ensure non floating point values
	initialValue := 22.33
	binary.Write(buf, binary.BigEndian, int16(initialValue*100))

	input := hub.Message{Payload: buf.Bytes(), Metadata: map[string]interface{}{}}

	e := New(utils.NewNoOpLogger())
	result, err := e.Execute(script, input)

	assert.Nil(t, err)

	resultAsMap := result.Output.(map[string]interface{})
	value := resultAsMap["value"]

	assert.Equal(t, initialValue, value)
}

func TestCSVDecodeValid(t *testing.T) {

	t.Parallel()

	script := Script{
		Main:    "xxx",
		Runtime: Javascript,
		Input:   CSV,
		Contents: `function xxx (header, lines) {
				return {
						'header' : header,
						'lines' : lines
					}
			}`,
	}

	csv := "column1, column2\none, two\nthree, four\n five,six"
	input := hub.Message{Payload: []byte(csv), Metadata: map[string]interface{}{}}

	e := New(utils.NewNoOpLogger())
	result, err := e.Execute(script, input)
	assert.Nil(t, err)
	resultAsMap := result.Output.(map[string]interface{})

	assert.Len(t, resultAsMap["header"], 2)
	assert.Len(t, resultAsMap["lines"], 3)

}

func TestJSONDecodeValid(t *testing.T) {

	t.Parallel()

	script := Script{
		Main:    "xxx",
		Runtime: Javascript,
		Input:   JSON,
		Contents: `function xxx (input) {
				return input
			}`,
	}

	json := "{ \"a\" : 1}"
	input := hub.Message{Payload: []byte(json), Metadata: map[string]interface{}{}}

	e := New(utils.NewNoOpLogger())
	result, err := e.Execute(script, input)

	assert.Nil(t, err)
	resultAsMap := result.Output.(map[string]interface{})

	assert.NotNil(t, resultAsMap["a"])
	assert.Equal(t, resultAsMap["a"], float64(1))
}

func TestEngineAuxFunction(t *testing.T) {

	t.Parallel()

	script := Script{
		Main:    "xxx",
		Runtime: Javascript,
		Input:   Raw,
		Contents: `function xxx () {
				return today()
			}`,
	}

	e := New(utils.NewNoOpLogger())

	e.AuxFuncs["today"] = func(call otto.FunctionCall) otto.Value {
		ti := time.Now().Format("01/02/2006")
		val, err := otto.ToValue(ti)
		assert.Nil(t, err)
		return val
	}

	result, err := e.Execute(script, hub.Message{})
	assert.Nil(t, err)

	assert.Equal(t, result.Output, time.Now().Format("01/02/2006"))
}

func TestEngineAuxObject(t *testing.T) {

	t.Parallel()

	script := Script{
		Main:    "xxx",
		Runtime: Javascript,
		Input:   Raw,
		Contents: `function xxx () {
				return geolocation.Coords
			}`,
	}

	e := New(utils.NewNoOpLogger())

	pos := new(position)
	pos.Coords.Latitude = 54.416333
	pos.Coords.Longitude = -4.36552

	e.AuxObjs["geolocation"] = pos

	result, err := e.Execute(script, hub.Message{})
	assert.Nil(t, err)

	coords := result.Output.(coordinates)
	assert.Equal(t, pos.Coords.Latitude, coords.Latitude)
	assert.Equal(t, pos.Coords.Longitude, coords.Longitude)
}
