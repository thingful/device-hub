package config

import (
	"fmt"

	"github.com/thingful/device-hub/engine"
)

type UID string

type Configuration struct {
	Listeners Endpoints `json:"listeners"`
	Endpoints Endpoints `json:"endpoints"`
	Profiles  Profiles  `json:"profiles"`
	Pipes     []Pipe    `json:"pipes"`
}

type Endpoint struct {
	Type          string    `json:"type"`
	UID           UID       `json:"uid"`
	Configuration ConfigMap `json:"configuration,omitempty"`
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
	Version string        `json:"version"`
	Script  engine.Script `json:"script"`
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

type ConfigMap map[string]interface{}

func (c ConfigMap) String(key string) (bool, string) {

	v, f := c[key]

	if !f {
		return false, ""
	}
	return true, v.(string)
}

func (c ConfigMap) MString(key string) string {

	found, v := c.String(key)
	if !found {
		panic(fmt.Errorf("value with key %s not found", key))
	}
	return v
}