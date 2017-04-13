package main

import "github.com/spf13/cobra"

var statsCommand = &cobra.Command{
	Use:   "stats",
	Short: "View statistics.",
}
