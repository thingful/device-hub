// Copyright © 2017 thingful

package store

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
)

// boltDBStore implements Storer
type boltDBStore struct {
	db *bolt.DB
}

// NewStore returns a initilised Storer instance using BoltDB as the backing store
func NewBoltDBStore(db *bolt.DB) Storer {
	return &boltDBStore{
		db: db,
	}
}

func (s *boltDBStore) MustCreateBuckets(buckets []bucket) {

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

func (s *boltDBStore) Insert(bucket bucket, uid []byte, data interface{}) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket.name)

		existing := b.Get(uid)

		if len(existing) > 0 {
			return ErrItemAlreadyExists
		}

		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return b.Put([]byte(uid), buf)

	})
	return err
}

func (s *boltDBStore) Delete(bucket bucket, uid []byte) error {

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket.name)

		return b.Delete(uid)
	})

	return err
}

func (s *boltDBStore) One(bucket bucket, uid []byte, out interface{}) error {

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

type m func() (map[string][]byte, error)

func fff(to interface{}, n m) error {
	ref := reflect.ValueOf(to)

	if ref.Kind() != reflect.Ptr || reflect.Indirect(ref).Kind() != reflect.Slice {
		return ErrSlicePtrNeeded
	}

	list, err := n()

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

func (s *boltDBStore) List(bucket bucket, to interface{}) error {

	f := func() (map[string][]byte, error) {
		list := map[string][]byte{}
		err := s.db.View(func(tx *bolt.Tx) error {

			b := tx.Bucket(bucket.name)
			b.ForEach(func(k, v []byte) error {
				list[string(k)] = v
				return nil
			})

			return nil
		})

		return list, err

	}

	return fff(to, f)
}

func (s *boltDBStore) Close() error {
	return s.db.Close()
}
