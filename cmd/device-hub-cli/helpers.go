package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

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

type roundTripFunc func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error

func roundTrip(sample interface{}, fn roundTripFunc) error {
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

	if cfg.PrintSampleRequest {
		return em.NewEncoder(os.Stdout).Encode(sample)
	}

	decoders := []iocodec.Decoder{}
	files := []*os.File{}

	// either no request file is not specified or set to std-in
	if (cfg.RequestFile == "" && cfg.RequestDir == "") || cfg.RequestFile == "-" {
		decoders = append(decoders, iocodec.DefaultDecoders["json"].NewDecoder(os.Stdin))

		// or request file is specified
	} else if cfg.RequestFile != "" {

		f, d, err := decoderFromPath(cfg.RequestFile)

		if err != nil {
			return err
		}

		decoders = append(decoders, d)
		files = append(files, f)

		// or request dir is specified
	} else if cfg.RequestDir != "" {
		listing, err := ioutil.ReadDir(cfg.RequestDir)

		if err != nil {
			return err
		}

		for _, fi := range listing {

			fmt.Println(fi.Name())

			f, d, err := decoderFromPath(cfg.RequestDir + fi.Name())

			if err != nil {
				return err
			}

			decoders = append(decoders, d)
			files = append(files, f)

		}
	}

	defer func() {
		for i, _ := range files {
			files[i].Close()
		}
	}()

	conn, client, err := dial()
	if err != nil {
		return err
	}

	defer conn.Close()

	for _, d := range decoders {
		err := fn(client, d, em.NewEncoder(os.Stdout))

		if err != nil {
			return err
		}
	}

	return nil
}

func decoderFromPath(filePath string) (*os.File, iocodec.Decoder, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("request file: %v", err)
	}

	ext := filepath.Ext(filePath)

	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	dm, ok := iocodec.DefaultDecoders[ext]

	if !ok {
		return nil, nil, fmt.Errorf("invalid request file format: %q", ext)
	}

	return f, dm.NewDecoder(f), nil
}
