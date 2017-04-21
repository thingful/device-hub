// Copyright © 2017 thingful

package server

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

type store struct {
	db *bolt.DB
}

type bucket []byte

var (
	endpoints = bucket([]byte("endpoints"))
	listeners = bucket([]byte("listeners"))
)

func NewStore(db *bolt.DB) (*store, error) {

	err := db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(endpoints)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists(listeners)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &store{
		db: db,
	}, nil
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

func (s *store) Get(b bucket, uid string, out interface{}) error {

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

func (s *store) List(bucket bucket, out []interface{}) error {

	err := s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)

		b.ForEach(func(k, v []byte) error {

			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
