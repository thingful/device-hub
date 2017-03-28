package profile

import (
	"encoding/json"
	"fmt"
	"testing"

	hub "github.com/thingful/device-hub"
)

func TestX(t *testing.T) {

	conf := Configuration{
		Listeners: []Endpoint{
			Endpoint{
				Type: "http",
				UID:  "http-one",
				Configuration: map[string]interface{}{
					"BindingAddress": ":8085",
				},
			},
		},
		Endpoints: []Endpoint{
			Endpoint{
				Type: "stdout",
				UID:  "stdout-one",
			},
		},
		Profiles: []Profile{
			Profile{
				Name:        "some-device-profile",
				Description: "blah-blah",
				UID:         "xxx",
				Version:     "one",
				Script: hub.Script{
					//Name:     "script-one",
					Main:     "decode",
					Runtime:  hub.Javascript,
					Input:    hub.JSON,
					Contents: "function decode(x) { return x };",
				},
			},
		},
		Pipes: []Pipe{
			Pipe{Uri: "/a",
				Profile:  "xxx",
				Listener: "http-one",
				Endpoint: "stdout-one"},
		},
	}

	bytes, _ := json.MarshalIndent(conf, "", "  ")
	fmt.Println(string(bytes))

}
