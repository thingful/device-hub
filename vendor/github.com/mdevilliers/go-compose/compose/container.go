package compose

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Container models the `docker inspect` command output.
type Container struct {
	ID              string           `json:"Id"`
	Name            string           `json:"Name,omitempty"`
	Created         time.Time        `json:"Created,omitempty"`
	Config          *Config          `json:"Config,omitempty"`
	State           State            `json:"State,omitempty"`
	Image           string           `json:"Image,omitempty"`
	NetworkSettings *NetworkSettings `json:"NetworkSettings,omitempty"`
}

// Config models the config section of the `docker inspect` command output.
type Config struct {
	Hostname     string              `json:"Hostname,omitempty"`
	ExposedPorts map[string]struct{} `json:"ExposedPorts,omitempty"`
	Env          []string            `json:"Env,omitempty"`
	Cmd          []string            `json:"Cmd"`
	Image        string              `json:"Image,omitempty"`
	Labels       map[string]string   `json:"Labels,omitempty"`
}

// State models the state section of the `docker inspect` command.
type State struct {
	Running    bool      `json:"Running,omitempty"`
	Paused     bool      `json:"Paused,omitempty"`
	Restarting bool      `json:"Restarting,omitempty"`
	OOMKilled  bool      `json:"OOMKilled,omitempty"`
	Pid        int       `json:"Pid,omitempty"`
	ExitCode   int       `json:"ExitCode,omitempty"`
	Error      string    `json:"Error,omitempty"`
	StartedAt  time.Time `json:"StartedAt,omitempty"`
	FinishedAt time.Time `json:"FinishedAt,omitempty"`
}

// NetworkSettings models the network settings section of the `docker inspect` command.
type NetworkSettings struct {
	Ports map[string][]PortBinding `json:"Ports,omitempty"`
}

// PortBinding models a port binding in the network settings section of the `docker inspect command.
type PortBinding struct {
	HostIP   string `json:"HostIP,omitempty"`
	HostPort string `json:"HostPort,omitempty"`
}

// Inspect inspects a container using the `docker inspect` command and returns a parsed version of its output.
func Inspect(id string) (*Container, error) {
	out, err := runCmd("docker", "inspect", id)
	if err != nil {
		return nil, fmt.Errorf("compose: error inspecting container: %v", err)
	}

	var inspect []*Container
	if err := json.Unmarshal([]byte(out), &inspect); err != nil {
		return nil, fmt.Errorf("compose: error parsing inspect output: %v", err)
	}
	if len(inspect) != 1 {
		return nil, fmt.Errorf("compose: inspect returned %v results, 1 expected", len(inspect))
	}

	return inspect[0], nil
}

// MustInspect is like Inspect, but panics on error.
func MustInspect(id string) *Container {
	container, err := Inspect(id)
	if err != nil {
		panic(err)
	}
	return container
}

// GetFirstPublicPort returns the first public public port mapped to the given exposedPort, for the given proto ("tcp", "udp", etc.), if found.
func (c *Container) GetFirstPublicPort(exposedPort uint32, proto string) (uint32, error) {
	if c.NetworkSettings == nil {
		return 0, fmt.Errorf("compose: no network settings for container '%v'", c.Name)
	}

	portSpec := fmt.Sprintf("%v/%v", exposedPort, strings.ToLower(proto))
	mapping, ok := c.NetworkSettings.Ports[portSpec]
	if !ok || len(mapping) == 0 {
		return 0, fmt.Errorf("compose: no public port for %v", portSpec)
	}

	port, err := strconv.ParseUint(mapping[0].HostPort, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("compose: error parsing port '%v'", mapping[0].HostPort)
	}

	return uint32(port), nil
}

// MustGetFirstPublicPort is like GetFirstPublicPort, but panics on error.
func (c *Container) MustGetFirstPublicPort(exposedPort uint32, proto string) uint32 {
	port, err := c.GetFirstPublicPort(exposedPort, proto)
	if err != nil {
		panic(err)
	}
	return port
}
