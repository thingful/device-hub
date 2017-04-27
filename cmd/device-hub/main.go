// Copyright © 2017 thingful

package main

import (
	"os"

	"google.golang.org/grpc"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	_ "github.com/thingful/device-hub/endpoint"
	_ "github.com/thingful/device-hub/listener"
)

var RootCmd = &cobra.Command{
	Use: "device-hub",
}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

var _config = newConfig()

type config struct {
	Binding    string `envconfig:"BINDING" default:":50051"`
	TLS        bool   `envconfig:"TLS"`
	ServerName string `envconfig:"TLS_SERVER_NAME"`
	CACertFile string `envconfig:"TLS_CA_CERT_FILE"`
	CertFile   string `envconfig:"TLS_CERT_FILE"`
	KeyFile    string `envconfig:"TLS_KEY_FILE"`
	Data       string `envconfig:"DATA" default:"."`
}

func newConfig() *config {
	c := &config{}
	envconfig.Process("", c)
	return c
}

func (o *config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Binding, "binding", "b", o.Binding, "binding address in form of {ip}:port")
	fs.BoolVar(&o.TLS, "tls", o.TLS, "enable tls")
	fs.StringVar(&o.CACertFile, "tls-ca-cert-file", o.CACertFile, "ca certificate file")
	fs.StringVar(&o.CertFile, "tls-cert-file", o.CertFile, "client certificate file")
	fs.StringVar(&o.KeyFile, "tls-key-file", o.KeyFile, "client key file")
	fs.StringVar(&o.Data, "data", o.Data, "path to db folder, defaults to current directory")
}

func init() {

	RootCmd.AddCommand(versionCommand)
	RootCmd.AddCommand(serverCommand)
	_config.AddFlags(RootCmd.PersistentFlags())
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
