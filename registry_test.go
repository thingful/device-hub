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

func TestEndpointsAreCached(t *testing.T) {

	count := 0

	RegisterEndpoint("simple", func(config utils.TypedMap) (Endpoint, error) {

		count++
		return mockEndpoint{count: count}, nil

	})

	one, err := EndpointByName("foo", "simple", map[string]interface{}{})

	assert.Nil(t, err)

	two, err := EndpointByName("foo", "simple", map[string]interface{}{})

	assert.Nil(t, err)
	assert.Equal(t, one, two)

	three, err := EndpointByName("bar", "simple", map[string]interface{}{})

	assert.Nil(t, err)
	assert.NotEqual(t, one, three)

}
