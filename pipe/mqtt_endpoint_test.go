// +build integration

package pipe

import (
	"fmt"
	"testing"

	// TODO : move import to upstream project
	"github.com/mdevilliers/go-compose/compose"
	"github.com/stretchr/testify/assert"
)

func TestMQTT_MultipleEndpoints(t *testing.T) {

	t.Parallel()

	composeYML := `version: '2'
services:
  mqtt:
    image: erlio/docker-vernemq:0.15.3
    ports:
      - 1883
    environment:
      - DOCKER_VERNEMQ_ALLOW_ANONYMOUS=on`

	c := compose.MustStartParallel(composeYML, false)
	defer c.Kill()

	mqttAddress := fmt.Sprintf("tcp://%s:%d", compose.MustInferDockerHost(), c.Containers["mqtt"].MustGetFirstPublicPort(1883, "tcp"))

	options := DefaultMQTTOptions(mqttAddress, "device-hub")
	client := DefaultClient(options)

	compose.MustConnectWithDefaults(func() error {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}

		return nil
	})

	defer client.Disconnect(1)

	channel1, err := NewMQTTChannel(client, "/a")
	assert.Nil(t, err)

	channel2, err := NewMQTTChannel(client, "/b")
	assert.Nil(t, err)

	client.Publish("/a", 0, false, "hello")
	client.Publish("/b", 0, false, "hello")

	message := <-channel1.Out()
	assert.Equal(t, message.Payload, []byte("hello"))

	message2 := <-channel2.Out()
	assert.Equal(t, message2.Payload, []byte("hello"))

}
