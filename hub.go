// Copyright © 2017 thingful

package hub

// Listener encapsulates various transports e.g. MQTT, HTTP
type Listener interface {
	NewChannel(string) (Channel, error)
	Close() error
}

// Channel exposes errors and Message channels
type Channel interface {
	Errors() chan error
	Out() chan Message
	Close() error
}

// Message contains a Payload, a processed Output and any Metadata collected
type Message struct {
	Payload  []byte                 `json:"payload"`
	Output   interface{}            `json:"output"`
	Schema   map[string]interface{} `json:"schema"`
	Metadata map[string]interface{} `json:"metadata"`
	Tags     map[string]string      `json:"tags,omitempty"`
}

// Endpoint takes a processed message and forwards to another service e.g. an HTTP endpoint, Kafka etc
type Endpoint interface {
	Write(message Message) error
}
