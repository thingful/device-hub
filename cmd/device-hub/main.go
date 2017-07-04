// Copyright © 2017 thingful

package main

import (
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"

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
	Binding    string `envconfig:"BINDING" default:":50051" yaml:"binding,omitempty"`
	TLS        bool   `envconfig:"TLS" yaml:"tls,omitempty"`
	ServerName string `envconfig:"TLS_SERVER_NAME" yaml:"server_name,omitempty"`
	CACertFile string `envconfig:"TLS_CA_CERT_FILE" yaml:"ca_cert_file,omitempty"`
	CertFile   string `envconfig:"TLS_CERT_FILE" yaml:"cert_file,omitempty"`
	KeyFile    string `envconfig:"TLS_KEY_FILE" yaml:"key_file,omitempty"`
	DataDir    string `envconfig:"DATA_DIR" default:"." yaml:"data_dir,omitempty"`
	DataImpl   string `envconfig:"DATA_IMPL" default:"boltdb" yaml:"data_impl,omitempty"`
	LogFile    bool   `envconfig:"LOG_FILE" yaml:"log_file,omitempty"`
	LogPath    string `envconfig:"LOG_PATH" default:"./device-hub.log" yaml:"log_path,omitempty"`
	Syslog     bool   `envconfig:"LOG_SYSLOG" yaml:"sys_log,omitempty"`
	ConfigFile bool   `envconfig:"CONFIG_FILE"`
	ConfigPath string `envconfig:"CONFIG_PATH" default:"./config.yaml"`
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
	fs.StringVar(&o.DataDir, "data-dir", o.DataDir, "path to db folder, defaults to current directory")
	fs.StringVar(&o.DataImpl, "data-impl", o.DataImpl, "datastore to use, valid values are 'boltdb' or 'filestore', defaults to boltdb")
	fs.BoolVar(&o.LogFile, "log-file", o.LogFile, "enable log to file")
	fs.StringVar(&o.LogPath, "log-path", o.LogPath, "path to log file, defaults to ./device-hub.log")
	fs.BoolVar(&o.Syslog, "log-syslog", o.Syslog, "enable log to local SYSLOG")
	fs.BoolVar(&o.ConfigFile, "config-file", o.ConfigFile, "enable config file overriding flags and env vars")
	fs.StringVar(&o.ConfigPath, "config-path", o.ConfigPath, "path to config file, defaults to ./config.yaml")
}

func (o *config) AddConfigFile(cf string) {
	content, err := ioutil.ReadFile(cf)
	if err != nil {
		log.Fatalf("Failed to read config file config.yaml: %s\n", err.Error())
	}
	// TODO check & test nontags fields!!!!!
	err = yaml.Unmarshal(content, &o)
	if err != nil {
		log.Fatalf("Error parsing config file: %s\n", err.Error())
	}
}

func init() {
	RootCmd.AddCommand(versionCommand)
	RootCmd.AddCommand(serverCommand)
	_config.AddFlags(RootCmd.PersistentFlags())
	if _config.ConfigFile {
		log.Println("Overriding settings with config file")
		_config.AddConfigFile(_config.ConfigPath)
	}
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
