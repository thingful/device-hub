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
	var checkConfig bool
	var configurationPath string

	flag.StringVar(&configurationPath, "config", "./config.json", "path to a json configuration file.")
	flag.BoolVar(&showVersion, "version", false, "show version.")
	flag.BoolVar(&checkConfig, "check", false, "validate configuration file.")

	flag.Parse()

	if showVersion {
		fmt.Println(fmt.Sprintf("device-hub.0.1.%s", hub.SourceVersion))
		return
	}

	if !file.Exists(configurationPath) {
		exitWithError(fmt.Errorf("configuration at %s doesn't exist", configurationPath))
	}

	configuration, err := config.LoadFromFile(configurationPath)

	if err != nil {
		exitWithError(err)
	}

	err = config.Validate(configuration)

	if checkConfig {
		if err != nil {
			fmt.Println("ERROR : ")
			fmt.Println(err.Error())
			return
		}
		fmt.Println("PASSED")
		return
	}

	// don't start with ANY error in the initial configuration file
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
