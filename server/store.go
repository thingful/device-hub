// Copyright © 2017 thingful

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
)

type store struct {
	db *bolt.DB
}

type bucket []byte

var (
	ErrSlicePtrNeeded = errors.New("slice ptr needed")
	ErrNotFound       = errors.New("not found")

	endpointsBucket = bucket([]byte("endpoints"))
	listenersBucket = bucket([]byte("listeners"))
	profilesBucket  = bucket([]byte("profiles"))
)

func NewStore(db *bolt.DB) (*store, error) {

	err := db.Update(func(tx *bolt.Tx) error {

		mustCreateBucket(tx, endpointsBucket)
		mustCreateBucket(tx, listenersBucket)
		mustCreateBucket(tx, profilesBucket)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &store{
		db: db,
	}, nil
}

func mustCreateBucket(tx *bolt.Tx, bucket bucket) {
	_, err := tx.CreateBucketIfNotExists(bucket)
	if err != nil {
		panic(fmt.Sprintf("create bucket failed : %s", err))
	}
}

func (s *store) Insert(bucket bucket, uid []byte, data interface{}) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)

		existing := b.Get(uid)

		if len(existing) > 0 {
			return errors.New("item already exists")
		}

		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return b.Put([]byte(uid), buf)

	})
	return err
}

func (s *store) Update(bucket bucket, uid string, data interface{}) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)

		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return b.Put([]byte(uid), buf)

	})
	return err
}

func (s *store) Delete(bucket bucket, uid string) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)

		return b.Delete([]byte(uid))
	})

	return err
}

func (s *store) One(b bucket, uid string, out interface{}) error {

	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(b)

		bytes := b.Get([]byte(uid))

		if len(bytes) > 0 {

			err := json.Unmarshal(bytes, &out)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (s *store) List(bucket bucket, to interface{}) error {

	ref := reflect.ValueOf(to)

	if ref.Kind() != reflect.Ptr || reflect.Indirect(ref).Kind() != reflect.Slice {
		return ErrSlicePtrNeeded
	}

	list := map[string][]byte{}
	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)
		b.ForEach(func(k, v []byte) error {
			list[string(k)] = v
			return nil
		})

		return nil
	})

	if err != nil {
		return err
	}

	results := reflect.MakeSlice(reflect.Indirect(ref).Type(), len(list), len(list))
	i := 0
	for k, _ := range list {
		raw := list[k]
		if raw == nil {
			return ErrNotFound
		}

		err = json.Unmarshal(raw, results.Index(i).Addr().Interface())
		if err != nil {
			return err
		}
		i++
	}

	reflect.Indirect(ref).Set(results)

	return nil
}
