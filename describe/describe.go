// Copyright © 2017 thingful

package describe

import (
	"fmt"

	"github.com/spf13/cast"
)

type typez string

const (
	Int32  = typez("int32")
	Int64  = typez("int64")
	String = typez("string")
	Url    = typez("url")
	Bool   = typez("bool")
)

// Parameter describes a configuration parameter
type Parameter struct {
	Name        string
	Type        typez
	Required    bool
	Description string
	Default     string
	Examples    []string
}

// Parameters are a collection of Paramater structs
type Parameters []Parameter

// Describe returns a string description of the Parameter
func (p Parameter) Describe() string {

	if p.Default == "" {
		return fmt.Sprintf("%s : %s (Required : %v) %s, Examples %v ", p.Name, p.Type, p.Required, p.Description, p.Examples)
	}

	return fmt.Sprintf("%s : %s (Required : %v, Default : %s) %s, Examples %v ", p.Name, p.Type, p.Required, p.Default, p.Description, p.Examples)
}

func Validate(config map[string]string, params Parameters) error {
	/*
		for k, v := range config {

			// is in list of params? no - fail

			// bit of type checking e.g. urls, ints etc.

		}

		for _, param := range params {

			// p.Required and not in config - fail
		}
	*/
	return nil
}

func CreateValues(config map[string]string, params Parameters) (Values, error) {

	return Values{}, nil
}

// Value contains the Parameter type description alongside its value
type Value struct {
	Parameter
	Value interface{}
}

// Values are a collection of Value structs
type Values map[string]Value

func (v Values) String(key string) (string, bool) {

	value, found := v[key]

	if !found {
		return "", false
	}

	str, err := cast.ToStringE(value.Value)

	if err != nil {
		return "", false
	}

	return str, true
}

func (v Values) MustString(key string) string {

	value, found := v.String(key)

	if !found {
		panic(fmt.Errorf("string value with key %s not found", key))
	}

	return value
}

func (v Values) Bool(key string) (bool, bool) {

	value, found := v[key]

	if !found {
		return false, false
	}

	b, err := cast.ToBoolE(value.Value)

	if err != nil {
		return false, false
	}

	return b, true
}

func (v Values) BoolWithDefault(key string, defaultValue bool) bool {

	value, found := v.Bool(key)

	if !found {

		return defaultValue
	}

	return value
}

func (v Values) Int(key string) (int, bool) {

	value, found := v[key]

	if !found {
		return 0, false
	}

	i, err := cast.ToIntE(value.Value)

	if err != nil {
		return 0, false
	}

	return i, true

}

func (v Values) IntWithDefault(key string, defaultValue int) int {

	value, found := v.Int(key)

	if !found {

		return defaultValue
	}

	return value
}
