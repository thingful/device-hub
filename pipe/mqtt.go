package pipe

import (
	"errors"
	"strings"
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
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)

	return opts
}

var (
	client      mqtt.Client
	client_lock sync.Mutex
)

func DefaultMQTTClient(options *mqtt.ClientOptions) mqtt.Client {

	client_lock.Lock()
	defer client_lock.Unlock()

	if client != nil {
		return client
	}

	client = mqtt.NewClient(options)

	return client

}

func NewMQTTListener(config map[string]interface{}, client mqtt.Client) (*mqttlistener, error) {

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

func (m *mqttlistener) NewChannel(topic string) (Channel, error) {

	if topic == "" {
		return nil, errors.New("mqtt topic is empty string")
	}

	if strings.Contains(topic, "#") {
		return nil, errors.New("mqtt wildcard (#) is not allowed")
	}

	errors := make(chan error)
	out := make(chan hub.Input)

	handler := func(client mqtt.Client, msg mqtt.Message) {
		input := hub.Input{
			Payload: msg.Payload(),
			Metadata: map[string]interface{}{
				TOPIC_NAME: msg.Topic(),
			},
		}

		out <- input

	}

	if token := m.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		return NoOpChannel{}, token.Error()
	}

	return defaultChannel{out: out, errors: errors}, nil
}

func (m *mqttlistener) Close() error {
	// TODO : set a sensible timeout
	m.client.Disconnect(1000)
	return nil
}
