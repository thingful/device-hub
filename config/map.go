package config

import (
	"fmt"

	"github.com/spf13/cast"
)

type configMap map[string]interface{}

func (c configMap) String(key string) (bool, string) {

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

func (c configMap) MString(key string) string {

	found, v := c.String(key)

	if !found {
		panic(fmt.Errorf("value with key %s not found", key))
	}

	return v
}

func (c configMap) DString(key, defaultValue string) string {

	found, v := c.String(key)

	if !found {
		return defaultValue
	}

	return v
}

func (c configMap) Bool(key string) (bool, bool) {

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

func (c configMap) MBool(key string) bool {

	found, v := c.Bool(key)

	if !found {

		panic(fmt.Errorf("value with key %s not found", key))

	}
	return v
}

func (c configMap) DBool(key string, defaultValue bool) bool {

	found, v := c.Bool(key)

	if !found {

		return defaultValue
	}

	return v
}
