// Copyright © 2017 thingful

package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/config"
	"github.com/thingful/device-hub/engine"
)

type manager struct {
	ctx   context.Context
	pipes map[string]*pipe
	conf  *config.Configuration
	sync.RWMutex
}

// state tracks the known state of a runtime pipe
type state string

const (
	UNKNOWN = state("UNKNOWN")
	RUNNING = state("RUNNING")
	STOPPED = state("STOPPED")
	ERRORED = state("ERRORED")
)

// pipe is an instance of a pipe containing runtime state information e.g. stats, state
type pipe struct {
	Listener  config.Endpoint
	Endpoints []config.Endpoint
	Profile   config.Profile
	Uri       string

	State   state
	Started time.Time

	MessageStatistics statistics

	cancel context.CancelFunc

	// TODO : add last error, debug etc
}

type Pipes []*pipe

type statistics struct {
	Total  uint64
	Errors uint64
	OK     uint64
}

func NewEndpointManager(ctx context.Context, state Pipes) (*manager, error) {

	pipes := map[string]*pipe{}

	for _, p := range state {
		pipes[p.Uri] = p
	}

	return &manager{
		pipes: pipes,
		ctx:   ctx,
	}, nil
}

// Start either ensures everything is running or errors
func (m *manager) Start() error {

	m.Lock()
	defer m.Unlock()

	for n, p := range m.pipes {

		if p.State != RUNNING {

			listener, err := hub.ListenerByName(string(p.Listener.UID), p.Listener.Kind, p.Listener.Configuration)

			if err != nil {
				return err
			}

			endpoints := []hub.Endpoint{}

			for _, e := range p.Endpoints {

				newendpoint, err := hub.EndpointByName(string(e.UID), e.Kind, e.Configuration)

				if err != nil {
					return err
				}

				endpoints = append(endpoints, newendpoint)

			}

			channel, err := listener.NewChannel(p.Uri)

			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(m.ctx)

			pp := m.pipes[n]

			go m.startOne(ctx, pp, listener, endpoints, channel)
			pp.cancel = cancel
			pp.State = RUNNING
			pp.Started = time.Now().UTC()
		}
	}
	return nil
}

func (m *manager) startOne(ctx context.Context, p *pipe, listener hub.Listener, endpoints []hub.Endpoint, channel hub.Channel) {

	scripter := engine.New()

	for {

		select {

		case <-ctx.Done():
			p.State = STOPPED
			err := channel.Close()

			if err != nil {
				log.Println(err)
			}

			return

		case err := <-channel.Errors():

			p.MessageStatistics.Total++
			p.MessageStatistics.Errors++
			log.Println(err)

		case input := <-channel.Out():

			output, err := scripter.Execute(p.Profile.Script, input)

			p.MessageStatistics.Total++

			if err != nil {
				p.MessageStatistics.Errors++
				log.Println(err)
			} else {
				p.MessageStatistics.OK++
			}

			output.Metadata[hub.PROFILE_NAME_KEY] = p.Profile.Name
			output.Metadata[hub.PROFILE_VERSION_KEY] = p.Profile.Version
			output.Metadata[hub.RUNTIME_VERSION_KEY] = hub.SourceVersion

			for e := range endpoints {

				err = endpoints[e].Write(output)

				if err != nil {
					log.Println(err)
				}

			}
		}
	}
}

func (m *manager) List() []pipe {

	m.RLock()
	defer m.RUnlock()

	r := []pipe{}

	for _, p := range m.pipes {
		r = append(r, *p)
	}

	return r
}

func (m *manager) DeletePipeByURI(uri string) error {

	if uri == "" {
		return errors.New("pipe uri not supplied")
	}

	m.Lock()
	defer m.Unlock()

	p, found := m.pipes[uri]
	if !found {
		return nil
	}

	p.cancel()

	delete(m.pipes, uri)
	fmt.Println(m.pipes)
	// TODO : keeps buffer of recently deleted pipes

	return nil
}

func (m *manager) AddPipe(pipe *pipe) error {
	m.Lock()

	_, alreadyExists := m.pipes[pipe.Uri]

	if alreadyExists {
		m.Unlock()
		return fmt.Errorf("pipe with uri %s already exists", pipe.Uri)
	}

	m.pipes[pipe.Uri] = pipe

	m.Unlock()
	return m.Start()
}
