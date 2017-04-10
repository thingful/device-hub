// Copyright © 2017 thingful

package listener

import hub "github.com/thingful/device-hub"

// defaultChannel is an implementation of Channel
type defaultChannel struct {
	errors chan error
	out    chan hub.Message
	close  func() error
}

func (m defaultChannel) Errors() chan error {
	return m.errors
}

func (m defaultChannel) Out() chan hub.Message {
	return m.out
}

func (m defaultChannel) Close() error {

	err := m.close()

	close(m.out)
	close(m.errors)

	return err
}
