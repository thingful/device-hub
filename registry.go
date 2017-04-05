// Copyright © 2017 thingful

package hub

import (
	"fmt"
	"sync"

	"github.com/thingful/device-hub/utils"
)

var (
	builders     = map[string]lazy{}
	buildersLock = sync.RWMutex{}
)

type endpointBuilder func(config utils.TypedMap) (Endpoint, error)
type listenerBuilder func(config utils.TypedMap) (Listener, error)

type builderFunc func(config utils.TypedMap) (interface{}, error)

type lazy struct {
	builder builderFunc
	built   interface{}
}

// RegisterEndpoint will store the builder with the correct name
func RegisterEndpoint(typez string, builder endpointBuilder) {

	buildersLock.Lock()
	defer buildersLock.Unlock()

	builders[typez] = lazy{
		builder: func(config utils.TypedMap) (interface{}, error) {
			i, err := builder(config)
			return i, err
		},
	}

}

// RegisterListener will store the builder with the correct name
func RegisterListener(typez string, builder listenerBuilder) {
	buildersLock.Lock()
	defer buildersLock.Unlock()

	builders[typez] = lazy{
		builder: func(config utils.TypedMap) (interface{}, error) {
			i, err := builder(config)
			return i, err
		},
	}
}

// EndpointByName returns or creates an Endpoint of specified type
func EndpointByName(uid, typez string, conf utils.TypedMap) (Endpoint, error) {

	f, err := genericByName(uid, typez, conf)

	if err != nil {
		return nil, err
	}

	e, ok := f.(Endpoint)

	if !ok {
		return nil, fmt.Errorf("builder registered with uid %s, type %s does not implement the Endpoint interface", uid, typez)
	}

	return e, nil
}

// ListenerByName returns or creates a Listener of specified type
func ListenerByName(uid, typez string, conf utils.TypedMap) (Listener, error) {

	f, err := genericByName(uid, typez, conf)

	if err != nil {
		return nil, err
	}

	l, ok := f.(Listener)

	if !ok {
		return nil, fmt.Errorf("builder registered with uid %s, type %s does not implement the Listener interface", uid, typez)
	}

	return l, nil
}

// genericByName exists instead of language support for generics!
func genericByName(uid, typez string, conf utils.TypedMap) (interface{}, error) {

	buildersLock.Lock()
	defer buildersLock.Unlock()

	// try and find by uid
	e, found := builders[uid]

	if !found {

		// try and find the builder for the type
		e, found = builders[typez]

		if !found {
			return nil, fmt.Errorf("builder with uid : %s, type : %s not found", uid, typez)
		}
	}

	// if already created return it
	if e.built != nil {
		return e.built, nil
	}

	// make a new endpoint storing it against the uid
	endpoint, err := e.builder(conf)

	builders[uid] = lazy{
		builder: e.builder,
		built:   endpoint,
	}

	return endpoint, err

}
