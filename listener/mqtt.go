// Copyright © 2017 thingful

package listener

import (
	"errors"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	hub "github.com/thingful/device-hub"
)

const (
	// mqttClientDisconnectTimeoutInMs is the disconnect
	// wait time to allow other operations to succeed
	mqttClientDisconnectTimeoutInMs = 1000
)

func newMQTTListener(options *mqtt.ClientOptions) (*mqttlistener, error) {

	listener := &mqttlistener{
		options:       options,
		subscriptions: map[string]defaultChannel{},
	}

	options.OnConnect = func(mqtt.Client) {
		log.Print("mqtt broker connected")
	}

	options.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Println("mqtt broker disconnected", err)

		err2 := listener.restartSubscriptions()

		if err2 != nil {
			log.Panicf("error restarting subscriptions - %v", err2)
		}
	}

	return listener, nil
}

type mqttlistener struct {

	// keep a copy of of the mqtt.ClientOptions as the mqtt.Client
	// is not too keen on being stopped and started
	options *mqtt.ClientOptions

	// connection_lock tracks the 'genesis' connection
	connection_lock sync.RWMutex

	client mqtt.Client

	// subscriptions are tracked so that we can resurrect state
	// on disconnections
	subscription_lock sync.RWMutex
	subscriptions     map[string]defaultChannel
}

func (m *mqttlistener) NewChannel(topic string) (hub.Channel, error) {

	if topic == "" {
		return nil, errors.New("mqtt topic is empty string")
	}

	m.subscription_lock.Lock()
	defer m.subscription_lock.Unlock()

	_, found := m.subscriptions[topic]

	if found {
		return nil, errors.New("unable to start subscription for existing topic")
	}

	// ensure we have a 'genisis' connection
	if err := m.ensureGenisisConnectionMade(); err != nil {
		return nil, err
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
		return m.closeDownChannel(topic)
	}}

	m.subscriptions[topic] = channel
	return channel, nil
}

func (m *mqttlistener) Close() error {
	m.client.Disconnect(mqttClientDisconnectTimeoutInMs)
	return nil
}

// restartSubscriptions is called to on reconnection and attempts to reinstate the
// exists subscriptions
func (m *mqttlistener) restartSubscriptions() error {

	m.subscription_lock.Lock()
	defer m.subscription_lock.Unlock()

	for topic, s := range m.subscriptions {

		log.Println("attempting to reconnect existing subscription - ", topic)
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

// closeDownChannel removes the mqtt subscription,
// closing the connection if no more subscriptions
func (m *mqttlistener) closeDownChannel(topic string) error {

	// tear down mqtt subscription
	token := m.client.Unsubscribe(topic)
	token.Wait()

	err := token.Error()

	if err != nil {
		return err
	}

	// clean up subscriptions
	m.subscription_lock.Lock()
	defer m.subscription_lock.Unlock()

	delete(m.subscriptions, topic)

	// if no subscriptions left disconnect the client
	if len(m.subscriptions) == 0 {

		m.connection_lock.Lock()
		defer m.connection_lock.Unlock()

		m.client.Disconnect(mqttClientDisconnectTimeoutInMs)
		m.client = nil
	}

	return nil

}

// ensureGenisisConnectionMade will connect the mqtt client on the first channel added.
// Once connected the mqtt.Client is clever enough to reconnect and the RestartSubscriptions
// method will be called on reconnects to resurrect the subscriptions
func (m *mqttlistener) ensureGenisisConnectionMade() error {

	m.connection_lock.Lock()
	defer m.connection_lock.Unlock()

	if m.client != nil {
		return nil
	}

	m.client = mqtt.NewClient(m.options)

	// TODO : set sensible wait time
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil

}
