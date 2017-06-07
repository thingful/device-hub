// Copyright © 2017 thingful

package describe

import (
	"fmt"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/cast"
)

type typez string

const (

	// Int32 type is an signed int32
	Int32 = typez("int32")

	// Int64 type is an signed int64
	Int64 = typez("int64")

	// String type is a string
	String = typez("string")

	// Url type is a value that is a valid URL
	Url = typez("url")

	// Bool type is a true/false value
	Bool = typez("bool")

	// Float32 type is a float32
	Float32 = typez("float32")
)

// Parameter describes a configuration parameter
type Parameter struct {
	Name        string
	Type        typez
	Required    bool
	Description string
	Default     interface{}
	Examples    []string
}

// Parameters are a collection of Paramater structs
type Parameters []Parameter

// Describe returns a string description of the Parameter
func (p Parameter) Describe() string {

	if p.Default == nil {
		return fmt.Sprintf("%s : %s (Type : %s, Required : %v, Examples %v )", p.Name, p.Description, p.Type, p.Required, p.Examples)
	}

	return fmt.Sprintf("%s : %s (Type : %s, Required : %v, Default : %v, Examples %v )", p.Name, p.Description, p.Type, p.Required, p.Default, p.Examples)
}

func NewValues(config map[string]string, params Parameters) (Values, error) {

	values := Values{
		collection: map[string]Value{},
	}

	// validate params
	for _, p := range params {

		if p.Required {

			v, found := config[p.Name]

			if !found {
				return values, fmt.Errorf("%s is Required but not supplied", p.Name)
			}
			values.collection[p.Name] = Value{Parameter: p, Value: v}
		} else {

			v, found := config[p.Name]

			if found {
				values.collection[p.Name] = Value{Parameter: p, Value: v}
			}
		}
	}
	// validate types
	for k, v := range values.collection {

		switch v.Type {

		case String:

			_, ok := values.String(k)

			if !ok {
				return values, fmt.Errorf("%s is not of type String", k)
			}

		case Int32:

			_, ok := values.Int32(k)

			if !ok {
				return values, fmt.Errorf("%s is not of type Int32", k)
			}

		case Int64:

			_, ok := values.Int64(k)

			if !ok {
				return values, fmt.Errorf("%s is not of type Int64", k)
			}

		case Bool:
			_, ok := values.Bool(k)

			if !ok {
				return values, fmt.Errorf("%s is not of type Bool", k)
			}

		case Url:
			_, ok := values.Url(k)

			if !ok {
				return values, fmt.Errorf("%s is not of type Url", k)
			}
		case Float32:

			_, ok := values.Float32(k)

			if !ok {
				return values, fmt.Errorf("%s is not of type Float32", k)
			}

		default:
			return values, fmt.Errorf("unknown type : %s", string(v.Type))
		}
	}

	return values, nil
}

// Value contains the Parameter type description alongside its value
type Value struct {
	Parameter
	Value interface{}
}

// Values are a collection of Value struct
type Values struct {
	collection map[string]Value
	sync.Mutex
}

// String returns a string and true if a value exists and can be cast to a string
func (v Values) String(key string) (string, bool) {

	v.Lock()
	value, found := v.collection[key]
	v.Unlock()

	if !found {
		return "", false
	}

	str, err := cast.ToStringE(value.Value)

	if err != nil {
		return "", false
	}

	return str, true
}

// MustString returns the value as a string or panics
func (v Values) MustString(key string) string {

	value, found := v.String(key)

	if !found {
		panic(fmt.Errorf("string value with key %s not found", key))
	}

	return value
}

// Bool returns a boolean and true if a value exists and can be cast to a boolean
func (v Values) Bool(key string) (bool, bool) {

	v.Lock()
	value, found := v.collection[key]
	v.Unlock()

	if !found {
		return false, false
	}

	b, err := cast.ToBoolE(value.Value)

	if err != nil {
		return false, false
	}

	return b, true
}

// BoolWithDefault returns a boolean or a default value if not found
func (v Values) BoolWithDefault(key string, defaultValue bool) bool {

	value, found := v.Bool(key)

	if !found {

		return defaultValue
	}

	return value
}

// Int32 returns a int32 and true if a value exists and can be cast to an int32
func (v Values) Int32(key string) (int32, bool) {

	v.Lock()
	value, found := v.collection[key]
	v.Unlock()

	if !found {
		return 0, false
	}

	i, err := cast.ToInt32E(value.Value)

	if err != nil {
		return 0, false
	}

	return i, true

}

// Int32WithDefault returns an integer or a default value if not found
func (v Values) Int32WithDefault(key string, defaultValue int32) int32 {

	value, found := v.Int32(key)

	if !found {

		return defaultValue
	}

	return value
}

// Int64 returns a int64 and true if a value exists and can be cast to an int64
func (v Values) Int64(key string) (int64, bool) {

	v.Lock()
	value, found := v.collection[key]
	v.Unlock()

	if !found {
		return 0, false
	}

	i, err := cast.ToInt64E(value.Value)

	if err != nil {
		return 0, false
	}

	return i, true

}

// Url returns a Url and true if a value exists and can be cast to an url
func (v Values) Url(key string) (string, bool) {

	value, found := v.String(key)

	if !found {
		return "", false
	}

	valid := govalidator.IsRequestURL(value)
	if !valid {
		return "", false
	}

	return value, true

}

// Float32 returns a float32 and true if a value exists and can be cast to an float32
func (v Values) Float32(key string) (float32, bool) {

	v.Lock()
	value, found := v.collection[key]
	v.Unlock()

	if !found {
		return 0, false
	}

	i, err := cast.ToFloat32E(value.Value)

	if err != nil {
		return 0, false
	}

	return i, true

}

// Float32WithDefault returns a float32 or a default value if not found
func (v Values) Float32WithDefault(key string, defaultValue float32) float32 {

	value, found := v.Float32(key)

	if !found {

		return defaultValue
	}

	return value
}
