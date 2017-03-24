package pipe

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	hub "github.com/thingful/device-hub"
)

var (
	r           *router
	router_lock sync.Mutex
)

func DefaultRouter() *router {

	router_lock.Lock()
	defer router_lock.Unlock()

	if r != nil {
		return r
	}

	r = &router{routes: map[string]Channel{}}
	return r

}

func StartDefaultHTTPListener(ctx context.Context, router *router, binding string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		ok, channel := router.Match(path)

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			channel.Errors() <- err
		}

		input := hub.Input{
			Payload: body,
		}

		channel.Out() <- input
		w.WriteHeader(http.StatusAccepted)
	})

	go func() {
		log.Fatal(http.ListenAndServe(binding, nil))

	}()
}

func NewHTTPChannel(uri string, router *router) Channel {

	errors := make(chan error)
	out := make(chan hub.Input)

	channel := defaultChannel{out: out, errors: errors}

	router.register(uri, channel)
	return channel
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
