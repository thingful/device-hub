package hub

// Listener encapsultates various transports e.g. MQTT, HTTP creating a stream of Inputs
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
