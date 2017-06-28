// Copyright © 2017 thingful

package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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

func (f *fileStore) Insert(bucket bucket, uid []byte, data interface{}) error {

	f.lock.Lock()
	defer f.lock.Unlock()

	p := path.Join(f.path, string(bucket.name), string(uid))

	// check file exists
	exists, err := exists(p)

	if exists {
		return ErrItemAlreadyExists
	}

	// ensure the full folder path exists
	folder := path.Dir(p)
	err = os.MkdirAll(folder, os.ModePerm)

	if err != nil {
		return err
	}

	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	file, err := os.Create(p)

	if err != nil {
		return err
	}

	_, err = file.Write(buf)
	return err
}

func (f *fileStore) Delete(bucket bucket, uid []byte) error {

	f.lock.Lock()
	defer f.lock.Unlock()

	p := path.Join(f.path, string(bucket.name), string(uid))

	_, err := os.Stat(p)

	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
	}

	return os.Remove(p)
}

func (f *fileStore) One(bucket bucket, uid []byte, out interface{}) error {

	f.lock.RLock()
	defer f.lock.RUnlock()

	p := path.Join(f.path, string(bucket.name), string(uid))

	bytes, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &out)
}

func (f *fileStore) List(bucket bucket, to interface{}) error {

	f.lock.RLock()
	defer f.lock.RUnlock()

	folder := path.Join(f.path, string(bucket.name))
	listing, err := ioutil.ReadDir(folder)

	if err != nil {
		return err
	}

	// create a collector to retrieve the bucket contents
	c := func() (map[string][]byte, error) {

		list := map[string][]byte{}

		for _, file := range listing {

			fullPath := path.Join(folder, file.Name())

			visit := func(path string, f os.FileInfo, err error) error {

				if !f.IsDir() {

					bytes, err := ioutil.ReadFile(path)
					if err != nil {
						return err
					}

					list[file.Name()] = bytes
				}

				return nil
			}
			err := filepath.Walk(fullPath, visit)

			if err != nil {
				return list, err
			}
		}
		return list, err
	}

	return deserialiseCollection(to, c)
}

func (f *fileStore) Close() error {
	return nil
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
