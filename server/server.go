// Copyright © 2017 thingful

package server

import (
	context "golang.org/x/net/context"

	"log"
	"net"

	"github.com/thingful/device-hub/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func Serve(manager *manager) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterHubServer(s, &server{
		manager: manager,
	})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	manager *manager
}

func (s *server) PipeList(context.Context, *proto.PipeListRequest) (*proto.PipeListReply, error) {

	pipes := s.manager.List()
	pipes_pb := []*proto.Pipe{}

	for _, p := range pipes {

		pipe_pb := &proto.Pipe{
			Uri: p.Uri,
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
}
