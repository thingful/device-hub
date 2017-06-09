// Copyright © 2017 thingful

package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
)

// Store is the entry point for the boltdb
type Store struct {
	db *bolt.DB
}

var (
	ErrSlicePtrNeeded = errors.New("slice ptr needed")

	ErrNotFound = errors.New("not found")
)

// NewStore returns a initilised Store instance
func NewStore(db *bolt.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) MustCreateBuckets(buckets []bucket) {

	err := s.db.Update(func(tx *bolt.Tx) error {

		for _, b := range buckets {

			_, err := tx.CreateBucketIfNotExists(b.name)
			if err != nil {
				panic(fmt.Sprintf("create bucket failed : %s", err))
			}

		}
		return nil
	})

	if err != nil {
		panic(err)
	}

}

func (s *Store) Insert(bucket bucket, uid []byte, data interface{}) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket.name)

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

func (s *Store) Update(bucket bucket, uid []byte, data interface{}) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket.name)

		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return b.Put(uid, buf)

	})
	return err
}

func (s *Store) Delete(bucket bucket, uid []byte) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket.name)

		return b.Delete(uid)
	})

	return err
}

func (s *Store) One(bucket bucket, uid []byte, out interface{}) error {

	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket.name)

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

func (s *Store) List(bucket bucket, to interface{}) error {

	ref := reflect.ValueOf(to)

	if ref.Kind() != reflect.Ptr || reflect.Indirect(ref).Kind() != reflect.Slice {
		return ErrSlicePtrNeeded
	}

	list := map[string][]byte{}
	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket.name)
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
