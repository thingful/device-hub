// Copyright © 2017 thingful

package runtime

import (
	"context"
	"errors"
	"sync"
	"testing"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/mocks"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/store"
	"github.com/thingful/device-hub/utils"

	"github.com/stretchr/testify/assert"
)

func TestLoopCancelledAndPipeStoppedOnContextDone(t *testing.T) {

	t.Parallel()

	ctx, closer := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	called := false

	mock := &mocks.Channel{
		Closer: func() error {

			called = true
			wg.Done()
			return nil
		},
		ErrorChannel:   make(chan error),
		MessageChannel: make(chan hub.Message),
	}

	pipe := &Pipe{}

	wg.Add(1)

	go loop(ctx, pipe, nil, map[string]hub.Endpoint{}, mock, utils.NewNoOpLogger(), map[string]string{})
	closer()

	wg.Wait()

	assert.True(t, called)
	assert.Equal(t, proto.Pipe_STOPPED, pipe.State)

}

func TestStatisticsOnChannelError(t *testing.T) {

	t.Parallel()

	ctx := context.Background()

	errorChannel := make(chan error)

	mock := &mocks.Channel{
		ErrorChannel:   errorChannel,
		MessageChannel: make(chan hub.Message),
	}

	pipe := newRuntimePipe(store.Pipe{})

	go loop(ctx, pipe, nil, map[string]hub.Endpoint{}, mock, utils.NewNoOpLogger(), map[string]string{})

	errorChannel <- errors.New("boo!")

	assert.Equal(t, uint64(1), pipe.Statistics.Received.Total)
	assert.Equal(t, uint64(1), pipe.Statistics.Received.Errors)
	assert.Equal(t, uint64(0), pipe.Statistics.Received.Ok)

}

func TestStatisticsOnChannelOut(t *testing.T) {

	t.Parallel()

	ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(3)

	messageChannel := make(chan hub.Message)

	mock := &mocks.Channel{
		ErrorChannel:   make(chan error),
		MessageChannel: messageChannel,
		Closer: func() error {
			defer wg.Done()
			return nil
		},
	}

	pipe := newRuntimePipe(store.Pipe{})

	endpoints := map[string]hub.Endpoint{
		"ok": &mocks.Endpoint{
			Writer: func(hub.Message) error {
				defer wg.Done()
				return nil
			},
		},
		"error": &mocks.Endpoint{
			Writer: func(hub.Message) error {
				defer wg.Done()
				return errors.New("boo")
			},
		},
	}

	go loop(ctx, pipe, nil, endpoints, mock, utils.NewNoOpLogger(), map[string]string{})

	message := hub.Message{
		Payload:  []byte("hello"),
		Metadata: map[string]interface{}{},
	}
	messageChannel <- message

	mock.Close()
	wg.Wait()

	assert.Equal(t, uint64(1), pipe.Statistics.Processed.Ok)
	assert.Equal(t, uint64(1), pipe.Statistics.Processed.Total)
	assert.Equal(t, uint64(0), pipe.Statistics.Processed.Errors)

	assert.Equal(t, uint64(1), pipe.Statistics.Received.Ok)
	assert.Equal(t, uint64(1), pipe.Statistics.Received.Total)
	assert.Equal(t, uint64(0), pipe.Statistics.Processed.Errors)

	assert.Equal(t, uint64(1), pipe.Statistics.Sent["ok"].Ok)
	assert.Equal(t, uint64(1), pipe.Statistics.Sent["ok"].Total)
	assert.Equal(t, uint64(0), pipe.Statistics.Sent["ok"].Errors)

	assert.Equal(t, uint64(0), pipe.Statistics.Sent["error"].Ok)
	assert.Equal(t, uint64(1), pipe.Statistics.Sent["error"].Total)
	assert.Equal(t, uint64(1), pipe.Statistics.Sent["error"].Errors)

}
