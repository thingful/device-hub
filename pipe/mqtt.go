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

func NewMQTTChannel(client mqtt.Client, topic string) (Channel, error) {

	if topic == "" {
		return nil, errors.New("mqtt topic is empty string")
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

	if token := client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		return NoOpChannel{}, token.Error()
	}

	return defaultChannel{out: out, errors: errors}, nil
}
