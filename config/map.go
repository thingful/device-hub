package config

import "fmt"

type ConfigMap map[string]interface{}

func (c ConfigMap) String(key string) (bool, string) {

	v, f := c[key]

	if !f {
		return false, ""
	}

	// TODO : need smarter coercian
	return true, v.(string)
}

func (c ConfigMap) MString(key string) string {

	found, v := c.String(key)
	if !found {
		panic(fmt.Errorf("value with key %s not found", key))
	}
	return v
}

func (c ConfigMap) Bool(key string) (bool, bool) {

	v, f := c[key]

	if !f {
		return false, false
	}

	// TODO : need smarter coercian
	return true, v.(bool)
}

func (c ConfigMap) MBool(key string) bool {

	found, v := c.Bool(key)
	if !found {
		panic(fmt.Errorf("value with key %s not found", key))
	}
	return v
}

func (c ConfigMap) DBool(key string, defaultValue bool) bool {

	found, v := c.Bool(key)
	if !found {
		return defaultValue
	}
	return v
}
