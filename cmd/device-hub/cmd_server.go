// Copyright © 2017 thingful

package main

import (
	"context"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"github.com/thingful/device-hub/server"
	"github.com/thingful/device-hub/store"
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

		manager, err := server.NewEndpointManager(ctx, repository)

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
		}, manager)

		return err
	},
}
