package pipe

import (
	hub "github.com/thingful/device-hub"
)

type NoOpChannel struct {
}

func (NoOpChannel) Errors() chan error {
	return make(chan error)
}

func (NoOpChannel) Out() chan hub.Input {
	return make(chan hub.Input)
}
