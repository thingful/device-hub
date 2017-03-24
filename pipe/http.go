package pipe

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	hub "github.com/thingful/device-hub"
)

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
