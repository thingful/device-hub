package pipe

import (
	"errors"

	"github.com/thingful/expando"
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

	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			errors <- err
		},
	})

	err := cli.Connect(options)

	if err != nil {
		return nil, err
	}

	return &mqttbroker{
		topic:  topic,
		errors: errors,
		client: cli,
	}, nil
}

type mqttbroker struct {
	topic  string
	errors chan error
	client *client.Client
}

func (m *mqttbroker) Channel() Channel {

	channel := make(chan expando.Input)

	// TODO : separate out subscriptions from channels
	err := m.client.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{

				TopicFilter: []byte(m.topic),
				QoS:         mqtt.QoS0,

				Handler: func(topicName, message []byte) {

					input := expando.Input{
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
		panic(err)
	}

	return mqttChannel{out: channel, errors: m.errors}
}

func (m *mqttbroker) Close() error {
	err := m.client.Disconnect()
	defer m.client.Terminate()
	return err
}

type mqttChannel struct {
	errors chan error
	out    chan expando.Input
}

// Errors returns a channel of errors
func (m mqttChannel) Errors() chan error {
	return m.errors
}

func (m mqttChannel) Out() chan expando.Input {
	return m.out
}
