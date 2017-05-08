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
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/utils"
)

type manager struct {
	ctx   context.Context
	pipes map[string]*pipe
	sync.RWMutex
}

type endpoint struct {
	Kind          string
	Type          string
	Uid           string
	Configuration utils.TypedMap
}

type profile struct {
	Uid         string
	Name        string
	Description string
	// TODO : make this a semantic triple
	Version string
	Script  engine.Script
}

// pipe is an instance of a pipe containing runtime state information e.g. stats, state
type pipe struct {
	Listener  endpoint
	Endpoints []endpoint
	Profile   profile
	Uri       string

	State   proto.Pipe_State
	Started time.Time

	MessageStatistics proto.Statistics

	cancel context.CancelFunc

	// TODO : add last error, debug etc
}

type Pipes []*pipe

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

		if p.State != proto.Pipe_RUNNING {

			listener, err := hub.ListenerByName(p.Listener.Uid, p.Listener.Kind, p.Listener.Configuration)

			if err != nil {
				return err
			}

			endpoints := []hub.Endpoint{}

			for _, e := range p.Endpoints {

				newendpoint, err := hub.EndpointByName(e.Uid, e.Kind, e.Configuration)

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
			pp.State = proto.Pipe_RUNNING
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
			p.State = proto.Pipe_STOPPED
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
				p.MessageStatistics.Ok++
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

			fmt.Print(p.MessageStatistics)
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
