package pipe

import hub "github.com/thingful/device-hub"

type Channel interface {
	Errors() chan error
	Out() chan hub.Input
}
