package endpoint

import (
	"encoding/json"
	"fmt"

	hub "github.com/thingful/device-hub"
)

func NewStdOutEndpoint() stdout {
	return stdout{}
}

type stdout struct{}

func (stdout) Write(message hub.Message) error {

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil

}
