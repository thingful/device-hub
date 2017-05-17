// Copyright © 2017 thingful

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/store"
)

type manager struct {
	Repository *store.Repository
	ctx        context.Context
	pipes      map[string]*pipe
	sync.RWMutex
}

// pipe is an instance of a pipe containing runtime state information e.g. stats, state
type pipe struct {
	store.Pipe

	State   proto.Pipe_State
	Started time.Time

	MessageStatistics proto.Statistics

	cancel context.CancelFunc

	// TODO : add last error, debug etc
}

type pipePredicate func(*pipe) bool

func NewEndpointManager(ctx context.Context, repository *store.Repository) (*manager, error) {

	// load any existing pipes from the database to
	// serve as the initial running state
	state, err := repository.Pipes.List()

	if err != nil {
		return nil, err
	}

	pipes := map[string]*pipe{}

	for _, p := range state {
		pipes[p.Uri] = &pipe{Pipe: p}
	}

	return &manager{
		Repository: repository,
		pipes:      pipes,
		ctx:        ctx,
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

			// TODO : fix error statisitics capturing
			// maybe separate in to processed, sending, receiving

			if err != nil {
				p.MessageStatistics.Errors++
				log.Println(err)
			} else {
				p.MessageStatistics.Ok++
			}

			output.Metadata[hub.PROFILE_NAME_KEY] = p.Profile.Name
			output.Metadata[hub.PROFILE_VERSION_KEY] = p.Profile.Version
			output.Metadata[hub.RUNTIME_VERSION_KEY] = hub.SourceVersion

			output.Schema = p.Profile.Schema

			for e := range endpoints {

				// TODO : do something more useful with this error
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

func (m *manager) Any(f pipePredicate) bool {

	m.RLock()
	defer m.RUnlock()

	for _, p := range m.pipes {
		if f(p) {
			return true
		}
	}
	return false
}

func (m *manager) DeletePipe(f pipePredicate) error {

	m.Lock()
	defer m.Unlock()

	uri := ""
	for _, p := range m.pipes {
		if f(p) {
			uri = p.Uri
			break
		}
	}

	p, found := m.pipes[uri]
	if !found {
		return nil
	}

	p.cancel()

	delete(m.pipes, uri)
	// TODO : keeps buffer of recently deleted pipes

	err := m.Repository.Pipes.Delete(uri)

	if err != nil {
		return err
	}

	return nil
}

func (m *manager) StartPipe(uri, listenerUID, profileUID string, endpointUIDs []string) error {

	listener, err := m.Repository.Listeners.One(listenerUID)

	if err != nil {
		if err == store.ErrNotFound {
			return fmt.Errorf("listener with uid : %s not found", listenerUID)
		}
		return err
	}

	temp, err := m.Repository.Profiles.One(profileUID)

	if err != nil {
		if err == store.ErrNotFound {
			return fmt.Errorf("profile with uid : %s not found", profileUID)
		}
		return err
	}

	profile, err := profileFromEntity(temp)

	if err != nil {
		return err
	}

	endpoints, err := m.Repository.Endpoints.Many(endpointUIDs)

	if err != nil {
		if err == store.ErrNotFound {
			return fmt.Errorf("endpoint with uids : %v not found", endpointUIDs)
		}
		return err
	}

	pipeconf := store.Pipe{
		Uri:       uri,
		Listener:  listener,
		Endpoints: endpoints,
		Profile:   *profile,
	}

	runtimepipe := &pipe{
		Pipe:  pipeconf,
		State: proto.Pipe_UNKNOWN,
	}

	err = m.Repository.Pipes.CreateOrUpdate(pipeconf)

	if err != nil {
		return err
	}

	err = m.addPipe(runtimepipe)

	if err != nil {

		deleteError := m.Repository.Pipes.Delete(pipeconf.Uri)

		if deleteError != nil {
			// TODO : review this!
			panic(deleteError)
		}
		return err
	}

	return nil
}

func profileFromEntity(entity *proto.Entity) (*store.Profile, error) {
	// TODO : give a monkeys about validation

	schema := map[string]interface{}{}

	err := json.Unmarshal([]byte(entity.Configuration["schema"]), &schema)

	if err != nil {
		return nil, err
	}

	return &store.Profile{
		Uid:         entity.Uid,
		Name:        entity.Configuration["profile-name"],
		Description: entity.Configuration["profile-description"],
		Version:     entity.Configuration["profile-version"],
		Schema:      schema,
		Script: engine.Script{
			Main:     entity.Configuration["script-main"],
			Runtime:  engine.Runtime(entity.Configuration["script-runtime"]),
			Input:    engine.InputType(entity.Configuration["script-input"]),
			Contents: entity.Configuration["script-contents"],
		},
	}, nil
}

func (m *manager) addPipe(pipe *pipe) error {
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
