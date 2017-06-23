// Copyright © 2017 thingful

package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"reflect"
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

	// TODO : check for existing item

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

	ref := reflect.ValueOf(to)

	if ref.Kind() != reflect.Ptr || reflect.Indirect(ref).Kind() != reflect.Slice {
		return ErrSlicePtrNeeded
	}

	folder := path.Join(f.path, string(bucket.name))
	listing, err := ioutil.ReadDir(folder)

	if err != nil {
		return err
	}

	list := map[string][]byte{}

	for _, file := range listing {

		fullPath := path.Join(folder, file.Name())

		bytes, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return err
		}

		list[file.Name()] = bytes
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
