// Copyright © 2017 thingful

package server

import (
	"fmt"
	"log"

	context "golang.org/x/net/context"

	"github.com/thingful/device-hub/proto"
)

type handler struct {
	manager *manager
	store   *store
}

func (s *handler) EndpointAdd(ctx context.Context, request *proto.EndpointAddRequest) (*proto.EndpointAddReply, error) {

	uid, err := s.store.Insert(endPoints, request.Endpoint)
	//TODO : handle update
	// if update then look at policy to reload existing

	if err != nil {
		return &proto.EndpointAddReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	//	s.store.List(endPoints, nil)
	pp := proto.Endpoint{}
	s.store.Get(endPoints, uid, &pp)

	fmt.Println(pp, pp.Type)

	return &proto.EndpointAddReply{
		Ok:  true,
		Uid: uid,
	}, nil
}

func (s *handler) EndpointDelete(ctx context.Context, request *proto.EndpointDeleteRequest) (*proto.EndpointDeleteReply, error) {
	panic("not implemented")
}

func (s *handler) EndpointList(ctx context.Context, request *proto.EndpointListRequest) (*proto.EndpointListReply, error) {
	panic("not implemented")
}

func (s *handler) ListenerAdd(ctx context.Context, request *proto.ListenerAddRequest) (*proto.ListenerAddReply, error) {
	panic("not implemented")
}

func (s *handler) ListenerDelete(ctx context.Context, request *proto.ListenerDeleteRequest) (*proto.ListenerDeleteReply, error) {
	panic("not implemented")
}

func (s *handler) ListenerList(ctx context.Context, request *proto.ListenerListRequest) (*proto.ListenerListReply, error) {
	panic("not implemented")
}

func (s *handler) PipeList(context.Context, *proto.PipeListRequest) (*proto.PipeListReply, error) {

	/*
		pipes := s.manager.List()
		pipes_pb := []*proto.Pipe{}

		for _, p := range pipes {

			pipe_pb := &proto.Pipe{
				Uri:   p.Uri,
				State: proto.PipeState(proto.PipeState_value[string(p.State)]),
				Profile: &proto.Profile{
					Name:        p.Profile.Name,
					Description: p.Profile.Description,
					Version:     p.Profile.Version,
				},
				Listener: &proto.Endpoint{
					Uid:  string(p.Listener.UID),
					Type: p.Listener.Type,
				},
				Endpoints: []*proto.Endpoint{},
				MessageStats: &proto.Statistics{
					Total:  p.MessageStatistics.Total,
					Errors: p.MessageStatistics.Errors,
					Ok:     p.MessageStatistics.OK,
				},
			}

			for _, e := range p.Endpoints {

				endpoint_pb := &proto.Endpoint{
					Uid:  string(e.UID),
					Type: e.Type,
				}

				pipe_pb.Endpoints = append(pipe_pb.Endpoints, endpoint_pb)
			}

			pipes_pb = append(pipes_pb, pipe_pb)
		}

		return &proto.PipeListReply{Pipes: pipes_pb}, nil
	*/
	return &proto.PipeListReply{}, nil
}

func (s *handler) PipeDelete(ctx context.Context, request *proto.PipeDeleteRequest) (*proto.PipeDeleteReply, error) {

	err := s.manager.DeletePipeByUID(request.Uri)

	if err != nil {
		log.Print("PipeDelete", err.Error())
		return &proto.PipeDeleteReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.PipeDeleteReply{
		Ok: true,
	}, nil
}

func (s *handler) PipeAdd(ctx context.Context, request *proto.PipeAddRequest) (*proto.PipeAddReply, error) {

	err := s.manager.AddPipe(request.Uri, request.ProfileUid, request.ListenerUid, request.EndpointUids)

	if err != nil {
		log.Print("PipeAdd", err.Error())
		return &proto.PipeAddReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.PipeAddReply{
		Ok: true,
	}, nil

}

func (s *handler) Stats(ctx context.Context, request *proto.StatsRequest) (*proto.StatsReply, error) {
	panic("not implemented")
}
