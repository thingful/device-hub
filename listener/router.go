// Copyright © 2017 thingful
package listener

import hub "github.com/thingful/device-hub"

func DefaultRouter() *router {
	return &router{routes: map[string]hub.Channel{}}
}

// TODO : add regex routing
// TODO : add header routing
type router struct {
	routes map[string]hub.Channel
}

func (r *router) register(uri string, channel hub.Channel) {
	// TODO : is overwriting good enough
	// TODO : should this be locked?
	r.routes[uri] = channel
}

func (r *router) Match(uri string) (bool, hub.Channel) {
	c, ok := r.routes[uri]
	return ok, c
}
