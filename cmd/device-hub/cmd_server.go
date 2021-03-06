// Copyright © 2017 thingful

package main

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/endpoint"
	"github.com/thingful/device-hub/listener"
	"github.com/thingful/device-hub/registry"
	"github.com/thingful/device-hub/runtime"
	"github.com/thingful/device-hub/server"
	"github.com/thingful/device-hub/store"
	"github.com/thingful/device-hub/utils"
)

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Start device hub.",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()

		var dataImpl store.Storer

		switch strings.ToLower(_config.DataImpl) {
		case "boltdb":

			dbFile := path.Join(_config.DataDir, "device-hub.db")
			db, err := bolt.Open(dbFile, 0600, nil)

			if err != nil {
				return err
			}

			dataImpl = store.NewBoltDBStore(db)

		case "filestore":
			dataImpl = store.NewFileStore(_config.DataDir)

		default:
			return fmt.Errorf("unknown data implementation : %s, valid values are 'boltdb' or 'filestore'", _config.DataImpl)
		}

		defer dataImpl.Close()

		register := registry.Default

		endpoint.Register(register)
		listener.Register(register)

		repository := store.NewRepository(dataImpl, register)

		options := runtime.Options{}

		if _config.GeoEnabled {
			options.GeoEnabled = true
			options.GeoLat = _config.GeoLat
			options.GeoLng = _config.GeoLng
		} else {
			options.GeoEnabled = false
		}

		manager, err := runtime.NewEndpointManager(ctx,
			repository,
			register,
			utils.NewLogger(hub.DaemonVersionString(),
				_config.Syslog,
				_config.LogFile,
				_config.LogPath),
			options,
		)

		if err != nil {
			return err
		}

		// starting the manager will ensure that the previous
		// running state is recreated
		err = manager.Start()

		if err != nil {
			return err
		}

		err = server.Serve(server.Options{
			Binding:           _config.Binding,
			UseTLS:            _config.TLS,
			CertFilePath:      _config.CertFile,
			KeyFilePath:       _config.KeyFile,
			TrustedCAFilePath: _config.CACertFile,
			LogFile:           _config.LogFile,
			LogPath:           _config.LogPath,
			Syslog:            _config.Syslog,
		}, manager)

		return err
	},
}
