// Copyright © 2017 thingful

package listener

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"

	"github.com/thingful/device-hub/registry"
)

func Register(r *registry.Registry) {

	mqtt_bindingAddress := describe.Parameter{
		Name:        "mqtt-broker-address",
		Type:        describe.Url,
		Required:    true,
		Description: "address to connect to",
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

	r.RegisterListener("mqtt",

		func(config describe.Values) (hub.Listener, error) {

			brokerAddress := config.MustString(mqtt_bindingAddress.Name)
			clientID := fmt.Sprintf("device-hub-%s", hub.SourceVersion)
			username, ufound := config.String(mqtt_username.Name)
			password, pfound := config.String(mqtt_password.Name)

			opts := mqtt.NewClientOptions()

			if ufound {
				opts.SetUsername(username)
			}

			if pfound {
				opts.SetPassword(password)
			}

			opts.AddBroker(brokerAddress)
			opts.SetClientID(clientID)

			opts.SetKeepAlive(2 * time.Second)
			opts.SetPingTimeout(10 * time.Second)
			opts.SetAutoReconnect(true)

			// Panic on connection lost until
			// https://github.com/thingful/device-hub/issues/27
			// is resolved
			opts.OnConnectionLost = func(client mqtt.Client, err error) {
				log.Panic("mqtt broker disconnected", err)
			}

			opts.OnConnect = func(mqtt.Client) {
				log.Print("mqtt broker connected")
			}

			client := mqtt.NewClient(opts)

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
		Type:        describe.String,
		Required:    true,
		Description: "address to bind to",
		Examples:    []string{"0.0.0.0:9090", "*:9090", ":8000"},
	}

	r.RegisterListener("http",
		func(config describe.Values) (hub.Listener, error) {

			binding := config.MustString(http_bindingAddress.Name)
			return newHTTPListener(binding)

		},
		describe.Parameters{
			http_bindingAddress,
		},
	)
}
