package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/device-hub/proto"
)

func TestUIDIsSetOnEntity(t *testing.T) {
	t.Parallel()

	listener := &proto.Entity{Type: "listener"}
	endpoint := &proto.Entity{Type: "endpoint"}
	profile := &proto.Entity{Type: "profile"}

	err := ensureEntityHasUID(listener)
	assert.Nil(t, err)

	err = ensureEntityHasUID(endpoint)
	assert.Nil(t, err)

	err = ensureEntityHasUID(profile)
	assert.Nil(t, err)

	// uids shouldn't be equal
	assert.NotEqual(t, listener.Uid, endpoint.Uid)
	assert.NotEqual(t, endpoint.Uid, profile.Uid)
	assert.NotEqual(t, profile.Uid, listener.Uid)
}

func TestUIDIsNotSetOnEntity(t *testing.T) {
	t.Parallel()

	listener := &proto.Entity{Type: "listener", Uid: "foo"}
	endpoint := &proto.Entity{Type: "endpoint", Uid: "bar"}
	profile := &proto.Entity{Type: "profile", Uid: "foobar"}

	err := ensureEntityHasUID(listener)
	assert.Nil(t, err)

	err = ensureEntityHasUID(endpoint)
	assert.Nil(t, err)

	err = ensureEntityHasUID(profile)
	assert.Nil(t, err)

	assert.Equal(t, "foo", listener.Uid)
	assert.Equal(t, "bar", endpoint.Uid)
	assert.Equal(t, "foobar", profile.Uid)
}

func TestProfileUIDDefaultsToProfileName(t *testing.T) {
	t.Parallel()

	profile := &proto.Entity{
		Type: "profile",
		Configuration: map[string]string{
			"profile-name": "foobar",
		}}

	err := ensureEntityHasUID(profile)
	assert.Nil(t, err)

	assert.Equal(t, "foobar", profile.Uid)
}
