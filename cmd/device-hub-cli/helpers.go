// Copyright © 2017 thingful

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/thingful/device-hub/proto"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func dial() (*grpc.ClientConn, proto.HubClient, error) {
	cfg := _config

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(cfg.Timeout),
	}

	if cfg.TLS {

		tlsConfig := &tls.Config{}

		if cfg.InsecureSkipVerify {
			tlsConfig.InsecureSkipVerify = true
		}

		if cfg.CACertFile != "" {
			cacert, err := ioutil.ReadFile(cfg.CACertFile)
			if err != nil {
				return nil, nil, fmt.Errorf("ca cert: %v", err)
			}
			certpool := x509.NewCertPool()
			certpool.AppendCertsFromPEM(cacert)
			tlsConfig.RootCAs = certpool
		}

		if cfg.CertFile != "" {
			if cfg.KeyFile == "" {
				return nil, nil, fmt.Errorf("missing key file")
			}
			pair, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
			if err != nil {
				return nil, nil, fmt.Errorf("cert/key: %v", err)
			}
			tlsConfig.Certificates = []tls.Certificate{pair}
		}

		if cfg.ServerName != "" {
			tlsConfig.ServerName = cfg.ServerName

		} else {

			addr, _, _ := net.SplitHostPort(cfg.ServerAddr)
			tlsConfig.ServerName = addr
		}

		cred := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(cred))

	} else {

		opts = append(opts, grpc.WithInsecure())
	}

	if cfg.AuthToken != "" {

		cred := oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: cfg.AuthToken,
			TokenType:   cfg.AuthTokenType,
		})
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	}
	if cfg.JWTKey != "" {

		cred, err := oauth.NewJWTAccessFromKey([]byte(cfg.JWTKey))
		if err != nil {
			return nil, nil, fmt.Errorf("jwt key: %v", err)
		}
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	}
	if cfg.JWTKeyFile != "" {

		cred, err := oauth.NewJWTAccessFromFile(cfg.JWTKeyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("jwt key file: %v", err)
		}
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	}

	conn, err := grpc.Dial(cfg.ServerAddr, opts...)

	if err != nil {
		return nil, nil, err
	}

	return conn, proto.NewHubClient(conn), nil
}

type roundTripFunc func(cli proto.HubClient, in rawContent, out iocodec.Encoder) error

func roundTrip(fn roundTripFunc) error {
	cfg := _config

	var em iocodec.EncoderMaker
	var ok bool

	if cfg.ResponseFormat == "" {
		em = iocodec.DefaultEncoders["json"]
	} else {

		em, ok = iocodec.DefaultEncoders[cfg.ResponseFormat]
		if !ok {
			return fmt.Errorf("invalid response format: %q", cfg.ResponseFormat)
		}
	}

	conn, client, err := dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = fn(client, nil, em.NewEncoder(os.Stdout))
	if err != nil {
		return err
	}

	return nil
}
