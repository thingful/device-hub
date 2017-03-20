package pipe

import (
	"context"
	"fmt"

	"github.com/thingful/expando"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func FromMQTT(ctx context.Context) *mqttbroker {

	errors := make(chan error)

	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			errors <- err
		},
	})

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "0.0.0.0:1883",
		ClientID: []byte("expando-client"),
	})

	if err != nil {
		panic(err)
	}

	return &mqttbroker{
		ctx:    ctx,
		errors: errors,
		client: cli,
	}
}

type mqttbroker struct {
	ctx    context.Context
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
					fmt.Println(string(topicName), string(message))
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
