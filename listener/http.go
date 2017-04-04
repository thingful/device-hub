// Copyright © 2017 thingful

package listener

import (
	"io/ioutil"
	"log"
	"net/http"

	hub "github.com/thingful/device-hub"
)

func NewHTTPListener(binding string) (*httpListener, error) {

	router := DefaultRouter()

	startDefaultHTTPListener(router, binding)

	return &httpListener{
		router: router,
	}, nil
}

type httpListener struct {
	router *router
}

func (h *httpListener) NewChannel(uri string) (hub.Channel, error) {

	errors := make(chan error)
	out := make(chan hub.Message)

	channel := defaultChannel{out: out, errors: errors}

	h.router.register(uri, channel)
	return channel, nil
}

func (h *httpListener) Close() error {
	return nil
}

func startDefaultHTTPListener(router *router, binding string) {

	http.HandleFunc("/", rootHandler(router))

	// TODO : shutdown nicely
	go func() {
		log.Fatal(http.ListenAndServe(binding, nil))

	}()
}

func rootHandler(router *router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		path := r.URL.Path

		ok, channel := router.Match(path)

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			channel.Errors() <- err
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		input := newHubMessage(body, "HTTP", r.URL.Path)

		channel.Out() <- input
		w.WriteHeader(http.StatusAccepted)
	}
}
