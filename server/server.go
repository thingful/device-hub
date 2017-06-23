// Copyright © 2017 thingful

package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"net"

	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Options struct {
	Binding           string
	UseTLS            bool
	CertFilePath      string
	KeyFilePath       string
	TrustedCAFilePath string
	LogFile           string
	Syslog            bool
}

func Serve(options Options, manager *runtime.Manager) error {

	var grpcServer *grpc.Server

	if options.UseTLS {

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

	} else {

		grpcServer = grpc.NewServer()

	}

	proto.RegisterHubServer(grpcServer, &handler{
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
