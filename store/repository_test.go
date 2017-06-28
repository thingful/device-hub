// Copyright © 2017 thingful

package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/device-hub/proto"

	testing_helper "github.com/thingful/device-hub/utils/testing"
)

func TestUIDIsSetOnEntity(t *testing.T) {
	t.Parallel()

	listener := &proto.Entity{Type: "listener"}
	endpoint := &proto.Entity{Type: "endpoint"}

	err := ensureEntityHasUID(listener)
	assert.Nil(t, err)

	err = ensureEntityHasUID(endpoint)
	assert.Nil(t, err)

	// uids shouldn't be equal
	assert.NotEqual(t, listener.Uid, endpoint.Uid)
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

func TestProfileUIDErrorsIfNoProfileName(t *testing.T) {
	t.Parallel()

	profile := &proto.Entity{
		Type:          "profile",
		Configuration: map[string]string{}}

	err := ensureEntityHasUID(profile)
	assert.NotNil(t, err)

}

func TestBoltDB_ListenersCreatedSearchedAndDeleted(t *testing.T) {

	t.Parallel()

	conn := testing_helper.MustDialBoltDB()
	defer conn.MustClose()

	store := NewStore(conn.DB)

	r := &mockregister{
		listenerRegistered: true,
	}

	repository := NewRepository(store, r)

	entity := proto.Entity{
		Type: "listener",
		Kind: "http",
		Uid:  "xxx",
	}

	uid, err := repository.UpdateOrCreateEntity(entity)
	assert.Nil(t, err)
	assert.Equal(t, entity.Uid, uid)

	listeners, err := repository.Search("l")

	assert.Nil(t, err)
	assert.Len(t, listeners, 1)

	l, err := repository.Listeners.One("xxx")

	assert.Nil(t, err)
	assert.Equal(t, entity.Uid, l.Uid)

	err = repository.Listeners.Delete("xxx")
	assert.Nil(t, err)

	listeners, err = repository.Search("l")

	assert.Nil(t, err)
	assert.Len(t, listeners, 0)

}

type mockregister struct {
	listenerRegistered bool
	endpointRegistered bool
}

func (r *mockregister) IsListenerRegistered(string) bool {
	return r.listenerRegistered
}

func (r *mockregister) IsEndpointRegistered(string) bool {
	return r.endpointRegistered
}
