// Copyright © 2017 thingful

package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	context "golang.org/x/net/context"

	"log"
	"net"

	"github.com/thingful/device-hub/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func Serve(options Options, manager *manager) error {

	var grpcServer *grpc.Server

	if options.Insecure {

		grpcServer = grpc.NewServer()

	} else {

		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(options.CertFilePath, options.KeyFilePath)
		if err != nil {
			return fmt.Errorf("could not load server key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(options.TrustedCAFilePath)
		if err != nil {
			return fmt.Errorf("could not read ca certificate: %s", err)
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return errors.New("failed to append client certs")
		}

		// Create the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})

		// Create the gRPC server with the credentials
		grpcServer = grpc.NewServer(grpc.Creds(creds))

	}

	proto.RegisterHubServer(grpcServer, &server{
		manager: manager,
	})

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", options.Binding)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		return err
	}
	return nil
}

type server struct {
	manager *manager
}

type Options struct {
	Binding           string
	Insecure          bool
	CertFilePath      string
	KeyFilePath       string
	TrustedCAFilePath string
}

func (s *server) EndpointAdd(ctx context.Context, request *proto.EndpointAddRequest) (*proto.EndpointAddReply, error) {
	panic("not implemented")
}

func (s *server) EndpointDelete(ctx context.Context, request *proto.EndpointDeleteRequest) (*proto.EndpointDeleteReply, error) {
	panic("not implemented")
}

func (s *server) EndpointList(ctx context.Context, request *proto.EndpointListRequest) (*proto.EndpointListReply, error) {
	panic("not implemented")
}

func (s *server) ListenerAdd(ctx context.Context, request *proto.ListenerAddRequest) (*proto.ListenerAddReply, error) {
	panic("not implemented")
}

func (s *server) ListenerDelete(ctx context.Context, request *proto.ListenerDeleteRequest) (*proto.ListenerDeleteReply, error) {
	panic("not implemented")
}

func (s *server) ListenerList(ctx context.Context, request *proto.ListenerListRequest) (*proto.ListenerListReply, error) {
	panic("not implemented")
}

func (s *server) PipeList(context.Context, *proto.PipeListRequest) (*proto.PipeListReply, error) {

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

func (s *server) PipeDelete(ctx context.Context, request *proto.PipeDeleteRequest) (*proto.PipeDeleteReply, error) {

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

func (s *server) PipeAdd(ctx context.Context, request *proto.PipeAddRequest) (*proto.PipeAddReply, error) {

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

func (s *server) Stats(ctx context.Context, request *proto.StatsRequest) (*proto.StatsReply, error) {
	panic("not implemented")
}
