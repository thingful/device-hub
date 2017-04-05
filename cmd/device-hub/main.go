// Copyright © 2017 thingful

package main

import (
	"flag"
	"fmt"
	"os"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/config"
	"github.com/thingful/go/file"
)

func main() {

	var showVersion bool
	var configurationPath string

	flag.StringVar(&configurationPath, "config", "./config.json", "path to a json configuration file.")
	flag.BoolVar(&showVersion, "version", false, "show version.")

	flag.Parse()

	if showVersion {
		fmt.Println(fmt.Sprintf("device-hub.0.1.%s", hub.SourceVersion))
		return
	}

	// TODO : ensure this path is constrained to a few well known paths
	if !file.Exists(configurationPath) {
		exitWithError(fmt.Errorf("configuration at %s doesn't exist", configurationPath))
	}

	configuration, err := config.LoadProfile(configurationPath)

	if err != nil {
		exitWithError(err)
	}

	app := NewDeviceHub(configuration)

	ctx, err := app.Run()

	if err != nil {
		exitWithError(err)
	}

	<-ctx.Done()
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
