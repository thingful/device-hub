package pipe

import (
	"bufio"
	"context"
	"errors"
	"os"

	"bitbucket.org/tsetsova/decode-prototype/hub/expando"
)

func FromStdIn(ctx context.Context) *stdin {

	channel := make(chan expando.Input)
	errors := make(chan error)

	return &stdin{
		ctx:     ctx,
		channel: channel,
		errors:  errors,
	}
}

type stdin struct {
	ctx     context.Context
	channel chan expando.Input
	errors  chan error
}

func (s *stdin) Channel() chan expando.Input {
	return s.channel
}

func (s *stdin) Errors() chan error {
	return s.errors
}

func (s *stdin) Next() {

	contents, err := getInputFromStdIn()

	if err != nil {
		s.errors <- err
	} else {
		s.channel <- expando.Input{Payload: []byte(contents)}
	}
}

// if we are being piped some input return it else error
func getInputFromStdIn() (string, error) {

	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return "", errors.New("input expected from stdin")
	}

	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}
