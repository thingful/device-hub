package config

import "github.com/thingful/device-hub/engine"

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
	Pipes []pipe `json:"pipes"`
}

// Endpoint contains a generic collection of configuration details
type Endpoint struct {
	Type          string    `json:"type"`
	UID           UID       `json:"uid"`
	Configuration configMap `json:"configuration,omitempty"`
}

type endpoints []Endpoint

// FindByUID looks for an endpoint by UID returning true if found, false if missing
func (e endpoints) FindByUID(uid UID) (bool, Endpoint) {

	for _, endpoint := range e {
		if endpoint.UID == uid {
			return true, endpoint
		}
	}
	return false, Endpoint{}
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
