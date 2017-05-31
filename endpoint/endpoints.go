// Copyright © 2017 thingful

package endpoint

import (
	"fmt"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"
)

func init() {

	hub.RegisterEndpoint("stdout", func(config describe.Values) (hub.Endpoint, error) {

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
			},
		})

	httpClientTimeoutMS := int32(1000)

	hub.RegisterEndpoint("http", func(config describe.Values) (hub.Endpoint, error) {

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
				Description: fmt.Sprintf("http client time out in ms, defaults to %s", httpClientTimeoutMS),
			},
		})

}
