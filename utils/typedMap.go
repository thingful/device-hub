// Copyright © 2017 thingful

package utils

import (
	"fmt"

	"github.com/spf13/cast"
)

type TypedMap map[string]interface{}

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

func (c TypedMap) MString(key string) string {

	found, v := c.String(key)

	if !found {
		panic(fmt.Errorf("value with key %s not found", key))
	}

	return v
}

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

		panic(fmt.Errorf("value with key %s not found", key))

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
