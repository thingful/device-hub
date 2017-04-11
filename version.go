package hub

import "fmt"

var (
	// SourceVersion is set via the makefile
	SourceVersion = "DEVELOPMENT"
)

// DaemonVersionString returns the long version of the daemon executable
func DaemonVersionString() string {
	return fmt.Sprintf("device-hub.0.1.%s", SourceVersion)
}

// ClientVersionString returns the long version of the client executable
func ClientVersionString() string {
	return fmt.Sprintf("device-hub-cli.0.1.%s", SourceVersion)
}
