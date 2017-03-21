package pipe

import "github.com/thingful/expando"

type NoOpChannel struct {
}

func (NoOpChannel) Errors() chan error {
	return make(chan error)
}

func (NoOpChannel) Out() chan expando.Input {
	return make(chan expando.Input)
}
