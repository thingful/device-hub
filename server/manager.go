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

type statistics struct {
	Total  uint64
	Errors uint64
	OK     uint64
}

func NewEndpointManager(ctx context.Context) (*manager, error) {
	/*
		c := &config.Configuration{}
		pipes := map[string]*pipe{}

		for _, p := range c.Pipes {

			pipe, err := newRuntimeInstance(p, c)

			if err != nil {
				return nil, err
			}

			pipes[p.Uri] = pipe

		}
	*/
	return &manager{
		pipes: map[string]*pipe{},
		ctx:   ctx,
		//		conf:  c,
	}, nil
}

// newRuntimeInstance turns a config.Pipe into an managed runtime pipe
func newRuntimeInstance(p config.Pipe, c *config.Configuration) (*pipe, error) {

	found, listenerConf := c.Listeners.FindByUID(p.Listener)

	if !found {
		return nil, fmt.Errorf("listener with UID %s not found", p.Listener)
	}

	endpoints := []config.Endpoint{}

	for e := range p.Endpoints {

		found, endpointConf := c.Endpoints.FindByUID(p.Endpoints[e])

		if !found {
			return nil, fmt.Errorf("endpoint with UID %s not found", p.Endpoints[e])
		}

		endpoints = append(endpoints, endpointConf)

	}

	found, profile := c.Profiles.FindByUID(p.Profile)

	if !found {
		return nil, fmt.Errorf("profile with UID %s not found", p.Profile)
	}

	return &pipe{
		Uri:       p.Uri,
		Listener:  listenerConf,
		Endpoints: endpoints,
		Profile:   profile,
		State:     UNKNOWN,
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

	m.Lock()
	defer m.Unlock()

	r := []pipe{}

	for _, p := range m.pipes {
		r = append(r, *p)
	}

	return r
}

func (m *manager) DeletePipeByUID(uri string) error {

	if uri == "" {
		return errors.New("pipe uri not supplied")
	}

	m.Lock()
	defer m.Unlock()

	p, found := m.pipes[uri]
	if !found {
		return fmt.Errorf("pipe with uri : %s not found", uri)
	}

	p.cancel()

	delete(m.pipes, uri)

	// TODO : delete pipe from configuration

	return nil
}

func (m *manager) AddPipe(uri, profile, listener string, endpoints []string) error {

	if uri == "" {
		return errors.New("pipe uri not supplied")
	}
	if profile == "" {
		return errors.New("pipe profile not supplied")
	}
	if listener == "" {
		return errors.New("pipe listener not supplied")
	}
	if len(endpoints) == 0 {
		return errors.New("no endpoints supplied")
	}

	endpointUIDs := []config.UID{}

	for _, e := range endpoints {
		endpointUIDs = append(endpointUIDs, config.UID(e))
	}

	configPipe := config.Pipe{
		Uri:       uri,
		Profile:   config.UID(profile),
		Listener:  config.UID(listener),
		Endpoints: endpointUIDs,
	}

	pipe, err := newRuntimeInstance(configPipe, m.conf)

	if err != nil {
		return err
	}

	// important to ensure this lock is released on any return
	m.Lock()

	_, alreadyExists := m.pipes[uri]

	if alreadyExists {
		m.Unlock()
		return fmt.Errorf("pipe with uri %s already exists", uri)
	}

	m.pipes[uri] = pipe
	m.conf.Pipes = append(m.conf.Pipes, configPipe)

	m.Unlock()

	return m.Start()
}
