// Copyright © 2017 thingful

package store

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strings"

	hashids "github.com/speps/go-hashids"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/proto"
)

type bucket struct {
	name  []byte
	store *Store
}

// Repository facilitates some higher level store interactions
type Repository struct {
	Listeners entityBucket
	Endpoints entityBucket
	Profiles  entityBucket
	Pipes     pipeBucket
	store     *Store
}

func NewRepository(store *Store) *Repository {
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
		store: store,
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

	var bucket bucket

	switch strings.ToLower(item.Type) {
	case "listener":
		bucket = e.Listeners.bucket

		exists := hub.IsListenerRegistered(item.Kind)

		if !exists {
			return "", fmt.Errorf("kind : %s not registered", item.Kind)
		}

	case "endpoint":
		bucket = e.Endpoints.bucket

		exists := hub.IsEndpointRegistered(item.Kind)
		if !exists {
			return "", fmt.Errorf("kind : %s not registered", item.Kind)
		}

	case "profile":
		bucket = e.Profiles.bucket

	default:
		return "", fmt.Errorf("type : %s not registered", item.Type)
	}

	err := ensureEntityHasUID(&item)

	if err != nil {
		return "", err
	}

	err = e.store.Insert(bucket, []byte(item.Uid), item)

	if err != nil {
		return "", err
	}

	return item.Uid, nil
}

func (e *Repository) Delete(typez, uid string) error {

	switch strings.ToLower(typez) {
	case "listener":
		return e.Listeners.Delete(uid)
	case "endpoint":
		return e.Endpoints.Delete(uid)
	case "profile":
		return e.Profiles.Delete(uid)
	case "pipes":
		return e.Pipes.Delete(uid)

	default:
		return fmt.Errorf("type : %s not found", typez)
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
