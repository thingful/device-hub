// Copyright © 2017 thingful

package endpoint

import (
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/registry"
)

func Register(r *registry.Registry) {

	r.RegisterEndpoint("stdout", func(config describe.Values) (hub.Endpoint, error) {

		prettyPrint := config.BoolWithDefault("pretty-print", false)

		return stdout{
			prettyPrint: prettyPrint,
		}, nil
	},
		describe.Parameters{
			describe.Parameter{
				Name:        "pretty-print",
				Type:        describe.Bool,
				Required:    false,
				Description: "pretty print the output",
				Default:     false,
			},
		})

	httpClientTimeoutMS := int32(1000)

	r.RegisterEndpoint("http", func(config describe.Values) (hub.Endpoint, error) {

		clientTimeOut := config.Int32WithDefault("http-client-timeout-ms", httpClientTimeoutMS)
		url := config.MustString("http-url")

		return NewHTTPEndpoint(url, clientTimeOut), nil
	},
		describe.Parameters{
			describe.Parameter{
				Name:        "http-url",
				Type:        describe.Url,
				Required:    true,
				Description: "url to post data to",
			},
			describe.Parameter{
				Name:        "http-client-timeout-ms",
				Type:        describe.Int32,
				Required:    false,
				Default:     httpClientTimeoutMS,
				Description: "http client time-out in ms",
			},
		})

}
