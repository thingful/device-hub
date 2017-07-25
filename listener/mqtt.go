// Copyright © 2017 thingful

package listener

import (
	"errors"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	hub "github.com/thingful/device-hub"
)

func newMQTTListener(client mqtt.Client) (*mqttlistener, error) {

	if !client.IsConnected() {
		return nil, errors.New("mqtt client is not connected")
	}

	return &mqttlistener{
		client:        client,
		subscriptions: map[string]defaultChannel{},
	}, nil
}

type mqttlistener struct {
	client mqtt.Client

	lock          sync.RWMutex
	subscriptions map[string]defaultChannel
}

func (m *mqttlistener) NewChannel(topic string) (hub.Channel, error) {

	if topic == "" {
		return nil, errors.New("mqtt topic is empty string")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	_, found := m.subscriptions[topic]

	if found {
		return nil, errors.New("unable to start subscription for existing topic")
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

	channel := defaultChannel{out: out, errors: errors, close: func() error {
		token := m.client.Unsubscribe(topic)
		token.Wait()
		return token.Error()
	}}

	m.subscriptions[topic] = channel

	return channel, nil
}

func (m *mqttlistener) Close() error {
	// TODO : set a sensible timeout
	m.client.Disconnect(1000)
	return nil
}

func (m *mqttlistener) RestartSubscriptions() error {

	m.lock.Lock()
	defer m.lock.Unlock()

	for topic, s := range m.subscriptions {

		handler := func(client mqtt.Client, msg mqtt.Message) {
			input := newHubMessage(msg.Payload(), "MQTT", msg.Topic())
			s.out <- input
		}

		if token := m.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	return nil
}
