/*
Package compose provides a Go wrapper around Docker Compose, useful for integration testing.

	// Define Compose config.
	var composeYML =`
	test_mockserver:
	  image: jamesdbloom/mockserver
	  ports:
	    - "10000:1080"
	    - "${SOME_ENV_VAR}" # This is replaced with the value of SOME_ENV_VAR.
	test_postgres:
	  image: postgres
	  ports:
	    - "5432"
	`

	// Start containers.
	c, err := compose.Start(composeYML, true, true)
	if err != nil {
		panic(err)
	}
	defer c.Kill()

	// Build MockServer public URL.
	mockServerURL := fmt.Sprintf(
		"http://%v:%v",
		compose.MustInferDockerHost(),
		c.Containers["test_mockserver"].MustGetFirstPublicPort(1080, "tcp"))

	// Wait for MockServer to start accepting connections.
	MustConnectWithDefaults(func() error {
		_, err := http.Get(mockServerURL)
		return err
	})
	...
*/
package compose

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// Compose is the main type exported by the package, used to interact with a running Docker Compose configuration.
type Compose struct {
	fileName           string
	composeProjectName string
	Containers         map[string]*Container
}

var (
	logger           = log.New(os.Stdout, "go-compose: ", log.LstdFlags)
	replaceEnvRegexp = regexp.MustCompile("\\$\\{[^\\}]+\\}")
	composeUpRegexp  = regexp.MustCompile("(?m:docker start <- \\(u'(.*)'\\)$)")
)

// Start starts a Docker Compose configuration.
// Fixes the Docker Compose project name to a known value so existing containers can be killed.
// If forcePull is true, it attempts do pull newer versions of the images.
// If rmFirst is true, it attempts to kill and delete containers before starting new ones.
func Start(dockerComposeYML string, forcePull, rmFirst bool) (*Compose, error) {
	return StartProject(dockerComposeYML, forcePull, rmFirst, "compose")
}

// StartParallel starts a Docker Compose configuration and is suitable for concurrent usage.
// The project name is defined at random to ensure multiple instances can be run.
// Note: that the docker services should not bind to localhost ports.
func StartParallel(dockerComposeYML string, forcePull bool) (*Compose, error) {
	return StartProject(dockerComposeYML, forcePull, false, randStringBytes(9))
}

// StartProject starts a Docker Compose configuration, giving fine grained control of all of the properties.
func StartProject(dockerComposeYML string, forcePull, rmFirst bool, projectName string) (*Compose, error) {

	logger.Println("initializing...")

	dockerComposeYML = replaceEnv(dockerComposeYML)

	fName, err := writeTmp(dockerComposeYML)
	if err != nil {
		return nil, err
	}

	ids, err := composeStart(fName, projectName, forcePull, rmFirst)
	if err != nil {
		return nil, err
	}

	containers := make(map[string]*Container)

	for _, id := range ids {
		container, err := Inspect(id)
		if err != nil {
			return nil, err
		}
		if !container.State.Running {
			return nil, fmt.Errorf("compose: container '%v' is not running", container.Name)
		}
		containers[container.Config.Labels["com.docker.compose.service"]] = container
	}

	return &Compose{fileName: fName, composeProjectName: projectName, Containers: containers}, nil
}

// MustStart is like Start, but panics on error.
func MustStart(dockerComposeYML string, forcePull, killFirst bool) *Compose {
	compose, err := Start(dockerComposeYML, forcePull, killFirst)
	if err != nil {
		panic(err)
	}
	return compose
}

// MustStartParallel is like StartParallel, but panics on error.
func MustStartParallel(dockerComposeYML string, forcePull bool) *Compose {
	compose, err := StartParallel(dockerComposeYML, forcePull)
	if err != nil {
		panic(err)
	}
	return compose
}

// Kill kills any running containers for the current configuration.
func (c *Compose) Kill() error {
	return composeKill(c.fileName, c.composeProjectName)
}

// MustKill is like Kill, but panics on error.
func (c *Compose) MustKill() {
	if err := c.Kill(); err != nil {
		panic(err)
	}
}

func replaceEnv(dockerComposeYML string) string {
	return replaceEnvRegexp.ReplaceAllStringFunc(dockerComposeYML, replaceEnvFunc)
}

func replaceEnvFunc(s string) string {
	return os.Getenv(strings.TrimSpace(s[2 : len(s)-1]))
}

func composeStart(fName, composeProjectName string, forcePull, rmFirst bool) ([]string, error) {
	if forcePull {
		logger.Println("pulling images...")
		if _, err := composeRun(fName, composeProjectName, "pull"); err != nil {
			return nil, fmt.Errorf("compose: error pulling images: %v", err)
		}
	}

	if rmFirst {
		if err := composeKill(fName, composeProjectName); err != nil {
			return nil, err
		}
		if err := composeRm(fName, composeProjectName); err != nil {
			return nil, err
		}
	}

	logger.Println("starting containers...")
	out, err := composeRun(fName, composeProjectName, "--verbose", "up", "-d")
	if err != nil {
		return nil, fmt.Errorf("compose: error starting containers: %v", err)
	}
	logger.Println("containers started")

	matches := composeUpRegexp.FindAllStringSubmatch(out, -1)
	ids := make([]string, 0, len(matches))
	for _, match := range matches {
		ids = append(ids, match[1])
	}

	return ids, nil
}

func composeKill(fName, composeProjectName string) error {
	logger.Println("killing stale containers...")
	_, err := composeRun(fName, composeProjectName, "kill")
	if err != nil {
		return fmt.Errorf("compose: error killing stale containers: %v", err)
	}
	return err
}

func composeRm(fName, composeProjectName string) error {
	logger.Println("removing stale containers...")
	_, err := composeRun(fName, composeProjectName, "rm", "--force")
	if err != nil {
		return fmt.Errorf("compose: error removing stale containers: %v", err)
	}
	return err
}

func composeRun(fName, composeProjectName string, otherArgs ...string) (string, error) {
	args := []string{"-f", fName, "-p", composeProjectName}
	args = append(args, otherArgs...)
	return runCmd("docker-compose", args...)
}
