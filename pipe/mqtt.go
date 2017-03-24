package pipe

import (
	"errors"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	hub "github.com/thingful/device-hub"
)

const (
	TOPIC_NAME = "topic"
)

func DefaultMQTTOptions(brokerAddress, clientID string) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerAddress)
	opts.SetClientID(clientID)

	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	return opts
}

var (
	client      mqtt.Client
	client_lock sync.Mutex
)

func DefaultClient(options *mqtt.ClientOptions) mqtt.Client {

	client_lock.Lock()
	defer client_lock.Unlock()

	if client != nil {
		return client
	}

	client = mqtt.NewClient(options)
	return client

}

func FromMQTT(client mqtt.Client, topic string) (*mqttbroker, error) {

	if topic == "" {
		return nil, errors.New("mqtt topic is empty string")
	}

	errors := make(chan error)
	return &mqttbroker{
		topic:  topic,
		errors: errors,
		client: client,
	}, nil
}

type mqttbroker struct {
	topic           string
	errors          chan error
	client          mqtt.Client
	connection_lock sync.Mutex
}

func (m *mqttbroker) Channel() (Channel, error) {

	err := m.connect()

	if err != nil {
		return NoOpChannel{}, err
	}

	channel := make(chan hub.Input)

	handler := func(client mqtt.Client, msg mqtt.Message) {
		input := hub.Input{
			Payload: msg.Payload(),
			Metadata: map[string]interface{}{
				TOPIC_NAME: msg.Topic(),
			},
		}

		channel <- input

	}

	if token := m.client.Subscribe(m.topic, 0, handler); token.Wait() && token.Error() != nil {
		return NoOpChannel{}, token.Error()
	}

	return mqttChannel{out: channel, errors: m.errors}, nil
}

func (m *mqttbroker) connect() error {

	m.connection_lock.Lock()
	defer m.connection_lock.Unlock()

	if m.client.IsConnected() {
		return nil
	}

	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *mqttbroker) Close() error {
	return nil
}

type mqttChannel struct {
	errors chan error
	out    chan hub.Input
}

// Errors returns a channel of errors
func (m mqttChannel) Errors() chan error {
	return m.errors
}

func (m mqttChannel) Out() chan hub.Input {
	return m.out
}
