// Copyright © 2017 thingful

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	ok, _ := isEmpty(nil)

	assert.False(t, ok)

	ok, _ = isEmpty(&Configuration{})
	assert.False(t, ok)

	// test no pipes
	ok, _ = isEmpty(&Configuration{
		Listeners: endpoints{
			Endpoint{},
		},
		Endpoints: endpoints{
			Endpoint{},
		},
	})

	assert.False(t, ok)

	// test no listeners
	ok, _ = isEmpty(&Configuration{
		Pipes: pipes{
			pipe{},
		},
		Endpoints: endpoints{
			Endpoint{},
		},
	})

	assert.False(t, ok)

	// test no endpoints
	ok, _ = isEmpty(&Configuration{
		Pipes: pipes{
			pipe{},
		},
		Listeners: endpoints{
			Endpoint{},
		},
	})

	assert.False(t, ok)

	// valid
	ok, _ = isEmpty(&Configuration{
		Pipes: pipes{
			pipe{},
		},
		Listeners: endpoints{
			Endpoint{},
		},
		Endpoints: endpoints{
			Endpoint{},
		},
	})

	assert.True(t, ok)

}

func TestUniqueUIDs(t *testing.T) {

	t.Parallel()

	config := &Configuration{
		Listeners: endpoints{
			Endpoint{UID: UID("foo")},
		},
		Endpoints: endpoints{
			Endpoint{UID: UID("foo")},
		},
	}

	ok, _ := uniqueUIDs(config)
	assert.False(t, ok)

	config2 := &Configuration{
		Listeners: endpoints{
			Endpoint{UID: UID("foo")},
		},
		Endpoints: endpoints{
			Endpoint{UID: UID("bar")},
		},
	}
	ok, _ = uniqueUIDs(config2)
	assert.True(t, ok)

}

func TestNoEmptyUIDs(t *testing.T) {

	t.Parallel()

	config := &Configuration{
		Listeners: endpoints{
			Endpoint{UID: UID("")},
		},
		Endpoints: endpoints{
			Endpoint{UID: UID("foo")},
		},
	}

	ok, _ := noEmptyUIDs(config)
	assert.False(t, ok)

	config2 := &Configuration{
		Listeners: endpoints{
			Endpoint{UID: UID("foo")},
		},
		Endpoints: endpoints{
			Endpoint{UID: UID("bar")},
		},
	}
	ok, _ = noEmptyUIDs(config2)
	assert.True(t, ok)

}

func TestUniquePipeURIs(t *testing.T) {

	t.Parallel()

	config := &Configuration{
		Pipes: pipes{
			pipe{Uri: "foo"},
			pipe{Uri: "foo"},
		},
	}

	ok, _ := uniquePipeURIs(config)
	assert.False(t, ok)

	config2 := &Configuration{
		Pipes: pipes{
			pipe{Uri: "foo"},
			pipe{Uri: "bar"},
		},
	}

	ok, _ = uniquePipeURIs(config2)
	assert.True(t, ok)

}

func TestValidatePipes_WithValidConfig(t *testing.T) {

	t.Parallel()

	config := &Configuration{
		Listeners: endpoints{
			Endpoint{UID: UID("foo")},
		},
		Endpoints: endpoints{
			Endpoint{UID: UID("bar")},
		},
		Profiles: profiles{
			Profile{UID: UID("cat")},
		},
		Pipes: pipes{
			pipe{
				Uri:       "/one",
				Profile:   UID("cat"),
				Listener:  UID("foo"),
				Endpoints: []UID{UID("bar")},
			},
		},
	}

	ok, _ := validatePipes(config)
	assert.True(t, ok)

}

func TestValidatePipes_WithMissingListener(t *testing.T) {

	t.Parallel()

	config := &Configuration{
		Listeners: endpoints{},
		Endpoints: endpoints{
			Endpoint{UID: UID("bar")},
		},
		Profiles: profiles{
			Profile{UID: UID("cat")},
		},
		Pipes: pipes{
			pipe{
				Uri:       "/one",
				Profile:   UID("cat"),
				Listener:  UID("foo"),
				Endpoints: []UID{UID("bar")},
			},
		},
	}

	ok, _ := validatePipes(config)
	assert.False(t, ok)

}

func TestValidatePipes_WithMissingProfile(t *testing.T) {

	t.Parallel()

	config := &Configuration{
		Listeners: endpoints{
			Endpoint{UID: UID("foo")},
		},
		Endpoints: endpoints{
			Endpoint{UID: UID("bar")},
		},
		Profiles: profiles{},
		Pipes: pipes{
			pipe{
				Uri:       "/one",
				Profile:   UID("cat"),
				Listener:  UID("foo"),
				Endpoints: []UID{UID("bar")},
			},
		},
	}

	ok, _ := validatePipes(config)
	assert.False(t, ok)

}

func TestValidatePipes_WithMissingEndpoint(t *testing.T) {

	t.Parallel()

	config := &Configuration{
		Listeners: endpoints{
			Endpoint{UID: UID("foo")},
		},
		Endpoints: endpoints{},
		Profiles: profiles{
			Profile{UID: UID("cat")},
		},
		Pipes: pipes{
			pipe{
				Uri:       "/one",
				Profile:   UID("cat"),
				Listener:  UID("foo"),
				Endpoints: []UID{UID("bar")},
			},
		},
	}

	ok, _ := validatePipes(config)
	assert.False(t, ok)

}
