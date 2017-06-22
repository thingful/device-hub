// Copyright © 2017 thingful

package registry

import (
	"fmt"
	"sync"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"
)

type Registry struct {
	endpoints          map[string]lazy
	endpointParameters map[string][]describe.Parameter
	endpointsLock      sync.RWMutex

	listeners          map[string]lazy
	listenerParameters map[string][]describe.Parameter
	listenersLock      sync.RWMutex
}

var (
	Default = New()
)

func New() *Registry {

	return &Registry{
		endpoints:          map[string]lazy{},
		endpointParameters: map[string][]describe.Parameter{},
		endpointsLock:      sync.RWMutex{},

		listeners:          map[string]lazy{},
		listenerParameters: map[string][]describe.Parameter{},
		listenersLock:      sync.RWMutex{},
	}
}

type endpointBuilder func(config describe.Values) (hub.Endpoint, error)
type listenerBuilder func(config describe.Values) (hub.Listener, error)

type builderFunc func(config describe.Values) (interface{}, error)

type lazy struct {
	builder builderFunc
	built   interface{}
}

// RegisterEndpoint will store the builder with the correct name
func (r *Registry) RegisterEndpoint(typez string, builder endpointBuilder, params describe.Parameters) {

	if len(params) == 0 {
		panic("endpoint registered without any parameters")
	}

	r.endpointsLock.Lock()
	defer r.endpointsLock.Unlock()

	r.endpoints[typez] = lazy{
		builder: func(config describe.Values) (interface{}, error) {
			i, err := builder(config)
			return i, err
		},
	}

	r.endpointParameters[typez] = params

}

// IsEndpointRegistered confirms if the endpoint has been registered
func (r *Registry) IsEndpointRegistered(typez string) bool {

	r.endpointsLock.Lock()
	_, found := r.endpoints[typez]
	r.endpointsLock.Unlock()

	return found
}

// DescribeEndpoint returns a collection of Parameter describing its configuration
func (r *Registry) DescribeEndpoint(typez string) (describe.Parameters, error) {

	r.endpointsLock.Lock()
	params, found := r.endpointParameters[typez]
	r.endpointsLock.Unlock()

	if !found {
		return nil, fmt.Errorf("no parameters found for endpoint : %s", typez)
	}

	return params, nil
}

// RegisterListener will store the builder with the correct name
func (r *Registry) RegisterListener(typez string, builder listenerBuilder, params describe.Parameters) {

	if len(params) == 0 {
		panic("listener registered without any parameters")
	}

	r.listenersLock.Lock()
	defer r.listenersLock.Unlock()

	r.listeners[typez] = lazy{
		builder: func(config describe.Values) (interface{}, error) {
			i, err := builder(config)
			return i, err
		},
	}

	r.listenerParameters[typez] = params
}

// IsListenerRegistered confirms if the listener has been registered
func (r *Registry) IsListenerRegistered(typez string) bool {

	r.listenersLock.Lock()
	_, found := r.listeners[typez]
	r.listenersLock.Unlock()

	return found
}

// DescribeListener returns a collection of Parameter describing its configuration
func (r *Registry) DescribeListener(typez string) (describe.Parameters, error) {

	r.listenersLock.Lock()
	params, found := r.listenerParameters[typez]
	r.listenersLock.Unlock()

	if !found {
		return nil, fmt.Errorf("no parameters found for listener : %s", typez)
	}

	return params, nil
}

// EndpointByName returns or creates an Endpoint of specified type
func (r *Registry) EndpointByName(uid, typez string, conf map[string]string) (hub.Endpoint, error) {

	r.endpointsLock.Lock()
	parameters, found := r.endpointParameters[typez]
	r.endpointsLock.Unlock()

	if !found {
		return nil, fmt.Errorf("parameters for type %s not found", typez)
	}

	values, err := describe.NewValues(conf, parameters)

	if err != nil {
		return nil, err
	}

	f, err := genericByName(r.endpoints, uid, typez, values)

	if err != nil {
		return nil, err
	}

	e, ok := f.(hub.Endpoint)

	if !ok {
		return nil, fmt.Errorf("builder registered with uid %s, type %s does not implement the Endpoint interface", uid, typez)
	}

	return e, nil
}

// ListenerByName returns or creates a Listener of specified type
func (r *Registry) ListenerByName(uid, typez string, conf map[string]string) (hub.Listener, error) {

	r.listenersLock.Lock()
	parameters, found := r.listenerParameters[typez]
	r.listenersLock.Unlock()

	if !found {
		return nil, fmt.Errorf("parameters for type %s not found", typez)
	}

	values, err := describe.NewValues(conf, parameters)

	if err != nil {
		return nil, err
	}

	f, err := genericByName(r.listeners, uid, typez, values)

	if err != nil {
		return nil, err
	}

	l, ok := f.(hub.Listener)

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
