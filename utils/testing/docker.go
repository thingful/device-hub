// Copyright © 2017 thingful
package testing

import "os"

var (
	werkerEnvVar = "WERCKER_STEP_NAME"
)

// IsRunningInWercker tests for a well know environmental variable and returns true if found, false if absent
func IsRunningInWercker() bool {
	return os.Getenv(werkerEnvVar) != ""
}
