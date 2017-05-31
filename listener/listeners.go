// Copyright © 2017 thingful

package listener

import (
	"fmt"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/utils/mqtt"
)

func init() {

	hub.RegisterListener("mqtt",

		func(config describe.Values) (hub.Listener, error) {

			brokerAddress := config.MustString("mqtt-broker-address")

			clientName := fmt.Sprintf("device-hub-%s", hub.SourceVersion)

			client := mqtt.DefaultMQTTClient(brokerAddress, clientName)

			// TODO : set sensible wait time
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}

			return newMQTTListener(client)

		},
		describe.Parameters{
			describe.Parameter{
				Name:        "mqtt-broker-address",
				Type:        describe.Url,
				Required:    true,
				Description: "address to bind to",
				Examples:    []string{"tcp://0.0.0.0:1883"}},
		},
	)

	hub.RegisterListener("http",
		func(config describe.Values) (hub.Listener, error) {

			binding := config.MustString("http-binding-address")
			return newHTTPListener(binding)

		},
		describe.Parameters{
			describe.Parameter{
				Name:        "http-binding-address",
				Type:        describe.Url,
				Required:    true,
				Description: "address to bind to",
				Examples:    []string{"tcp://0.0.0.0:9090", "tcp://*:9090"}},
		},
	)
}
