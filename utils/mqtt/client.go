// Copyright © 2017 thingful

package mqtt

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// DefaultMQTTClient returns an mqtt.Client with some sensible defaults
func DefaultMQTTClient(brokerAddress, clientID string) mqtt.Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerAddress)
	opts.SetClientID(clientID)

	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Panic("mqtt broker disconnected", err)
	}

	opts.OnConnect = func(mqtt.Client) {
		log.Print("mqtt broker connected")
	}

	return mqtt.NewClient(opts)
}
