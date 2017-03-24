package pipe

import (
	"bufio"
	"context"
	"errors"
	"os"

	hub "github.com/thingful/device-hub"
)

func NewStdInChannel(cancel context.CancelFunc) *stdinChannel {

	errors := make(chan error)
	out := make(chan hub.Input)

	channel := &stdinChannel{
		defaultChannel: defaultChannel{
			errors: errors,
			out:    out,
		}, cancel: cancel,
	}

	go channel.next()
	return channel
}

type stdinChannel struct {
	defaultChannel
	cancel context.CancelFunc
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
