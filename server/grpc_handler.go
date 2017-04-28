// Copyright © 2017 thingful

package server

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strings"

	hashids "github.com/speps/go-hashids"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/config"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/proto"
	context "golang.org/x/net/context"
)

type handler struct {
	manager *manager
	store   *store
}

func (s *handler) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateReply, error) {

	hash, err := hash(request)

	if err != nil {
		return &proto.CreateReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	var bucket bucket

	switch strings.ToLower(request.Type) {
	case "listener":
		bucket = listenersBucket

		exists := hub.IsListenerRegistered(request.Kind)

		if !exists {
			return &proto.CreateReply{
				Ok:    false,
				Error: fmt.Sprintf("kind : %s not registered", request.Kind),
			}, nil

		}

	case "endpoint":
		bucket = endpointsBucket

		exists := hub.IsEndpointRegistered(request.Kind)

		if !exists {
			return &proto.CreateReply{
				Ok:    false,
				Error: fmt.Sprintf("kind : %s not registered", request.Kind),
			}, nil
		}
	case "profile":
		bucket = profilesBucket

		if request.Configuration["profile-name"] != "" {
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

	var bucket bucket

	switch strings.ToLower(request.Type) {
	case "listener":
		bucket = listenersBucket
	case "endpoint":
		bucket = endpointsBucket
	case "profile":
		bucket = profilesBucket
	default:
		return &proto.DeleteReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s not found", request.Type),
		}, nil
	}

	//TODO : error if any running pipes
	err := s.store.Delete(bucket, request.Uid)

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

			err := s.store.List(listenersBucket, &listeners)

			if err != nil {
				return nil, err
			}
			all = append(all, listeners...)

		case "endpoint", "e":

			endpoints := []*proto.Entity{}

			err := s.store.List(endpointsBucket, &endpoints)

			if err != nil {
				return nil, err
			}
			all = append(all, endpoints...)
		case "profile", "p":

			profiles := []*proto.Entity{}

			err := s.store.List(profilesBucket, &profiles)

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

	listener := config.Endpoint{}
	endpoints := make([]config.Endpoint, len(request.Endpoints), len(request.Endpoints))

	err := s.store.One(listenersBucket, request.Listener, &listener)

	if err != nil {
		if err == ErrNotFound {
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
	err = s.store.One(profilesBucket, request.Profile, &temp)

	if err != nil {
		if err == ErrNotFound {
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
		err = s.store.One(endpointsBucket, e, &endpoints[i])

		if err != nil {
			if err == ErrNotFound {
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
		State:     UNKNOWN,
	}

	// TODO : replace with method
	s.manager.pipes[request.Uri] = pipe
	err = s.manager.Start()

	if err != nil {
		return &proto.StartReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.StartReply{Ok: true}, nil
}

func profileFromEntity(entity proto.Entity) (*config.Profile, error) {

	return &config.Profile{
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
	return &proto.StopReply{Ok: true}, nil
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
