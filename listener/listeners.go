// Copyright © 2017 thingful

package listener

import (
	"fmt"

	hub "github.com/thingful/device-hub"
	d "github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/utils"
	"github.com/thingful/device-hub/utils/mqtt"
)

func init() {

	hub.RegisterListener("mqtt",

		func(config utils.TypedMap) (hub.Listener, error) {

			clientName := fmt.Sprintf("device-hub-%s", hub.SourceVersion)

			brokerAddress := config.MString("mqtt-broker-address")

			client := mqtt.DefaultMQTTClient(brokerAddress, clientName)

			// TODO : set sensible wait time
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}

			return newMQTTListener(client)

		},
		[]d.Parameter{
			d.Parameter{Name: "mqtt-broker-address",
				Type:        d.Url,
				Required:    true,
				Description: "address to bind to",
				Examples:    []string{"tcp://0.0.0.0:1883"}},
		},
	)

	hub.RegisterListener("http",
		func(config utils.TypedMap) (hub.Listener, error) {

			binding := config.MString("http-binding-address")
			return newHTTPListener(binding)

		},
		[]d.Parameter{
			d.Parameter{Name: "http-binding-address",
				Type:        d.Url,
				Required:    true,
				Description: "address to bind to",
				Examples:    []string{"0.0.0.0:9090", "*:9090"}},
		},
	)
}
