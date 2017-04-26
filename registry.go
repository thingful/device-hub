// Copyright © 2017 thingful

package hub

import (
	"fmt"
	"sync"

	"github.com/thingful/device-hub/utils"
)

var (
	endpoints     = map[string]lazy{}
	endpointsLock = sync.RWMutex{}

	listeners     = map[string]lazy{}
	listenersLock = sync.RWMutex{}
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

	endpointsLock.Lock()
	defer endpointsLock.Unlock()

	endpoints[typez] = lazy{
		builder: func(config utils.TypedMap) (interface{}, error) {
			i, err := builder(config)
			return i, err
		},
	}

}

// IsEndpointRegistered confirms if the endpoint has been registered
func IsEndpointRegistered(typez string) bool {

	endpointsLock.Lock()
	defer endpointsLock.Unlock()

	_, found := endpoints[typez]
	return found
}

// RegisterListener will store the builder with the correct name
func RegisterListener(typez string, builder listenerBuilder) {

	listenersLock.Lock()
	defer listenersLock.Unlock()

	listeners[typez] = lazy{
		builder: func(config utils.TypedMap) (interface{}, error) {
			i, err := builder(config)
			return i, err
		},
	}
}

// IsListenerRegistered confirms if the listener has been registered
func IsListenerRegistered(typez string) bool {

	listenersLock.Lock()
	defer listenersLock.Unlock()

	_, found := listeners[typez]
	return found
}

// EndpointByName returns or creates an Endpoint of specified type
func EndpointByName(uid, typez string, conf utils.TypedMap) (Endpoint, error) {

	endpointsLock.Lock()
	defer endpointsLock.Unlock()

	f, err := genericByName(endpoints, uid, typez, conf)

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

	listenersLock.Lock()
	defer listenersLock.Unlock()

	f, err := genericByName(listeners, uid, typez, conf)

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
func genericByName(builders map[string]lazy, uid, typez string, conf utils.TypedMap) (interface{}, error) {

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
