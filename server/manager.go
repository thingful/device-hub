// Copyright © 2017 thingful

package server

import (
	"context"
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
	pipes []pipe
	sync.RWMutex
}

type state string

const (
	UNKNOWN = state("UNKNOWN")
	RUNNING = state("RUNNING")
	STOPPED = state("STOPPED")
	ERRORED = state("ERRORED")
)

type pipe struct {
	Listener  config.Endpoint
	Endpoints []config.Endpoint
	Profile   config.Profile
	Uri       string

	State   state
	Started time.Time

	MessageStatistics statistics

	// TODO : add last error, debug etc
}

type statistics struct {
	Total  uint64
	Errors uint64
	OK     uint64
}

func NewEndpointManager(ctx context.Context, c *config.Configuration) (*manager, error) {

	pipes := []pipe{}

	for _, p := range c.Pipes {

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

		pipes = append(pipes, pipe{
			Uri:       p.Uri,
			Listener:  listenerConf,
			Endpoints: endpoints,
			Profile:   profile,
			State:     UNKNOWN})

	}

	return &manager{
		pipes: pipes,
		ctx:   ctx,
	}, nil
}

func (m *manager) Start() error {

	m.Lock()
	defer m.Unlock()

	for n, p := range m.pipes {

		if p.State != RUNNING {

			listener, err := hub.ListenerByName(string(p.Listener.UID), p.Listener.Type, p.Listener.Configuration)

			if err != nil {
				return err
			}

			endpoints := []hub.Endpoint{}

			for _, e := range p.Endpoints {

				newendpoint, err := hub.EndpointByName(string(e.UID), e.Type, e.Configuration)

				if err != nil {
					return err
				}

				endpoints = append(endpoints, newendpoint)

			}

			channel, err := listener.NewChannel(p.Uri)

			if err != nil {
				return err
			}

			go m.startOne(&m.pipes[n], listener, endpoints, channel)
			m.pipes[n].State = RUNNING
		}
	}
	return nil
}

func (m *manager) startOne(p *pipe, listener hub.Listener, endpoints []hub.Endpoint, channel hub.Channel) {

	scripter := engine.New()

	for {

		select {

		case <-m.ctx.Done():
			p.State = STOPPED
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
	return m.pipes
}
