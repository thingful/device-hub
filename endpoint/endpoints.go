// Copyright © 2017 thingful

package endpoint

import (
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
}
