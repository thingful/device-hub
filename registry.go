// Copyright © 2017 thingful

package hub

import (
	"fmt"
	"sync"

	"github.com/thingful/device-hub/describe"
)

var (
	endpoints          = map[string]lazy{}
	endpointParameters = map[string][]describe.Parameter{}
	endpointsLock      = sync.RWMutex{}

	listeners          = map[string]lazy{}
	listenerParameters = map[string][]describe.Parameter{}
	listenersLock      = sync.RWMutex{}
)

type endpointBuilder func(config describe.Values) (Endpoint, error)
type listenerBuilder func(config describe.Values) (Listener, error)

type builderFunc func(config describe.Values) (interface{}, error)

type lazy struct {
	builder builderFunc
	built   interface{}
}

// RegisterEndpoint will store the builder with the correct name
func RegisterEndpoint(typez string, builder endpointBuilder, params []describe.Parameter) {

	endpointsLock.Lock()
	defer endpointsLock.Unlock()

	endpoints[typez] = lazy{
		builder: func(config describe.Values) (interface{}, error) {
			i, err := builder(config)
			return i, err
		},
	}

	endpointParameters[typez] = params

}

// IsEndpointRegistered confirms if the endpoint has been registered
func IsEndpointRegistered(typez string) bool {

	endpointsLock.Lock()
	defer endpointsLock.Unlock()

	_, found := endpoints[typez]
	return found
}

// RegisterListener will store the builder with the correct name
func RegisterListener(typez string, builder listenerBuilder, params []describe.Parameter) {

	listenersLock.Lock()
	defer listenersLock.Unlock()

	listeners[typez] = lazy{
		builder: func(config describe.Values) (interface{}, error) {
			i, err := builder(config)
			return i, err
		},
	}

	listenerParameters[typez] = params
}

// IsListenerRegistered confirms if the listener has been registered
func IsListenerRegistered(typez string) bool {

	listenersLock.Lock()
	defer listenersLock.Unlock()

	_, found := listeners[typez]
	return found
}

func DescribeListener(typez string) ([]describe.Parameter, error) {

	// TODO: deal with not being found
	return listenerParameters[typez], nil
}

// EndpointByName returns or creates an Endpoint of specified type
func EndpointByName(uid, typez string, conf map[string]string) (Endpoint, error) {

	endpointsLock.Lock()
	defer endpointsLock.Unlock()

	parameters, found := endpointParameters[typez]

	if !found {
		return nil, fmt.Errorf("parameters for type %s not found", typez)
	}

	values, err := describe.NewValues(conf, parameters)

	if err != nil {
		return nil, err
	}

	f, err := genericByName(endpoints, uid, typez, values)

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
func ListenerByName(uid, typez string, conf map[string]string) (Listener, error) {

	listenersLock.Lock()
	defer listenersLock.Unlock()

	parameters, found := listenerParameters[typez]

	if !found {
		return nil, fmt.Errorf("parameters for type %s not found", typez)
	}

	values, err := describe.NewValues(conf, parameters)

	if err != nil {
		return nil, err
	}

	f, err := genericByName(listeners, uid, typez, values)

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
func genericByName(builders map[string]lazy, uid, typez string, conf describe.Values) (interface{}, error) {

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
