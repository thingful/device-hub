// Copyright © 2017 thingful

package mocks

import hub "github.com/thingful/device-hub"

// Channel allows you to create a mock Channel
// Implements hub.Channel
type Channel struct {
	ErrorChannel   chan error
	MessageChannel chan hub.Message
	Closer         func() error
}

func (m *Channel) Errors() chan error {
	return m.ErrorChannel
}

func (m *Channel) Out() chan hub.Message {
	return m.MessageChannel
}

func (m *Channel) Close() error {
	return m.Closer()
}

// Endpoint allows to create a mock Endpoint
// Implements hub.Endpoint
type Endpoint struct {
	Error error
}

func (e *Endpoint) Write(message hub.Message) error {
	return e.Error
}
