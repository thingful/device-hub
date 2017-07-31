// Copyright © 2017 thingful

package listener

import (
	"encoding/hex"
	"log"
	"os"
	"time"

	"crypto/sha256"

	"github.com/rs/xid"

	hub "github.com/thingful/device-hub"
)

func newHubMessage(payload []byte, protocol, uri string) hub.Message {

	host, err := os.Hostname()
	if err != nil {
		host = "UNKNOWN"
	}
	hasher := sha256.New()
	if _, err := hasher.Write(payload); err != nil {
		log.Printf("failed to write to hash in message: %s", err.Error())
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	m := hub.Message{
		Payload: payload,
		Metadata: map[string]interface{}{
			hub.SHA256_SUM:         hash,
			hub.PIPE_URI_NAME_KEY:  uri,
			hub.PIPE_PROTOCOL_KEY:  protocol,
			hub.PIPE_TIMESTAMP_KEY: time.Now().UTC(),
			hub.HOSTNAME_KEY:       host,
			hub.MESSAGE_ID:         xid.New().String(),
		},
	}
	return m
}
