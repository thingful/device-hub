// Copyright © 2017 thingful

package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
)

type Store struct {
	db *bolt.DB
}

type Bucket []byte

var (
	ErrSlicePtrNeeded = errors.New("slice ptr needed")
	ErrNotFound       = errors.New("not found")

	Endpoints = Bucket([]byte("endpoints"))
	Listeners = Bucket([]byte("listeners"))
	Profiles  = Bucket([]byte("profiles"))
	Pipes     = Bucket([]byte("pipes"))
)

func NewStore(db *bolt.DB) (*Store, error) {

	err := db.Update(func(tx *bolt.Tx) error {

		mustCreateBucket(tx, Endpoints)
		mustCreateBucket(tx, Listeners)
		mustCreateBucket(tx, Profiles)
		mustCreateBucket(tx, Pipes)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func mustCreateBucket(tx *bolt.Tx, bucket Bucket) {
	_, err := tx.CreateBucketIfNotExists(bucket)
	if err != nil {
		panic(fmt.Sprintf("create bucket failed : %s", err))
	}
}

func (s *Store) Insert(bucket Bucket, uid []byte, data interface{}) error {

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

func (s *Store) Update(bucket Bucket, uid []byte, data interface{}) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)

		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return b.Put(uid, buf)

	})
	return err
}

func (s *Store) Delete(bucket Bucket, uid []byte) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)

		return b.Delete(uid)
	})

	return err
}

func (s *Store) One(bucket Bucket, uid []byte, out interface{}) error {

	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)

		bytes := b.Get(uid)

		if len(bytes) > 0 {

			err := json.Unmarshal(bytes, &out)

			if err != nil {
				return err
			}

			return nil
		}

		return ErrNotFound
	})

	return err
}

func (s *Store) List(bucket Bucket, to interface{}) error {

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
