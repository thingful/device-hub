// Copyright © 2017 thingful

package listener

import (
	"fmt"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/utils"
)

func init() {

	hub.RegisterListener("mqtt", func(config utils.TypedMap) (hub.Listener, error) {

		clientName := fmt.Sprintf("device-hub-%s", hub.SourceVersion)

		brokerAddress := config.MString("mqtt-broker-address")

		options := defaultMQTTOptions(brokerAddress, clientName)
		client := defaultMQTTClient(options)

		// TODO : set sensible wait time
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return nil, token.Error()
		}

		return newMQTTListener(client)

	})

	hub.RegisterListener("http", func(config utils.TypedMap) (hub.Listener, error) {

		binding := config.MString("http-binding-address")
		return newHTTPListener(binding)

	})
}
