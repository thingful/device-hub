// Copyright © 2017 thingful

package mqtt

import (
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

	return mqtt.NewClient(opts)
}
