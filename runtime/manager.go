// Copyright © 2017 thingful

package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/registry"
	"github.com/thingful/device-hub/store"
	"github.com/thingful/device-hub/utils"
)

type Options struct {
	GeoEnabled bool
	GeoLat     float64
	GeoLng     float64
}

// Manager holds the running instances of all the pipes
type Manager struct {
	Repository *store.Repository
	ctx        context.Context
	pipes      map[string]*Pipe
	sync.RWMutex
	register *registry.Registry
	logger   utils.Logger
	options  Options
}

// pipe holds runtime state information including various counters
// including messages received, message processed and messages sent
// Each counter contains the total, errors and ok
type Pipe struct {
	// wraps a store.Pipe
	store.Pipe

	// current state - see proto/devicehub.proto for details
	State proto.Pipe_State

	// capture the time started
	Started time.Time

	// keep runtime statistics
	// NOTE : statistics are not persisted or cumulative across shutdowns
	Statistics *proto.Statistics

	// allow the pipe to cancelled and give it an opportunity
	// to shut down nicely
	cancel context.CancelFunc
}

// newRuntimePipe returns an initialised runtime.Pipe
func newRuntimePipe(p store.Pipe) *Pipe {
	return &Pipe{
		Pipe:  p,
		State: proto.Pipe_UNKNOWN,
		Statistics: &proto.Statistics{
			Processed: &proto.Counters{},
			Received:  &proto.Counters{},
			Sent:      map[string]*proto.Counters{},
		},
	}
}

// PipePredicate is a function to facilitate predicating the collection
type PipePredicate func(*Pipe) bool

// NewEndpointManager returns a manager instance or an error
func NewEndpointManager(ctx context.Context,
	repository *store.Repository,
	registry *registry.Registry,
	logger utils.Logger,
	options Options) (*Manager, error) {

	// load any existing pipes from the database to
	// serve as the initial running state
	state, err := repository.Pipes.List()

	if err != nil {
		return nil, err
	}

	pipes := map[string]*Pipe{}

	for _, p := range state {
		pipes[p.Uri] = newRuntimePipe(p)
	}
	return &Manager{
		Repository: repository,
		pipes:      pipes,
		ctx:        ctx,
		register:   registry,
		logger:     logger,
		options:    options,
	}, nil
}

// Start either ensures everything is running or errors
// TODO: allow an allowance for pipe start up failure and mark as unstartable
func (m *Manager) Start() error {

	m.Lock()
	defer m.Unlock()

	for n, p := range m.pipes {

		if p.State != proto.Pipe_RUNNING {

			listener, err := m.register.ListenerByName(p.Listener.Uid, p.Listener.Kind, p.Listener.Configuration)

			if err != nil {
				return err
			}

			endpoints := map[string]hub.Endpoint{}

			for _, e := range p.Endpoints {

				newendpoint, err := m.register.EndpointByName(e.Uid, e.Kind, e.Configuration)

				if err != nil {
					return err
				}

				endpoints[e.Uid] = newendpoint

			}

			channel, err := listener.NewChannel(p.Uri)

			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(m.ctx)

			pp := m.pipes[n]

			go loop(ctx, pp, listener, endpoints, channel, m.logger, p.Tags, m.options)
			pp.cancel = cancel
			pp.State = proto.Pipe_RUNNING
			pp.Started = time.Now().UTC()
		}
	}
	return nil
}

// Status returns the set of known pipes
func (m *Manager) Status() []Pipe {

	m.RLock()
	defer m.RUnlock()

	r := []Pipe{}

	for _, p := range m.pipes {
		r = append(r, *p)
	}

	return r
}

// Any facilitates matching on a predicate
// returns on first match or iterates to completion
func (m *Manager) Any(f PipePredicate) bool {

	m.RLock()
	defer m.RUnlock()

	for _, p := range m.pipes {
		if f(p) {
			return true
		}
	}
	return false
}

// DeletePipe cancels and removes the first match
func (m *Manager) DeletePipe(f PipePredicate) error {

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

// StartPipe ensures all components can be found and the runtime information persisted before starting an instance
func (m *Manager) StartPipe(uri, listenerUID, profileUID string, endpointUIDs []string, tags map[string]string) error {

	listener, err := m.Repository.Listeners.One(listenerUID)

	if err != nil {
		if err == store.ErrNotFound {
			return fmt.Errorf("listener with uid : %s not found", listenerUID)
		}
		return err
	}

	temp, err := m.Repository.Profiles.One(profileUID)

	if err != nil {

		fmt.Println("err", err)
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
		Tags:      tags,
	}

	runtimepipe := newRuntimePipe(pipeconf)
	err = m.Repository.Pipes.Insert(pipeconf)

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

	if entity.Configuration["schema"] != "" {

		err := json.Unmarshal([]byte(entity.Configuration["schema"]), &schema)

		if err != nil {
			return nil, err
		}
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

func (m *Manager) addPipe(pipe *Pipe) error {
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
