// Copyright © 2017 thingful

package main

import (
	"context"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/runtime"
	"github.com/thingful/device-hub/server"
	"github.com/thingful/device-hub/store"
	"github.com/thingful/device-hub/utils"
)

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Start device hub.",
	RunE: func(cmd *cobra.Command, args []string) error {

		dbFile := fmt.Sprintf("%s/device-hub.db", _config.Data)

		db, err := bolt.Open(dbFile, 0600, nil)

		if err != nil {
			return err
		}
		defer db.Close()

		ctx := context.Background()

		s := store.NewStore(db)
		repository := store.NewRepository(s)

		manager, err := runtime.NewEndpointManager(ctx, repository, utils.NewLogger(hub.DaemonVersionString(), _config.LogFile))

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
		}, manager)

		return err
	},
}
