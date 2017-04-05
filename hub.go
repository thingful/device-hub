// Copyright © 2017 thingful

package hub

// Listener encapsulates various transports e.g. MQTT, HTTP creating a stream of Inputs
type Listener interface {
	NewChannel(string) (Channel, error)
	Close() error
}

// Channel exposes errors and Message channels
type Channel interface {
	Errors() chan error
	Out() chan Message
}

// Message contains a Payload, a processed Output and any Metadata collected
type Message struct {
	Payload  []byte                 `json:"payload"`
	Output   interface{}            `json:"output"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Endpoint takes a processed message and forwards to another service e.g. an HTTP endpoint, Kafka etc
type Endpoint interface {
	Write(message Message) error
}

var (
	// SourceVersion is set via the makefile
	SourceVersion = "DEVELOPMENT"
)
