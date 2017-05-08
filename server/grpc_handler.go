// Copyright © 2017 thingful

package server

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strings"

	hashids "github.com/speps/go-hashids"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/store"
	context "golang.org/x/net/context"
)

type handler struct {
	manager *manager
	store   *store.Store
}

func (s *handler) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateReply, error) {

	hash, err := hash(request)

	if err != nil {
		return &proto.CreateReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	var bucket store.Bucket

	switch strings.ToLower(request.Type) {
	case "listener":
		bucket = store.Listeners

		exists := hub.IsListenerRegistered(request.Kind)

		if !exists {
			return &proto.CreateReply{
				Ok:    false,
				Error: fmt.Sprintf("kind : %s not registered", request.Kind),
			}, nil

		}

	case "endpoint":
		bucket = store.Endpoints

		exists := hub.IsEndpointRegistered(request.Kind)

		if !exists {
			return &proto.CreateReply{
				Ok:    false,
				Error: fmt.Sprintf("kind : %s not registered", request.Kind),
			}, nil
		}
	case "profile":
		bucket = store.Profiles

		if request.Configuration["profile-name"] != "" {
			// TODO : consider adding version to the profile-name? Would be useful for
			// having multiple profiles running at the same time.
			hash = []byte(request.Configuration["profile-name"])
		}

	default:
		return &proto.CreateReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s not registered", request.Type),
		}, nil
	}

	err = s.store.Insert(bucket, hash, proto.Entity{
		Uid:           string(hash),
		Type:          request.Type,
		Kind:          request.Kind,
		Configuration: request.Configuration,
	})

	if err != nil {
		return &proto.CreateReply{
			Error: err.Error(),
		}, nil
	}

	return &proto.CreateReply{
		Uid: string(hash),
		Ok:  true,
	}, nil
}

func (s *handler) Delete(ctx context.Context, request *proto.DeleteRequest) (*proto.DeleteReply, error) {

	var bucket store.Bucket

	switch strings.ToLower(request.Type) {
	case "listener":
		bucket = store.Listeners
	case "endpoint":
		bucket = store.Endpoints
	case "profile":
		bucket = store.Profiles
	default:
		return &proto.DeleteReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s not found", request.Type),
		}, nil
	}

	//TODO : error if any running pipes
	err := s.store.Delete(bucket, []byte(request.Uid))

	if err != nil {
		return &proto.DeleteReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.DeleteReply{Ok: true}, nil
}

func (s *handler) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetReply, error) {

	if strings.ToLower(request.Filter) == "all" {
		request.Filter = "e,l,p"
	}

	keys := strings.Split(request.Filter, ",")

	all := []*proto.Entity{}

	for _, key := range keys {

		switch strings.ToLower(key) {
		case "listener", "l":

			listeners := []*proto.Entity{}

			err := s.store.List(store.Listeners, &listeners)

			if err != nil {
				return nil, err
			}
			all = append(all, listeners...)

		case "endpoint", "e":

			endpoints := []*proto.Entity{}

			err := s.store.List(store.Endpoints, &endpoints)

			if err != nil {
				return nil, err
			}
			all = append(all, endpoints...)
		case "profile", "p":

			profiles := []*proto.Entity{}

			err := s.store.List(store.Profiles, &profiles)

			if err != nil {
				return nil, err
			}
			all = append(all, profiles...)

		default:
			return &proto.GetReply{
				Ok:    false,
				Error: fmt.Sprintf("filter of type : %s not registered", key),
			}, nil
		}
	}

	return &proto.GetReply{Ok: true, Entities: all}, nil
}

func (s *handler) Start(ctx context.Context, request *proto.StartRequest) (*proto.StartReply, error) {

	listener := endpoint{}
	endpoints := make([]endpoint, len(request.Endpoints), len(request.Endpoints))

	err := s.store.One(store.Listeners, []byte(request.Listener), &listener)

	if err != nil {
		if err == store.ErrNotFound {
			return &proto.StartReply{
				Ok:    false,
				Error: fmt.Sprintf("listener with uid : %s not found", request.Listener),
			}, nil
		}

		return &proto.StartReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	temp := proto.Entity{}
	err = s.store.One(store.Profiles, []byte(request.Profile), &temp)

	if err != nil {
		if err == store.ErrNotFound {
			return &proto.StartReply{
				Ok:    false,
				Error: fmt.Sprintf("profile with uid : %s not found", request.Profile),
			}, nil
		}

		return &proto.StartReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	profile, _ := profileFromEntity(temp)

	for i, e := range request.Endpoints {
		err = s.store.One(store.Endpoints, []byte(e), &endpoints[i])

		if err != nil {
			if err == store.ErrNotFound {
				return &proto.StartReply{
					Ok:    false,
					Error: fmt.Sprintf("endpoint with uid : %s not found", request.Profile),
				}, nil
			}

			return &proto.StartReply{
				Ok:    false,
				Error: err.Error(),
			}, nil
		}
	}

	pipe := &pipe{
		Uri:       request.Uri,
		Listener:  listener,
		Endpoints: endpoints,
		Profile:   *profile,
		State:     proto.Pipe_UNKNOWN,
	}

	err = s.store.Insert(store.Pipes, []byte(pipe.Uri), pipe)

	if err != nil {

		return &proto.StartReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	err = s.manager.AddPipe(pipe)

	if err != nil {

		deleteError := s.store.Delete(store.Pipes, []byte(pipe.Uri))

		if deleteError != nil {
			// TODO : review this!
			panic(deleteError)
		}

		return &proto.StartReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.StartReply{Ok: true}, nil
}

func profileFromEntity(entity proto.Entity) (*profile, error) {

	return &profile{
		Uid:         entity.Uid,
		Name:        entity.Configuration["profile-name"],
		Description: entity.Configuration["profile-description"],
		Version:     entity.Configuration["profile-version"],
		Script: engine.Script{
			Main:     entity.Configuration["script-main"],
			Runtime:  engine.Runtime(entity.Configuration["script-runtime"]),
			Input:    engine.InputType(entity.Configuration["script-input"]),
			Contents: entity.Configuration["script-contents"],
		},
	}, nil
}

func (s *handler) Stop(ctx context.Context, request *proto.StopRequest) (*proto.StopReply, error) {

	err := s.manager.DeletePipeByURI(request.Uri)

	if err != nil {
		return &proto.StopReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	err = s.store.Delete(store.Pipes, []byte(request.Uri))

	if err != nil {
		return &proto.StopReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.StopReply{Ok: true}, nil
}

func (s *handler) List(ctx context.Context, request *proto.ListRequest) (*proto.ListReply, error) {

	pipes := s.manager.List()

	ppipes := []*proto.Pipe{}

	for _, pipe := range pipes {

		endpoints := []string{}

		for _, e := range pipe.Endpoints {
			endpoints = append(endpoints, e.Uid)
		}

		ppipe := &proto.Pipe{
			Uri:       pipe.Uri,
			Profile:   pipe.Profile.Uid,
			Listener:  pipe.Listener.Uid,
			Endpoints: endpoints,
			Stats: &proto.Statistics{
				Total:  pipe.MessageStatistics.Total,
				Errors: pipe.MessageStatistics.Errors,
				Ok:     pipe.MessageStatistics.OK,
			},
			State: pipe.State,
		}

		ppipes = append(ppipes, ppipe)

	}

	return &proto.ListReply{Ok: true, Pipes: ppipes}, nil
}

func hash(data interface{}) ([]byte, error) {

	j, err := json.Marshal(data)

	if err != nil {
		return []byte{}, err
	}

	checksum := crc32.ChecksumIEEE(j)
	h := hashids.New()

	uid, err := h.Encode([]int{int(checksum)})

	if err != nil {
		return []byte{}, err
	}
	return []byte(uid), nil

}
