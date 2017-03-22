package pipe

import hub "github.com/thingful/device-hub"

type Broker interface {
	Channel() (Channel, error)
	Close() error
}

type Channel interface {
	Errors() chan error
	Out() chan hub.Input
}
