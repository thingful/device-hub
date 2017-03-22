package pipe

import (
	"errors"
	"sync"

	hub "github.com/thingful/device-hub"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

const (
	TOPIC_NAME = "topic"
)

func FromMQTT(options *client.ConnectOptions, topic string) (*mqttbroker, error) {

	if topic == "" {
		return nil, errors.New("mqtt topic is empty string")
	}

	errors := make(chan error)
	return &mqttbroker{
		topic:   topic,
		errors:  errors,
		options: options,
	}, nil
}

type mqttbroker struct {
	topic           string
	errors          chan error
	client          *client.Client
	options         *client.ConnectOptions
	connection_lock sync.Mutex
}

func (m *mqttbroker) Channel() (Channel, error) {

	err := m.connect()

	if err != nil {
		return NoOpChannel{}, err
	}

	channel := make(chan hub.Input)

	err = m.client.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{

				TopicFilter: []byte(m.topic),
				QoS:         mqtt.QoS0,

				Handler: func(topicName, message []byte) {

					input := hub.Input{
						Payload: message,
						Metadata: map[string]interface{}{
							TOPIC_NAME: topicName,
						},
					}

					channel <- input
				},
			},
		},
	})
	if err != nil {
		return NoOpChannel{}, err
	}

	return mqttChannel{out: channel, errors: m.errors}, nil
}

func (m *mqttbroker) connect() error {

	m.connection_lock.Lock()
	defer m.connection_lock.Unlock()

	if m.client != nil {
		return nil
	}

	m.client = client.New(&client.Options{
		ErrorHandler: func(err error) {
			m.errors <- err
		},
	})

	err := m.client.Connect(m.options)

	if err != nil {
		m.client = nil
		return err
	}

	return nil
}

func (m *mqttbroker) Close() error {
	err := m.client.Disconnect()
	defer m.client.Terminate()
	return err
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
