// Copyright © 2017 thingful

package hub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/device-hub/utils"
)

type mockEndpoint struct {
	count int
}

func (m mockEndpoint) Write(v Message) error {
	return nil
}

func TestBuildersAreCached(t *testing.T) {

	count := 0

	RegisterEndpoint("simple", func(config utils.TypedMap) (Endpoint, error) {

		count++
		return mockEndpoint{count: count}, nil

	})

	one, err := EndpointByName("foo", "simple", utils.TypedMap{})
	assert.Nil(t, err)

	two, err := EndpointByName("foo", "simple", utils.TypedMap{})
	assert.Nil(t, err)

	assert.Equal(t, one, two)
	assert.Equal(t, one.(mockEndpoint).count, two.(mockEndpoint).count)

	three, err := EndpointByName("bar", "simple", utils.TypedMap{})

	assert.Nil(t, err)
	assert.NotEqual(t, one, three)
	assert.NotEqual(t, one.(mockEndpoint).count, three.(mockEndpoint).count)

}

func TestErrorThrownForIncorrectType(t *testing.T) {

	RegisterEndpoint("endpoint", func(config utils.TypedMap) (Endpoint, error) {

		return mockEndpoint{}, nil

	})

	_, err := ListenerByName("foo", "endpoint", utils.TypedMap{})

	assert.NotNil(t, err)
}
