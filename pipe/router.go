package pipe

func DefaultRouter() *router {
	return &router{routes: map[string]Channel{}}
}

// TODO : add regex routing
// TODO : add header routing
type router struct {
	routes map[string]Channel
}

func (r *router) register(uri string, channel Channel) {
	// TODO : is overwriting good enough
	// TODO : should this be locked?
	r.routes[uri] = channel
}

func (r *router) Match(uri string) (bool, Channel) {
	c, ok := r.routes[uri]
	return ok, c
}
