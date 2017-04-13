package main

import "github.com/spf13/cobra"

var listenerCommand = &cobra.Command{
	Use:   "listener",
	Short: "Add, Delete and List listeners.",
}
