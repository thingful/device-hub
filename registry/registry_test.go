// Copyright © 2017 thingful

package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"
)

type mockEndpoint struct {
	count int
}

func (m mockEndpoint) Write(v hub.Message) error {
	return nil
}

func TestBuildersAreCached(t *testing.T) {

	count := 0
	registry := New()

	registry.RegisterEndpoint("simple", func(config describe.Values) (hub.Endpoint, error) {

		count++
		return mockEndpoint{count: count}, nil

	}, describe.Parameters{
		describe.Parameter{},
	})

	one, err := registry.EndpointByName("foo", "simple", map[string]string{})
	assert.Nil(t, err)

	two, err := registry.EndpointByName("foo", "simple", map[string]string{})
	assert.Nil(t, err)

	assert.Equal(t, one, two)
	assert.Equal(t, one.(mockEndpoint).count, two.(mockEndpoint).count)

	three, err := registry.EndpointByName("bar", "simple", map[string]string{})

	assert.Nil(t, err)
	assert.NotEqual(t, one, three)
	assert.NotEqual(t, one.(mockEndpoint).count, three.(mockEndpoint).count)

}

func TestErrorThrownForIncorrectType(t *testing.T) {

	registry := New()

	registry.RegisterEndpoint("endpoint", func(config describe.Values) (hub.Endpoint, error) {

		return mockEndpoint{}, nil

	}, describe.Parameters{
		describe.Parameter{},
	})

	_, err := registry.ListenerByName("foo", "endpoint", map[string]string{})

	assert.NotNil(t, err)
}
