// Copyright © 2017 thingful

package store

import (
	"encoding/json"
	"os"
	"path"
	"sync"
)

// fileStore implements Storer
type fileStore struct {
	path string
	lock sync.RWMutex
}

func NewFileStore(path string) Storer {
	return &fileStore{
		path: path,
	}
}

func (f *fileStore) MustCreateBuckets(buckets []bucket) {

	f.lock.Lock()
	defer f.lock.Unlock()

	for _, b := range buckets {

		p := path.Join(f.path, string(b.name))

		_, err := os.Stat(p)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(p, os.ModePerm)

				if err != nil {
					panic(err)
				}

			} else {
				panic(err)
			}

		}
	}
}

func (f *fileStore) InsertOrUpdate(bucket bucket, uid []byte, data interface{}) error {

	f.lock.Lock()
	defer f.lock.Unlock()

	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	p := path.Join(f.path, string(bucket.name), string(uid))

	// ensure the folder exists
	folder := path.Dir(p)
	err = os.MkdirAll(folder, os.ModePerm)

	if err != nil {
		return err
	}

	file, err := os.Create(p)

	if err != nil {
		return err
	}

	_, err = file.Write(buf)

	if err != nil {
		return err
	}

	return nil
}

func (f *fileStore) Delete(bucket bucket, uid []byte) error {

	f.lock.Lock()
	defer f.lock.Unlock()

	return nil
}

func (f *fileStore) One(bucket bucket, uid []byte, out interface{}) error {

	f.lock.RLock()
	defer f.lock.RUnlock()

	return nil
}

func (f *fileStore) List(bucket bucket, to interface{}) error {

	f.lock.RLock()
	defer f.lock.RUnlock()

	return nil
}
