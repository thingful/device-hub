// Copyright © 2017 thingful

package runtime

import (
	"context"
	"errors"
	"sync"
	"testing"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/store"

	"github.com/stretchr/testify/assert"
)

type mockChannel struct {
	errors   chan error
	messages chan hub.Message
	closer   func() error
}

func (m *mockChannel) Errors() chan error {
	return m.errors
}

func (m *mockChannel) Out() chan hub.Message {
	return m.messages
}

func (m *mockChannel) Close() error {
	return m.closer()
}

func TestLoopCancelledAndPipeStoppedOnContextDone(t *testing.T) {

	t.Parallel()

	ctx, closer := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	called := false

	mock := &mockChannel{
		closer: func() error {

			called = true
			wg.Done()
			return nil
		},
		errors:   make(chan error),
		messages: make(chan hub.Message),
	}

	pipe := &Pipe{}

	wg.Add(1)

	go loop(ctx, pipe, nil, map[string]hub.Endpoint{}, mock, map[string]string{})
	closer()

	wg.Wait()

	assert.True(t, called)
	assert.Equal(t, proto.Pipe_STOPPED, pipe.State)

}

func TestStatisticsOnChannelError(t *testing.T) {

	t.Parallel()

	ctx := context.Background()

	//	wg := sync.WaitGroup{}

	errorChannel := make(chan error)

	mock := &mockChannel{
		errors:   errorChannel,
		messages: make(chan hub.Message),
	}

	pipe := newRuntimePipe(store.Pipe{})

	go loop(ctx, pipe, nil, map[string]hub.Endpoint{}, mock, map[string]string{})

	errorChannel <- errors.New("boo!")

	assert.Equal(t, uint64(1), pipe.Statistics.Received.Total)
	assert.Equal(t, uint64(1), pipe.Statistics.Received.Errors)
	assert.Equal(t, uint64(0), pipe.Statistics.Received.Ok)

}
