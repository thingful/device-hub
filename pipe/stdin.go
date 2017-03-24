package pipe

import (
	"bufio"
	"context"
	"errors"
	"os"

	hub "github.com/thingful/device-hub"
)

func NewStdInChannel(cancel context.CancelFunc) (Channel, error) {

	errors := make(chan error)
	out := make(chan hub.Input)

	channel := &stdinChannel{
		cancel: cancel,
		errors: errors,
		out:    out,
	}

	go channel.next()
	return channel, nil
}

type stdinChannel struct {
	errors chan error
	out    chan hub.Input
	cancel context.CancelFunc
}

// Errors returns a channel of errors
func (s stdinChannel) Errors() chan error {
	return s.errors
}

// Out returns a channel of expando.Input
func (s stdinChannel) Out() chan hub.Input {
	return s.out
}

func (s stdinChannel) next() {

	contents, err := getInputFromStdIn()

	if err != nil {
		s.errors <- err
	} else {
		s.out <- hub.Input{Payload: contents}
	}
	s.cancel()
}

// if we are being piped some input return it else error
func getInputFromStdIn() ([]byte, error) {

	fi, err := os.Stdin.Stat()

	if err != nil {
		return []byte{}, err
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		return []byte{}, errors.New("input expected from stdin")
	}

	reader := bufio.NewReader(os.Stdin)

	return reader.ReadBytes('\n')
}
