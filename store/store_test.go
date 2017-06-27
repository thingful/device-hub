package store

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/device-hub/proto"
	testing_helper "github.com/thingful/device-hub/utils/testing"
)

// TestStores will test the complete set of store tests against each implementation
func TestStores(t *testing.T) {

	t.Parallel()

	// test bolt db first
	conn := testing_helper.MustDialBoltDB()
	defer conn.MustClose()

	boltDB := NewBoltDBStore(conn.DB)

	ListenersCreatedSearchedAndDeleted(t, boltDB)
	EndpointsCreatedSearchedAndDeleted(t, boltDB)

	// crate a folder for the file store
	f, err := ioutil.TempDir("", "bolt-test-")
	if err != nil {
		t.Fatal(err)
	}

	fileStore := NewFileStore(f)

	ListenersCreatedSearchedAndDeleted(t, fileStore)
	EndpointsCreatedSearchedAndDeleted(t, boltDB)

}

func EndpointsCreatedSearchedAndDeleted(t *testing.T, store Storer) {

	r := &mockregister{
		endpointRegistered: true,
	}

	repository := NewRepository(store, r)

	entity := proto.Entity{
		Type: "endpoint",
		Kind: "http",
		Uid:  "xxx",
	}

	uid, err := repository.UpdateOrCreateEntity(entity)
	assert.Nil(t, err)
	assert.Equal(t, entity.Uid, uid)

	endpoints, err := repository.Search("e")

	assert.Nil(t, err)
	assert.Len(t, endpoints, 1)

	e, err := repository.Endpoints.One("xxx")

	assert.Nil(t, err)
	assert.Equal(t, entity.Uid, e.Uid)

	err = repository.Endpoints.Delete("xxx")
	assert.Nil(t, err)

	endpoints, err = repository.Search("e")

	assert.Nil(t, err)
	assert.Len(t, endpoints, 0)

}

func ListenersCreatedSearchedAndDeleted(t *testing.T, store Storer) {

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
