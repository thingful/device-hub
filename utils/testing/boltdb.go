// Copyright © 2017 thingful
package testing

import (
	"io/ioutil"
	"os"

	"github.com/boltdb/bolt"
)

type boltDBConnection struct {
	*bolt.DB
	path string
}

// DialBoltDB returns an isolated connection to a boltDB instance. Suitable for parallel testing.
func DialBoltDB() (*boltDBConnection, error) {

	// Generate temporary filename.
	f, err := ioutil.TempFile("", "bolt-test-")
	if err != nil {
		return nil, err
	}
	f.Close()

	db, err := bolt.Open(f.Name(), 0600, nil)

	if err != nil {
		return nil, err
	}
	return &boltDBConnection{
		path: f.Name(),
		DB:   db,
	}, nil
}

// MustDialBoltDB returns an isolated connection or panics.
func MustDialBoltDB() *boltDBConnection {

	conn, err := DialBoltDB()

	if err != nil {
		panic(err)
	}
	return conn
}

// Close cleans up the boltdb connection.
func (c *boltDBConnection) Close() error {
	defer os.Remove(c.path)
	return c.DB.Close()

}

// MustClose clears up the connection or panics.
func (c *boltDBConnection) MustClose() {
	err := c.Close()
	if err != nil {
		panic(err)
	}
}
