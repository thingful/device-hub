// Copyright © 2017 thingful

package listener

import (
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	hub "github.com/thingful/device-hub"
)

func newMQTTListener(client mqtt.Client) (*mqttlistener, error) {

	if !client.IsConnected() {
		return nil, errors.New("mqtt client is not connected")
	}

	return &mqttlistener{
		client: client,
	}, nil
}

type mqttlistener struct {
	client mqtt.Client
}

func (m *mqttlistener) NewChannel(topic string) (hub.Channel, error) {

	if topic == "" {
		return nil, errors.New("mqtt topic is empty string")
	}

	errors := make(chan error)
	out := make(chan hub.Message)

	handler := func(client mqtt.Client, msg mqtt.Message) {
		input := newHubMessage(msg.Payload(), "MQTT", msg.Topic())
		out <- input
	}

	if token := m.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return defaultChannel{out: out, errors: errors, close: func() error {
		token := m.client.Unsubscribe(topic)
		token.Wait()
		return token.Error()
	}}, nil
}

func (m *mqttlistener) Close() error {
	// TODO : set a sensible timeout
	m.client.Disconnect(1000)
	return nil
}
