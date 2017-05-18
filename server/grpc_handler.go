// Copyright © 2017 thingful

package server

import (
	"fmt"
	"strings"

	"github.com/thingful/device-hub/proto"
	context "golang.org/x/net/context"
)

type handler struct {
	manager *manager
}

// Create inserts or overwrites listener, endpoint and profile entities
func (s *handler) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateReply, error) {

	entity := proto.Entity{
		Type:          request.Type,
		Kind:          request.Kind,
		Configuration: request.Configuration,
	}

	uid, err := s.manager.Repository.UpdateOrCreateEntity(entity)

	if err != nil {
		return &proto.CreateReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.CreateReply{
		Uid: uid,
		Ok:  true,
	}, nil
}

// Delete removes listener, endpoint and profile entities
func (s *handler) Delete(ctx context.Context, request *proto.DeleteRequest) (*proto.DeleteReply, error) {

	// first check if the item is being used in
	// any running pipes by consulting the manager
	// TODO : consider moving this into the manager
	running := true

	switch strings.ToLower(request.Type) {
	case "listener":

		running = s.manager.Any(func(p *pipe) bool {
			return p.Listener.Uid == request.Uid
		})

	case "endpoint":
		running = s.manager.Any(func(p *pipe) bool {
			for _, e := range p.Endpoints {
				if e.Uid == request.Uid {
					return true
				}
			}
			return false
		})
	case "profile":

		running = s.manager.Any(func(p *pipe) bool {
			return p.Profile.Uid == request.Uid
		})

	default:
		return &proto.DeleteReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s not found", request.Type),
		}, nil
	}

	if running {
		return &proto.DeleteReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s, uid : %s is currently being used, stop pipe before deleting", request.Type, request.Uid),
		}, nil

	}

	err := s.manager.Repository.Delete(request.Type, request.Uid)

	if err != nil {
		return &proto.DeleteReply{
			Ok:    false,
			Error: err.Error(),
		}, nil

	}

	return &proto.DeleteReply{Ok: true}, nil
}

// Get is a generic method to list listener, endpoint and profile entities
func (s *handler) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetReply, error) {

	if strings.ToLower(request.Filter) == "all" {
		request.Filter = "e,l,p"
	}

	all, err := s.manager.Repository.Search(request.Filter)

	if err != nil {
		return &proto.GetReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.GetReply{Ok: true, Entities: all}, nil
}

// Start will start a 'pipe'
func (s *handler) Start(ctx context.Context, request *proto.StartRequest) (*proto.StartReply, error) {

	err := s.manager.StartPipe(request.Uri,
		request.Listener,
		request.Profile,
		request.Endpoints,
		request.Tags)

	if err != nil {
		return &proto.StartReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.StartReply{Ok: true}, nil
}

// Stop will stop a 'pipe'
func (s *handler) Stop(ctx context.Context, request *proto.StopRequest) (*proto.StopReply, error) {

	if request.Uri == "" {
		return &proto.StopReply{
			Ok:    false,
			Error: "pipe uri not supplied",
		}, nil
	}

	err := s.manager.DeletePipe(func(p *pipe) bool {
		return p.Uri == request.Uri
	})

	if err != nil {
		return &proto.StopReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.StopReply{Ok: true}, nil
}

// List returns all running 'pipes'
func (s *handler) List(ctx context.Context, request *proto.ListRequest) (*proto.ListReply, error) {

	pipes := s.manager.List()

	ppipes := []*proto.Pipe{}

	for _, pipe := range pipes {

		endpoints := []string{}

		for _, e := range pipe.Endpoints {
			endpoints = append(endpoints, e.Uid)
		}

		ffs := pipe.MessageStatistics

		ppipe := &proto.Pipe{
			Uri:       pipe.Uri,
			Profile:   pipe.Profile.Uid,
			Listener:  pipe.Listener.Uid,
			Endpoints: endpoints,
			Stats:     &ffs,
			State:     pipe.State,
		}

		ppipes = append(ppipes, ppipe)

	}

	return &proto.ListReply{Ok: true, Pipes: ppipes}, nil
}
