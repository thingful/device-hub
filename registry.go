// Copyright © 2017 thingful

package hub

import (
	"fmt"
	"sync"

	"github.com/thingful/device-hub/utils"
)

var (
	endpointBuilders     = map[string]lazyEndpoint{}
	endpointBuildersLock = sync.RWMutex{}
)

type endpointBuilder func(config utils.TypedMap) (Endpoint, error)

type lazyEndpoint struct {
	builder  endpointBuilder
	endpoint Endpoint
}

// RegisterEndpoint will store the builder with the correct name
func RegisterEndpoint(typez string, builder endpointBuilder) {

	endpointBuildersLock.Lock()
	defer endpointBuildersLock.Unlock()

	endpointBuilders[typez] = lazyEndpoint{
		builder: builder,
	}

}

// EndpointByName retrieves the builder by name, and initilising with the configuration
func EndpointByName(uid, typez string, conf utils.TypedMap) (Endpoint, error) {

	endpointBuildersLock.Lock()
	defer endpointBuildersLock.Unlock()

	// try and find by uid
	e, found := endpointBuilders[uid]

	if !found {

		// try and find the builder for the type
		e, found = endpointBuilders[typez]

		if !found {
			return nil, fmt.Errorf("endpoint of uid : %s, type : %s  not found", uid, typez)
		}
	}

	// if already created return it
	if e.endpoint != nil {
		return e.endpoint, nil
	}

	// make a new endpoint storing it against the uid
	endpoint, err := e.builder(conf)

	endpointBuilders[uid] = lazyEndpoint{
		builder:  e.builder,
		endpoint: endpoint,
	}

	return endpoint, err

}
