package pipe

import (
	"github.com/thingful/expando"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func FromMQTT(options *client.ConnectOptions) (*mqttbroker, error) {

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
		errors: errors,
		client: cli,
	}, nil
}

type mqttbroker struct {
	errors chan error
	client *client.Client
}

func (m *mqttbroker) Channel() Channel {

	channel := make(chan expando.Input)

	err := m.client.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("#"),
				QoS:         mqtt.QoS0,

				Handler: func(topicName, message []byte) {

					// TODO : add topic name to metadata
					channel <- expando.Input{Payload: message}
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
