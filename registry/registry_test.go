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

//"github.com/thingful/device-hub/listener"

/*
func TestConfigurationFile(t *testing.T) {

	//for each file in ./test-configurations/
	folder := "./test-configurations/"
	listing, err := ioutil.ReadDir(folder)
	assert.Nil(t, err)

	for _, fi := range listing {

		folderPath := path.Join(folder, fi.Name())

		dm := iocodec.DefaultDecoders["yaml"]

		f, err := os.Open(folderPath)

		assert.Nil(t, err)

		in := dm.NewDecoder(f)
		entity := proto.Entity{}

		err = in.Decode(&entity)

		assert.Nil(t, err)
		assert.NotEmpty(t, entity.Kind)
		assert.NotEmpty(t, entity.Type)

		register := New()
		//listener.Register(register)
	}

	//load file

	//parse and describe

}*/
