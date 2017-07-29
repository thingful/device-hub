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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load Config File if is enabled
		if _config.ConfigFile {
			_config.AddConfigFile(_config.ConfigPath)

		}
	},
}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

var _config = newConfig()

type config struct {
	Binding    string  `envconfig:"BINDING" default:":50051" yaml:"binding"`
	TLS        bool    `envconfig:"TLS" yaml:"tls"`
	ServerName string  `envconfig:"TLS_SERVER_NAME" yaml:"tls_server_name"`
	CACertFile string  `envconfig:"TLS_CA_CERT_FILE" yaml:"tls_ca_cert_file"`
	CertFile   string  `envconfig:"TLS_CERT_FILE" yaml:"tls_cert_file"`
	KeyFile    string  `envconfig:"TLS_KEY_FILE" yaml:"tls_key_file"`
	DataDir    string  `envconfig:"DATA_DIR" default:"." yaml:"data_dir"`
	DataImpl   string  `envconfig:"DATA_IMPL" default:"boltdb" yaml:"data_impl"`
	LogFile    bool    `envconfig:"LOG_FILE" yaml:"log_file"`
	LogPath    string  `envconfig:"LOG_PATH" default:"./device-hub.log" yaml:"log_path"`
	Syslog     bool    `envconfig:"LOG_SYSLOG" yaml:"sys_log"`
	ConfigFile bool    `envconfig:"CONFIG_FILE" yaml:"-"`
	ConfigPath string  `envconfig:"CONFIG_PATH" default:"./config.yaml" yaml:"-"`
	GeoEnabled bool    `envconfig:"GEO_ENABLED" yaml:"geo_enabled"`
	GeoLat     float64 `envconfig:"GEO_LAT" yaml:"geo_lat"`
	GeoLng     float64 `envconfig:"GEO_LNG" yaml:"geo_lng"`
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
	fs.BoolVarP(&o.LogFile, "log-file", "l", o.LogFile, "enable log to file")
	fs.StringVar(&o.LogPath, "log-path", o.LogPath, "path to log file, defaults to ./device-hub.log")
	fs.BoolVar(&o.Syslog, "log-syslog", o.Syslog, "enable log to local SYSLOG")
	fs.BoolVarP(&o.ConfigFile, "config-file", "c", o.ConfigFile, "enable config file overriding flags and env vars")
	fs.StringVar(&o.ConfigPath, "config-path", o.ConfigPath, "path to config file, defaults to ./config.yaml")
	fs.BoolVar(&o.GeoEnabled, "geo-enabled", o.GeoEnabled, "enable geo location")
	fs.Float64Var(&o.GeoLat, "geo-lat", o.GeoLat, "device-hub geo latitude")
	fs.Float64Var(&o.GeoLng, "geo-lng", o.GeoLng, "device-hub geo longitude")
	fs.Parse(os.Args)
}

func (o *config) AddConfigFile(cf string) {

	var readFields map[string]interface{}

	content, err := ioutil.ReadFile(cf)
	if err != nil {
		log.Fatalf("Failed to read config file config.yaml: %s\n", err.Error())
	}

	err = yaml.Unmarshal(content, &o)
	if err != nil {
		log.Fatalf("Error parsing config file: %s\n", err.Error())
	}

	err = yaml.Unmarshal(content, &readFields)
	if err != nil {
		log.Fatalf("Error parsing config file: %s\n", err.Error())
	}

	keys := make([]string, 0)

	for k, _ := range readFields {
		keys = append(keys, k)
	}
	log.Printf("Overriding keys: %v with config file\n", keys)
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
