package pipe

import (
	"bufio"
	"context"
	"errors"
	"os"

	"github.com/thingful/expando"
)

// FromStdIn pipes data from stdin
func FromStdIn(cancel context.CancelFunc) *stdin {

	errors := make(chan error)

	return &stdin{
		cancel: cancel,
		errors: errors,
	}
}

type stdin struct {
	cancel context.CancelFunc
	errors chan error
}

// Channel returns a new channel to start processing messages
func (s *stdin) Channel() Channel {
	out := make(chan expando.Input)

	channel := stdinChannel{cancel: s.cancel,
		out:    out,
		errors: s.errors}

	go channel.next()

	return channel
}

func (s *stdin) Close() error {
	return nil
}

type stdinChannel struct {
	errors chan error
	out    chan expando.Input
	cancel context.CancelFunc
}

// Errors returns a channel of errors
func (s stdinChannel) Errors() chan error {
	return s.errors
}

// Out returns a channel of expando.Input
func (s stdinChannel) Out() chan expando.Input {
	return s.out
}

func (s stdinChannel) next() {

	contents, err := getInputFromStdIn()

	if err != nil {
		s.errors <- err
	} else {
		s.out <- expando.Input{Payload: contents}
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
