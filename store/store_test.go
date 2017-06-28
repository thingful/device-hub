package store

import (
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/device-hub/proto"
	testing_helper "github.com/thingful/device-hub/utils/testing"
)

// storeTester is the signature for store tests
type storeTester func(t *testing.T, store Storer)

var (

	// AllTests contain the complete list of tests that must pass for each Storer
	AllTests = []storeTester{
		ListenersCreatedSearchedAndDeleted,
		EndpointsCreatedSearchedAndDeleted,
		InsertShouldReturnErrorIfSameItemAddedTwice,
	}
)

// TestStores will test the complete set of store tests against each implementation
func TestStores(t *testing.T) {

	t.Parallel()

	// test bolt db

	for _, test := range AllTests {

		t.Logf("boltdb : %s", runtime.FuncForPC(reflect.ValueOf(test).Pointer()).Name())

		conn := testing_helper.MustDialBoltDB()
		defer conn.MustClose()

		s := NewBoltDBStore(conn.DB)

		test(t, s)

	}

	// test fileStore
	for _, test := range AllTests {

		t.Logf("fileStore : %s", runtime.FuncForPC(reflect.ValueOf(test).Pointer()).Name())

		// create a temp folder for the file store
		f, err := ioutil.TempDir("", "bolt-test-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f)

		s := NewFileStore(f)
		test(t, s)

	}
}

func InsertShouldReturnErrorIfSameItemAddedTwice(t *testing.T, store Storer) {

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

	_, err = repository.UpdateOrCreateEntity(entity)
	assert.NotNil(t, err)
	assert.Equal(t, ErrItemAlreadyExists, err)

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
