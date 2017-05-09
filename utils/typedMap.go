// Copyright © 2017 thingful

package utils

import (
	"fmt"

	"github.com/spf13/cast"
)

// TypedMap provides strongly typed access to a map[string]string
type TypedMap map[string]string

// String return true if a string value exists in the collection along with the value
func (c TypedMap) String(key string) (bool, string) {

	v, f := c[key]

	if !f {
		return false, ""
	}

	str, err := cast.ToStringE(v)

	if err != nil {
		return false, ""
	}

	return true, str
}

// MString returns the value as a string for the key or panics
func (c TypedMap) MString(key string) string {

	found, v := c.String(key)

	if !found {
		panic(fmt.Errorf("string value with key %s not found", key))
	}

	return v
}

// DString returns the value as a string of the key or a default value
func (c TypedMap) DString(key, defaultValue string) string {

	found, v := c.String(key)

	if !found {
		return defaultValue
	}

	return v
}

func (c TypedMap) Bool(key string) (bool, bool) {

	v, f := c[key]

	if !f {
		return false, false
	}

	b, err := cast.ToBoolE(v)

	if err != nil {
		return false, false
	}

	return true, b
}

func (c TypedMap) MBool(key string) bool {

	found, v := c.Bool(key)

	if !found {

		panic(fmt.Errorf("bool value with key %s not found", key))

	}
	return v
}

func (c TypedMap) DBool(key string, defaultValue bool) bool {

	found, v := c.Bool(key)

	if !found {

		return defaultValue
	}

	return v
}
