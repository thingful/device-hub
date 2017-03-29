package pipe

import (
	"os"
	"time"

	hub "github.com/thingful/device-hub"
)

const (
	URI_NAME_KEY  = "uri"
	TIMESTAMP_KEY = "received_at"
	PROTOCOL_KEY  = "protocol"
	HOSTNAME_KEY  = "host"
)

func newHubMessage(payload []byte, protocol, uri string) hub.Message {

	host, err := os.Hostname()

	if err != nil {
		host = "UNKNOWN"
	}

	m := hub.Message{
		Payload: payload,
		Metadata: map[string]interface{}{
			URI_NAME_KEY:  uri,
			PROTOCOL_KEY:  protocol,
			TIMESTAMP_KEY: time.Now().UTC(),
			HOSTNAME_KEY:  host,
		},
	}
	return m
}
