package profile

import (
	"encoding/json"
	"io/ioutil"

	hub "github.com/thingful/device-hub"
)

type UID string

type Configuration struct {
	Listeners Endpoints `json:"listeners"`
	Endpoints Endpoints `json:"endpoints"`
	Profiles  Profiles  `json:"profiles"`
	Pipes     []Pipe    `json:"pipes"`
}

type Endpoint struct {
	Type          string                 `json:"type"`
	UID           UID                    `json:"uid"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

type Endpoints []Endpoint

func (e Endpoints) FindByUID(uid UID) (bool, Endpoint) {

	for _, endpoint := range e {
		if endpoint.UID == uid {
			return true, endpoint
		}
	}
	return false, Endpoint{}
}

type Profile struct {
	UID         UID    `json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// TODO : make this a semantic triple
	Version string     `json:"version"`
	Script  hub.Script `json:"script"`
}

type Profiles []Profile

func (p Profiles) FindByUID(uid UID) (bool, Profile) {

	for _, profile := range p {
		if profile.UID == uid {
			return true, profile
		}
	}
	return false, Profile{}
}

type Pipe struct {
	Uri      string `json:"uri"`
	Profile  UID    `json:"profile"`
	Listener UID    `json:"listener"`
	Endpoint UID    `json:"endpoint"`
}

func Unmarshal(bytes []byte) (*Configuration, error) {

	conf := Configuration{}
	err := json.Unmarshal(bytes, &conf)

	return &conf, err
}

func Marshal(conf Configuration) ([]byte, error) {
	return json.Marshal(conf)
}

func LoadProfile(path string) (*Configuration, error) {

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return Unmarshal(bytes)
}
