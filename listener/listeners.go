// Copyright © 2017 thingful

package listener

import (
	"fmt"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/utils/mqtt"
)

func init() {

	mqtt_bindingAddress := describe.Parameter{
		Name:        "mqtt-broker-address",
		Type:        describe.Url,
		Required:    true,
		Description: "address to bind to",
		Examples:    []string{"tcp://0.0.0.0:1883"},
	}

	mqtt_username := describe.Parameter{
		Name:        "mqtt-username",
		Type:        describe.String,
		Required:    false,
		Description: "user name for mqtt server",
	}

	mqtt_password := describe.Parameter{
		Name:        "mqtt-password",
		Type:        describe.String,
		Required:    false,
		Description: "user password for mqtt server",
	}

	hub.RegisterListener("mqtt",

		func(config describe.Values) (hub.Listener, error) {

			brokerAddress := config.MustString(mqtt_bindingAddress.Name)

			clientName := fmt.Sprintf("device-hub-%s", hub.SourceVersion)

			client := mqtt.DefaultMQTTClient(brokerAddress, clientName)

			// TODO : set sensible wait time
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				return nil, token.Error()
			}

			return newMQTTListener(client)

		},
		describe.Parameters{
			mqtt_bindingAddress,
			mqtt_username,
			mqtt_password,
		},
	)

	http_bindingAddress := describe.Parameter{
		Name:        "http-binding-address",
		Type:        describe.Url,
		Required:    true,
		Description: "address to bind to",
		Examples:    []string{"tcp://0.0.0.0:9090", "tcp://*:9090"},
	}

	hub.RegisterListener("http",
		func(config describe.Values) (hub.Listener, error) {

			binding := config.MustString(http_bindingAddress.Name)
			return newHTTPListener(binding)

		},
		describe.Parameters{
			http_bindingAddress,
		},
	)
}
