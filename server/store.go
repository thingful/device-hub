// Copyright © 2017 thingful

package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

type store struct {
	db *bolt.DB
}

type bucket []byte

var (
	endPoints = bucket([]byte("endpoints"))
)

func NewStore(db *bolt.DB) (*store, error) {

	err := db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(endPoints)
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

func (s *store) Insert(bucket bucket, data interface{}) (string, error) {

	var uid string

	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket)

		id, err := b.NextSequence()

		if err != nil {
			return err
		}

		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}

		uid = fmt.Sprintf("%s-%d", bucket, id)

		return b.Put(itob(id), buf)

	})
	return uid, err

}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
