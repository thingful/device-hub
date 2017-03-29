package pipe

import (
	"encoding/json"
	"fmt"

	hub "github.com/thingful/device-hub"
)

func WriteToStdOut(message hub.Message) error {

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil

}
