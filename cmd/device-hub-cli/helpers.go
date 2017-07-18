// Copyright © 2017 thingful

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"

	yaml "gopkg.in/yaml.v2"

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

type rawConf []byte

func (r rawConf) Decode(target interface{}) error {
	err := yaml.Unmarshal(r, target)
	if err != nil {
		return fmt.Errorf("error decoding data: %s", err.Error())
	}
	return nil
}

// Represent a conf file, Data is basically used to order a cliConf Slice
// Raw contains the file content
type cliConf struct {
	Data map[string]interface{}
	Raw  rawConf
}

// Load the configuration file to Data
func (c *cliConf) Load(filePath string) (err error) {

	c.Raw, err = ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file [%s]: %s", filePath, err.Error())
	}
	err = yaml.Unmarshal(c.Raw, &c.Data)
	if err != nil {
		return fmt.Errorf("error parsing file [%s]: %s", filePath, err.Error())
	}
	return nil
}

type cliConfSlice struct {
	C []cliConf
}

func (c *cliConfSlice) Append(e cliConf) {

	c.C = append(c.C, e)
}
func (c cliConfSlice) Len() int {
	return len(c.C)
}
func (c cliConfSlice) Less(i, j int) bool {
	if c.C[j].Data["type"] == "process" {
		return true
	}
	return false
}
func (c cliConfSlice) Swap(i, j int) {
	c.C[i], c.C[j] = c.C[j], c.C[i]
}
func (c cliConfSlice) Print() {
	for k, v := range c.C {
		fmt.Println(k, v.Data)
	}
}

type roundTripFunc func(cli proto.HubClient, in rawConf, out iocodec.Encoder) error

func roundTrip(sample interface{}, caller string, fn roundTripFunc) error {
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

	var dataSlice cliConfSlice

	// either no request file is not specified or set to std-in
	if (cfg.RequestFile == "" && cfg.RequestDir == "") || cfg.RequestFile == "-" {
		// Add empty item to iterate - TODO refactor
		dataSlice.Append(cliConf{})
		// or request file is specified
	} else if cfg.RequestFile != "" {

		var data cliConf

		err := data.Load(cfg.RequestFile)
		if err != nil {
			return err
		}

		dataSlice.Append(data)

		// or request dir is specified
	} else if cfg.RequestDir != "" {
		listing, err := ioutil.ReadDir(cfg.RequestDir)

		if err != nil {
			return err
		}

		for _, fi := range listing {

			fmt.Println(fi.Name())

			folderPath := path.Join(cfg.RequestDir, fi.Name())

			var data cliConf
			err = data.Load(folderPath)
			if err != nil {
				return err
			}
			dataSlice.Append(data)
		}
	}

	conn, client, err := dial()
	if err != nil {
		return err
	}
	defer conn.Close()
	// sort the config items to create & delete cases
	switch caller {
	case "create":
		sort.Sort(cliConfSlice(dataSlice))
	case "delete":
		sort.Sort(sort.Reverse(cliConfSlice(dataSlice)))
	}

	for _, d := range dataSlice.C {
		err := fn(client, d.Raw, em.NewEncoder(os.Stdout))

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
		if ext != "yaml" {
			return nil, nil, fmt.Errorf("invalid request file format: %q", ext)

		}
	}

	dm, _ := iocodec.DefaultDecoders["yaml"]

	return f, dm.NewDecoder(f), nil
}

func yamlDecoder(filePath string, target interface{}) error {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file [%s]: %s", filePath, err.Error())
	}
	err = yaml.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("error parsing file [%s]: %s", filePath, err.Error())
	}
	return nil
}
