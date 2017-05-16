// Copyright © 2017 thingful

package endpoint

import (
	"errors"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/utils"
)

func init() {

	hub.RegisterEndpoint("stdout", func(config utils.TypedMap) (hub.Endpoint, error) {

		prettyPrint := config.DBool("pretty-print", false)

		return stdout{
			prettyPrint: prettyPrint,
		}, nil
	})

	hub.RegisterEndpoint("http", func(config utils.TypedMap) (hub.Endpoint, error) {

		clientTimeOut := config.DInt("http-client-timeout-ms", 1000)

		found, url := config.String("http-url")

		if !found {
			return nil, errors.New("http endpoint needs 'http-url' set")
		}

		return NewHTTPEndpoint(url, clientTimeOut), nil
	})

}
