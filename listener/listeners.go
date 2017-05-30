// Copyright © 2017 thingful

package listener

import (
	"fmt"

	hub "github.com/thingful/device-hub"
	d "github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/utils/mqtt"
)

func init() {

	hub.RegisterListener("mqtt",

		func(config d.Values) (hub.Listener, error) {

			brokerAddress := config.MustString("mqtt-broker-address")

			clientName := fmt.Sprintf("device-hub-%s", hub.SourceVersion)

			client := mqtt.DefaultMQTTClient(brokerAddress, clientName)

			// TODO : set sensible wait time
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}

			return newMQTTListener(client)

		},
		d.Parameters{
			d.Parameter{
				Name:        "mqtt-broker-address",
				Type:        d.Url,
				Required:    true,
				Description: "address to bind to",
				Examples:    []string{"tcp://0.0.0.0:1883"}},
		},
	)

	hub.RegisterListener("http",
		func(config d.Values) (hub.Listener, error) {

			binding := config.MustString("http-binding-address")
			return newHTTPListener(binding)

		},
		d.Parameters{
			d.Parameter{
				Name:        "http-binding-address",
				Type:        d.Url,
				Required:    true,
				Description: "address to bind to",
				Examples:    []string{"0.0.0.0:9090", "*:9090"}},
		},
	)
}
