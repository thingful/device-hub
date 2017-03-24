package pipe

import hub "github.com/thingful/device-hub"

// Channel exposes errors and hub.Input channels
type Channel interface {
	Errors() chan error
	Out() chan hub.Input
}

type defaultChannel struct {
	errors chan error
	out    chan hub.Input
}

func (m defaultChannel) Errors() chan error {
	return m.errors
}

func (m defaultChannel) Out() chan hub.Input {
	return m.out
}
