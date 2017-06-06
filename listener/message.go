// Copyright © 2017 thingful

package listener

import (
	"os"
	"time"

	hub "github.com/thingful/device-hub"
)

func NewHubMessage(payload []byte, protocol, uri string) hub.Message {

	host, err := os.Hostname()

	if err != nil {
		host = "UNKNOWN"
	}

	m := hub.Message{
		Payload: payload,
		Metadata: map[string]interface{}{
			hub.PIPE_URI_NAME_KEY:  uri,
			hub.PIPE_PROTOCOL_KEY:  protocol,
			hub.PIPE_TIMESTAMP_KEY: time.Now().UTC(),
			hub.HOSTNAME_KEY:       host,
		},
	}
	return m
}
