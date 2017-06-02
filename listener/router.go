// Copyright © 2017 thingful
package listener

import (
	"sync"

	hub "github.com/thingful/device-hub"
)

func DefaultRouter() *router {
	return &router{routes: map[string]hub.Channel{}}
}

// TODO : add regex routing
// TODO : add header routing
type router struct {
	routes map[string]hub.Channel
	sync.Mutex
}

func (r *router) register(uri string, channel hub.Channel) {
	// TODO : is overwriting good enough
	r.Lock()
	r.routes[uri] = channel
	r.Unlock()
}

func (r *router) delete(uri string) error {
	delete(r.routes, uri)
	return nil
}

func (r *router) Match(uri string) (bool, hub.Channel) {
	r.Lock()
	c, ok := r.routes[uri]
	r.Unlock()

	return ok, c
}
