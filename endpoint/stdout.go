package endpoint

import (
	"encoding/json"
	"fmt"

	hub "github.com/thingful/device-hub"
)

func NewStdOutEndpoint(prettyPrint bool) stdout {
	return stdout{
		prettyPrint: prettyPrint,
	}
}

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
