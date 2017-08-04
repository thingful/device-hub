// Copyright © 2017 thingful
package testing

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	// TODO : move import to upstream project
	"github.com/mdevilliers/go-compose/compose"
)

type testingEnvironment struct {
	MQTTClientOptions *mqtt.ClientOptions
	MQTTClient        mqtt.Client
	compose           *compose.Compose
}

func MustUp() *testingEnvironment {
	t, err := Up()

	if err != nil {
		panic(err)
	}
	return t
}

func Up() (*testingEnvironment, error) {

	var c *compose.Compose

	composeYML := `version: '2'
services:
  mqtt:
    image: erlio/docker-vernemq:0.15.3
    ports:
      - 1883
    environment:
      - DOCKER_VERNEMQ_ALLOW_ANONYMOUS=on`

	c = compose.MustStartParallel(composeYML, false)

	brokerAddress := fmt.Sprintf("tcp://%s:%d", compose.MustInferDockerHost(), c.Containers["mqtt"].MustGetFirstPublicPort(1883, "tcp"))

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerAddress)
	opts.SetClientID("device-hub-test-client")

	client := mqtt.NewClient(opts)

	compose.MustConnectWithDefaults(func() error {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}

		return nil
	})

	return &testingEnvironment{
		compose:           c,
		MQTTClientOptions: mqtt.NewClientOptions().AddBroker(brokerAddress),
		MQTTClient:        client,
	}, nil

}

func (t *testingEnvironment) Down() {
	t.MQTTClient.Disconnect(1)
	t.compose.Kill()
}
