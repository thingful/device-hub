package config

import (
	"encoding/json"
	"io/ioutil"
)

// Unmarshal de-serialises a configuration from JSON form
func Unmarshal(bytes []byte) (*Configuration, error) {

	conf := Configuration{}
	err := json.Unmarshal(bytes, &conf)

	return &conf, err
}

// Marshal serialises a configuration to JSON
func Marshal(conf Configuration) ([]byte, error) {
	return json.Marshal(conf)
}

// LoadProfile loads a configuration from a location
func LoadProfile(path string) (*Configuration, error) {

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return Unmarshal(bytes)
}
