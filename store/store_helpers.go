package store

import (
	"encoding/json"
	"reflect"
)

// collector retunrns a map of []byte arrays from the Storer implementations
type collector func() (map[string][]byte, error)

// deserialiseCollection executes a collector deserialising it into a passed in interface
func deserialiseCollection(to interface{}, c collector) error {
	ref := reflect.ValueOf(to)

	if ref.Kind() != reflect.Ptr || reflect.Indirect(ref).Kind() != reflect.Slice {
		return ErrSlicePtrNeeded
	}

	list, err := c()

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
