// Copyright © 2017 thingful
package describe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValues_AllValuesRequired(t *testing.T) {

	t.Parallel()

	config := map[string]string{
		"an-int":   "1",
		"an-int64": "2",
		"a-string": "hello",
		"a-bool":   "true",
		"a-url":    "https://abc.com",
	}

	params := Parameters{
		Parameter{Name: "an-int", Type: Int32, Required: true},
		Parameter{Name: "an-int64", Type: Int64, Required: true},
		Parameter{Name: "a-string", Type: String, Required: true},
		Parameter{Name: "a-bool", Type: Bool, Required: true},
		Parameter{Name: "a-url", Type: Url, Required: true},
	}

	values, err := NewValues(config, params)

	assert.Nil(t, err)

	i, ifound := values.Int32("an-int")

	assert.True(t, ifound)
	assert.Equal(t, int32(1), i)

	i64, i64found := values.Int64("an-int64")

	assert.True(t, i64found)
	assert.Equal(t, int64(2), i64)

	s, sfound := values.String("a-string")
	assert.True(t, sfound)
	assert.Equal(t, "hello", s)

	b, bfound := values.Bool("a-bool")
	assert.True(t, bfound)
	assert.True(t, b)

}

func TestNewValues_NoValuesRequired(t *testing.T) {

	t.Parallel()

	config := map[string]string{
		"an-int":   "1",
		"an-int64": "2",
		"a-string": "hello",
		"a-bool":   "true",
		"a-url":    "https://abc.com",
	}

	params := Parameters{
		Parameter{Name: "an-int", Type: Int32, Required: false},
		Parameter{Name: "an-int64", Type: Int64, Required: true},
		Parameter{Name: "a-string", Type: String, Required: false},
		Parameter{Name: "a-bool", Type: Bool, Required: false},
		Parameter{Name: "a-url", Type: Url, Required: false},
	}

	values, err := NewValues(config, params)

	assert.Nil(t, err)

	i, ifound := values.Int32("an-int")

	assert.True(t, ifound)
	assert.Equal(t, int32(1), i)

	i64, i64found := values.Int64("an-int64")

	assert.True(t, i64found)
	assert.Equal(t, int64(2), i64)

	s, sfound := values.String("a-string")
	assert.True(t, sfound)
	assert.Equal(t, "hello", s)

	b, bfound := values.Bool("a-bool")
	assert.True(t, bfound)
	assert.True(t, b)

}

func TestNewValues_InvalidURL(t *testing.T) {

	t.Parallel()

	config := map[string]string{
		"a-url": "not a url",
	}

	params := Parameters{
		Parameter{Name: "a-url", Type: Url, Required: true},
	}

	_, err := NewValues(config, params)

	assert.NotNil(t, err)

}
