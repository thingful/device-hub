// Copyright © 2017 thingful

package main

import (
	"os"
	"time"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	_ "github.com/thingful/device-hub/endpoint"
	_ "github.com/thingful/device-hub/listener"
)

var RootCmd = &cobra.Command{
	Use: "device-hub-cli",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := _resources.SetResources(_config)
		if err != nil {
			return err
		}
		return nil
	},
}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

var (
	_config    = newConfig()
	_resources = newResources()
	_encoder   iocodec.Encoder
)

type config struct {
	ServerAddr         string        `envconfig:"SERVER_ADDR" default:"127.0.0.1:50051"`
	RequestFile        string        `envconfig:"REQUEST_FILE"`
	RequestDir         string        `envconfig:"REQUEST_DIR"`
	PrintSampleRequest bool          `envconfig:"PRINT_SAMPLE_REQUEST"`
	ResponseFormat     string        `envconfig:"RESPONSE_FORMAT" default:"json"`
	Timeout            time.Duration `envconfig:"TIMEOUT" default:"10s"`
	TLS                bool          `envconfig:"TLS"`
	ServerName         string        `envconfig:"TLS_SERVER_NAME"`
	InsecureSkipVerify bool          `envconfig:"TLS_INSECURE_SKIP_VERIFY"`
	CACertFile         string        `envconfig:"TLS_CA_CERT_FILE"`
	CertFile           string        `envconfig:"TLS_CERT_FILE"`
	KeyFile            string        `envconfig:"TLS_KEY_FILE"`
	AuthToken          string        `envconfig:"AUTH_TOKEN"`
	AuthTokenType      string        `envconfig:"AUTH_TOKEN_TYPE" default:"Bearer"`
	JWTKey             string        `envconfig:"JWT_KEY"`
	JWTKeyFile         string        `envconfig:"JWT_KEY_FILE"`
	ProcessFile        processConf
}

func newConfig() *config {
	c := &config{}
	envconfig.Process("", c)
	return c
}

func (o *config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.ServerAddr, "server-addr", "s", o.ServerAddr, "server address in form of host:port")
	fs.StringVarP(&o.RequestFile, "request-file", "f", o.RequestFile, "client request file (must be yaml); use \"-\" for stdin + json")
	fs.StringVarP(&o.RequestDir, "request-dir", "d", o.RequestDir, "directory containing client request file(s) (must be yaml)")
	fs.BoolVarP(&o.PrintSampleRequest, "print-sample-request", "p", o.PrintSampleRequest, "print sample request file and exit")
	fs.StringVarP(&o.ResponseFormat, "response-format", "o", o.ResponseFormat, "response format (json, prettyjson or yaml)")
	fs.DurationVar(&o.Timeout, "timeout", o.Timeout, "client connection timeout")
	fs.BoolVar(&o.TLS, "tls", o.TLS, "enable tls")
	fs.StringVar(&o.ServerName, "tls-server-name", o.ServerName, "tls server name override")
	fs.BoolVar(&o.InsecureSkipVerify, "tls-insecure-skip-verify", o.InsecureSkipVerify, "INSECURE: skip tls checks")
	fs.StringVar(&o.CACertFile, "tls-ca-cert-file", o.CACertFile, "ca certificate file")
	fs.StringVar(&o.CertFile, "tls-cert-file", o.CertFile, "client certificate file")
	fs.StringVar(&o.KeyFile, "tls-key-file", o.KeyFile, "client key file")
	fs.StringVar(&o.AuthToken, "auth-token", o.AuthToken, "authorization token")
	fs.StringVar(&o.AuthTokenType, "auth-token-type", o.AuthTokenType, "authorization token type")
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "jwt key")
	fs.StringVar(&o.JWTKeyFile, "jwt-key-file", o.JWTKeyFile, "jwt key file")
}

func init() {
	RootCmd.AddCommand(versionCommand)

	RootCmd.AddCommand(createCommand)
	RootCmd.AddCommand(showCommand)
	RootCmd.AddCommand(deleteCommand)
	RootCmd.AddCommand(startCommand())
	RootCmd.AddCommand(stopCommand)
	RootCmd.AddCommand(statusCommand)
	RootCmd.AddCommand(describeCommand)

	_config.AddFlags(RootCmd.PersistentFlags())

	var em iocodec.EncoderMaker
	em = iocodec.DefaultEncoders["json"]
	_encoder = em.NewEncoder(os.Stdout)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
