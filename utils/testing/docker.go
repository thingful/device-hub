// Copyright © 2017 thingful
package testing

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	mqtt_helper "github.com/thingful/device-hub/utils/mqtt"

	// TODO : move import to upstream project
	"github.com/mdevilliers/go-compose/compose"
)

var (
	ciEnvVar = "CI"
)

// IsRunningInCI tests for a well know environmental variable and returns true if found, false if absent
func IsRunningInCI() bool {
	return os.Getenv(ciEnvVar) != ""
}

type testingEnvironment struct {
	MQTTClient mqtt.Client
	compose    *compose.Compose
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

	mqttAddress := fmt.Sprintf("tcp://%s:%d", "0.0.0.0", 1883)
	client := mqtt_helper.DefaultMQTTClient(mqttAddress, "device-hub")

	if IsRunningInCI() {

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return nil, token.Error()
		}

	} else {
		composeYML := `version: '2'
services:
  mqtt:
    image: erlio/docker-vernemq:0.15.3
    ports:
      - 1883
    environment:
      - DOCKER_VERNEMQ_ALLOW_ANONYMOUS=on`

		c = compose.MustStartParallel(composeYML, false)

		mqttAddress = fmt.Sprintf("tcp://%s:%d", compose.MustInferDockerHost(), c.Containers["mqtt"].MustGetFirstPublicPort(1883, "tcp"))
		client = mqtt_helper.DefaultMQTTClient(mqttAddress, "device-hub")

		compose.MustConnectWithDefaults(func() error {
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				return token.Error()
			}

			return nil
		})
	}

	return &testingEnvironment{
		compose:    c,
		MQTTClient: client,
	}, nil

}

func (t *testingEnvironment) Down() {
	t.MQTTClient.Disconnect(1)
	t.compose.Kill()
}
