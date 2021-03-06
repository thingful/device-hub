// Copyright © 2017 thingful

package endpoint

import (
	"encoding/json"
	"fmt"

	hub "github.com/thingful/device-hub"
)

type stdout struct {
	prettyPrint bool
}

func (s stdout) Write(message hub.Message) error {

	var bytes []byte
	var err error

	if s.prettyPrint {

		bytes, err = json.MarshalIndent(message, "", "    ")

	} else {

		bytes, err = json.Marshal(message)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil

}
