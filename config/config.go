// Copyright © 2017 thingful

package config

import (
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/utils"
)

// UID is a unique string-like value
type UID string

// Configuration is the root object for a device-hub configuration
type Configuration struct {
	// Listeners are the entry point to the system and understand protocols e.g. http, coap, mqtt
	Listeners endpoints `json:"listeners"`
	// Endpoints are the destinations for processed data e.g. another web service, a message queue
	Endpoints endpoints `json:"endpoints"`
	// Profiles are device specific, versioned processors of data and act inbetween the listeners and the endpoints
	Profiles profiles `json:"profiles"`
	// Pipes wire the above into a shallow series of steps
	Pipes pipes `json:"pipes"`
}

// Endpoint contains a generic collection of configuration details
type Endpoint struct {
	Type          string         `json:"type"`
	UID           UID            `json:"uid"`
	Configuration utils.TypedMap `json:"configuration,omitempty"`
}

type endpoints []Endpoint

// FindByUID looks for an endpoint by UID returning true if found, false if missing
func (e endpoints) FindByUID(uid UID) (bool, Endpoint) {

	found := false
	endpoint := Endpoint{}

	e.mapper(func(a Endpoint) {
		if a.UID == uid {
			found = true
			endpoint = a
		}
	})
	return found, endpoint
}

// endpoint mapper is an action specialised for Endpoints
type endpointMapper func(e Endpoint)

func (e endpoints) mapper(f endpointMapper) {
	for _, endpoint := range e {
		f(endpoint)
	}
}

// Profile is a device profile
type Profile struct {
	UID         UID    `json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// TODO : make this a semantic triple
	Version string        `json:"version"`
	Script  engine.Script `json:"script"`
}

type profiles []Profile

// FindByUID looks for an profile by UID returning true if found, false if missing
func (p profiles) FindByUID(uid UID) (bool, Profile) {

	for _, profile := range p {
		if profile.UID == uid {
			return true, profile
		}
	}
	return false, Profile{}
}

type pipe struct {
	Uri       string `json:"uri"`
	Profile   UID    `json:"profile"`
	Listener  UID    `json:"listener"`
	Endpoints []UID  `json:"endpoints"`
}

type pipes []pipe

type pipeMapper func(p pipe)

func (p pipes) mapper(f pipeMapper) {
	for _, pipe := range p {
		f(pipe)
	}
}
