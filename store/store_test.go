// Copyright © 2017 thingful

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
		ProfilesCreatedSearchedAndDeleted,
		PipesCreatedSearchedAndDeleted,
		InsertShouldReturnErrorIfSameItemAddedTwice,
		SearchAllShouldReturnAllEntities,
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

func SearchAllShouldReturnAllEntities(t *testing.T, store Storer) {
	r := &mockregister{
		endpointRegistered: true,
		listenerRegistered: true,
	}

	repository := NewRepository(store, r)

	one := proto.Entity{
		Type: "endpoint",
		Kind: "http",
		Uid:  "one",
	}

	two := proto.Entity{
		Type: "listener",
		Kind: "http",
		Uid:  "two",
	}

	three := proto.Entity{
		Type: "profile",
		Kind: "script",
		Configuration: map[string]string{
			"profile-name": "foo/bar",
		},
	}

	uid, err := repository.Insert(one)
	assert.Nil(t, err)
	assert.Equal(t, one.Uid, uid)

	uid, err = repository.Insert(two)
	assert.Nil(t, err)
	assert.Equal(t, two.Uid, uid)

	uid, err = repository.Insert(three)
	assert.Nil(t, err)
	assert.Equal(t, "foo/bar", uid)

	all, err := repository.Search("e,l,p")

	assert.Nil(t, err)
	assert.Len(t, all, 3)

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

	uid, err := repository.Insert(entity)
	assert.Nil(t, err)
	assert.Equal(t, entity.Uid, uid)

	_, err = repository.Insert(entity)
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

	uid, err := repository.Insert(entity)
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

	uid, err := repository.Insert(entity)
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

func ProfilesCreatedSearchedAndDeleted(t *testing.T, store Storer) {

	r := &mockregister{}

	repository := NewRepository(store, r)

	entity := proto.Entity{
		Type: "profile",
		Kind: "script",
		Configuration: map[string]string{
			"profile-name": "foo/bar",
		},
	}

	uid, err := repository.Insert(entity)

	assert.Nil(t, err)
	assert.Equal(t, "foo/bar", uid)

	profiles, err := repository.Search("p")

	assert.Nil(t, err)
	assert.Len(t, profiles, 1)

	p, err := repository.Profiles.One("foo/bar")

	assert.Nil(t, err)
	assert.Equal(t, "foo/bar", p.Uid)

	err = repository.Profiles.Delete("foo/bar")
	assert.Nil(t, err)

	profiles, err = repository.Search("p")
	assert.Nil(t, err)
	assert.Len(t, profiles, 0)
}

func PipesCreatedSearchedAndDeleted(t *testing.T, store Storer) {

	r := &mockregister{}

	repository := NewRepository(store, r)

	pipe := Pipe{
		Profile:   Profile{},
		Listener:  &proto.Entity{},
		Endpoints: []*proto.Entity{},
		Uri:       "/foo",
		Tags:      map[string]string{},
	}

	err := repository.Pipes.Insert(pipe)

	assert.Nil(t, err)

	pipes, err := repository.Pipes.List()

	assert.Nil(t, err)
	assert.Len(t, pipes, 1)

	err = repository.Pipes.Delete("/foo")

	assert.Nil(t, err)

	pipes, err = repository.Pipes.List()

	assert.Nil(t, err)
	assert.Len(t, pipes, 0)

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
