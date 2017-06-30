// Copyright © 2017 thingful

package server

import (
	"fmt"
	"strings"

	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/runtime"
	context "golang.org/x/net/context"
)

type handler struct {
	manager *runtime.Manager
}

// Create inserts listener, endpoint and profile entities
func (s *handler) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateReply, error) {

	entity := proto.Entity{
		Uid:           request.Uid,
		Type:          request.Type,
		Kind:          request.Kind,
		Configuration: request.Configuration,
	}

	uid, err := s.manager.Repository.Insert(entity)

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

	entity := proto.Entity{
		Uid:           request.Uid,
		Type:          request.Type,
		Configuration: request.Configuration,
	}

	switch strings.ToLower(entity.Type) {
	case "listener":

		running = s.manager.Any(func(p *runtime.Pipe) bool {
			return p.Listener.Uid == entity.Uid
		})

	case "endpoint":
		running = s.manager.Any(func(p *runtime.Pipe) bool {
			for _, e := range p.Endpoints {
				if e.Uid == entity.Uid {
					return true
				}
			}
			return false
		})
	case "profile":

		running = s.manager.Any(func(p *runtime.Pipe) bool {
			return p.Profile.Uid == entity.Uid
		})

	default:
		return &proto.DeleteReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s not found", entity.Type),
		}, nil
	}

	if running {
		return &proto.DeleteReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s, uid : %s is currently being used, stop pipe before deleting", entity.Type, entity.Uid),
		}, nil

	}

	err := s.manager.Repository.Delete(entity)

	if err != nil {
		return &proto.DeleteReply{
			Ok:    false,
			Error: err.Error(),
		}, nil

	}

	return &proto.DeleteReply{Ok: true}, nil
}

// Show is a generic method to list listener, endpoint and profile entities
func (s *handler) Show(ctx context.Context, request *proto.ShowRequest) (*proto.ShowReply, error) {

	if strings.ToLower(request.Filter) == "all" {
		request.Filter = "e,l,p"
	}

	all, err := s.manager.Repository.Search(request.Filter)

	if err != nil {
		return &proto.ShowReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.ShowReply{Ok: true, Entities: all}, nil
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

	err := s.manager.DeletePipe(func(p *runtime.Pipe) bool {
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

// Status returns all running 'pipes'
func (s *handler) Status(ctx context.Context, request *proto.StatusRequest) (*proto.StatusReply, error) {

	pipes := s.manager.Status()

	ppipes := []*proto.Pipe{}

	for _, pipe := range pipes {

		endpoints := []string{}

		for _, e := range pipe.Endpoints {
			endpoints = append(endpoints, e.Uid)
		}

		ffs := pipe.Statistics

		ppipe := &proto.Pipe{
			Uri:       pipe.Uri,
			Profile:   pipe.Profile.Uid,
			Listener:  pipe.Listener.Uid,
			Endpoints: endpoints,
			Stats:     ffs,
			State:     pipe.State,
		}

		ppipes = append(ppipes, ppipe)

	}

	return &proto.StatusReply{Ok: true, Pipes: ppipes}, nil
}
