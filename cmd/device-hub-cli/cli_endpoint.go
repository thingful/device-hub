package main

import "github.com/spf13/cobra"

var endpointCommand = &cobra.Command{
	Use:   "endpoint",
	Short: "Add, Delete and List endpoints.",
}
