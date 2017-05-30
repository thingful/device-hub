// Copyright © 2017 thingful

package describe

import "fmt"

type typez string

const (
	Int32  = typez("int32")
	Int64  = typez("int64")
	String = typez("string")
	Url    = typez("url")
	Bool   = typez("bool")
)

type Parameter struct {
	Name        string
	Type        typez
	Required    bool
	Description string
	Default     string
	Examples    []string
}

type Parameters []Parameter

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

type Value struct {
	Parameter
	Value interface{}
}

type Values []Values

func (v Values) String(key string) (bool, string) {
	return true, "booya"
}
func (v Values) MustString(key string) string {
	return "booya"
}

func (v Values) BoolWithDefault(key string, defaultv bool) bool {
	return false
}

func (v Values) IntWithDefault(key string, defaultv int) int {
	return 1
}
