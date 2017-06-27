// Copyright © 2017 thingful

package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"strings"

	hashids "github.com/speps/go-hashids"
	"github.com/thingful/device-hub/proto"
)

type bucket struct {
	name  []byte
	store Storer
}

// Repository facilitates some higher level store interactions
type Repository struct {
	Listeners entityBucket
	Endpoints entityBucket
	Profiles  entityBucket
	Pipes     pipeBucket
	store     Storer
	register  register
}

// Storer is the interface a low level storage mechanism needs to implement
type Storer interface {
	// MustCreateBuckets will ensure the underlying storage exists for the entities
	MustCreateBuckets(buckets []bucket)

	// Insert will insert an entity or error
	Insert(bucket bucket, uid []byte, data interface{}) error

	// Delete will remove an entity or error
	Delete(bucket bucket, uid []byte) error

	// One will return an entity of error
	One(bucket bucket, uid []byte, out interface{}) error

	// List will return an array of entities or error
	List(bucket bucket, to interface{}) error
}

var (

	// ErrSlicePtrNeeded is returned when deserialisation requires an array of pointers
	ErrSlicePtrNeeded = errors.New("slice ptr needed")

	// ErrNotFound is returned if the item searched for doesn not exist
	ErrNotFound = errors.New("not found")

	// ErrItemAlreadyExists is returned if inserting an item that already exists
	ErrItemAlreadyExists = errors.New("item already exists")
)

type register interface {
	IsEndpointRegistered(string) bool
	IsListenerRegistered(string) bool
}

func NewRepository(store Storer, register register) *Repository {
	r := &Repository{
		Listeners: entityBucket{
			bucket: bucket{name: []byte("listeners"),
				store: store,
			}},
		Endpoints: entityBucket{
			bucket: bucket{name: []byte("endpoints"),
				store: store,
			}},
		Profiles: entityBucket{
			bucket{name: []byte("profiles"),
				store: store,
			}},
		Pipes: pipeBucket{
			bucket{name: []byte("pipes"),
				store: store,
			}},
		store:    store,
		register: register,
	}

	store.MustCreateBuckets([]bucket{
		r.Listeners.bucket,
		r.Endpoints.bucket,
		r.Profiles.bucket,
		r.Pipes.bucket,
	})

	return r
}

func (e *Repository) UpdateOrCreateEntity(item proto.Entity) (string, error) {

	var b bucket

	switch strings.ToLower(item.Type) {
	case "listener":
		b = e.Listeners.bucket

		exists := e.register.IsListenerRegistered(item.Kind)

		if !exists {
			return "", fmt.Errorf("kind : %s not registered", item.Kind)
		}

	case "endpoint":
		b = e.Endpoints.bucket

		exists := e.register.IsEndpointRegistered(item.Kind)
		if !exists {
			return "", fmt.Errorf("kind : %s not registered", item.Kind)
		}

	case "profile":
		b = e.Profiles.bucket

	default:
		return "", fmt.Errorf("type : %s not registered", item.Type)
	}

	err := ensureEntityHasUID(&item)

	if err != nil {
		return "", err
	}

	err = e.store.Insert(b, []byte(item.Uid), item)

	if err != nil {
		return "", err
	}

	return item.Uid, nil
}

func (e *Repository) Delete(entity proto.Entity) error {

	err := ensureEntityHasUID(&entity)

	if err != nil {
		return err
	}

	switch strings.ToLower(entity.Type) {
	case "listener":
		return e.Listeners.Delete(entity.Uid)
	case "endpoint":
		return e.Endpoints.Delete(entity.Uid)
	case "profile":
		return e.Profiles.Delete(entity.Uid)
	case "pipes":
		return e.Pipes.Delete(entity.Uid)

	default:
		return fmt.Errorf("type : %s not found", entity.Type)
	}
}

func (e *Repository) Search(filter string) ([]*proto.Entity, error) {

	keys := strings.Split(filter, ",")

	all := []*proto.Entity{}

	for _, key := range keys {

		switch strings.ToLower(key) {
		case "listener", "l":

			l := []*proto.Entity{}

			err := e.store.List(e.Listeners.bucket, &l)

			if err != nil {
				return nil, err
			}
			all = append(all, l...)

		case "endpoint", "e":

			en := []*proto.Entity{}

			err := e.store.List(e.Endpoints.bucket, &en)

			if err != nil {
				return nil, err
			}
			all = append(all, en...)
		case "profile", "p":

			p := []*proto.Entity{}

			err := e.store.List(e.Profiles.bucket, &p)

			if err != nil {
				return nil, err
			}
			all = append(all, p...)

		default:
			return nil, fmt.Errorf("filter of type : %s not registered", key)
		}
	}

	return all, nil
}

func ensureEntityHasUID(entity *proto.Entity) error {

	if entity.Uid != "" {
		return nil
	}

	switch strings.ToLower(entity.Type) {

	// profile entities should default to the profile-name
	case "profile":

		if entity.Configuration["profile-name"] != "" {
			// TODO : consider adding version to the profile-name?
			// Would be useful for having multiple profiles running
			// at the same time.
			entity.Uid = entity.Configuration["profile-name"]
			return nil
		}

		return errors.New("profile uid cannot be created - no 'profile-name'")

	default:
		hash, err := hash(entity)

		if err != nil {
			return err
		}

		entity.Uid = string(hash)
		return nil

	}
	return nil
}

func hash(data interface{}) ([]byte, error) {

	j, err := json.Marshal(data)

	if err != nil {
		return []byte{}, err
	}

	checksum := crc32.ChecksumIEEE(j)
	h, err := hashids.New()

	if err != nil {
		return []byte{}, err
	}

	uid, err := h.Encode([]int{int(checksum)})

	if err != nil {
		return []byte{}, err
	}
	return []byte(uid), nil
}
