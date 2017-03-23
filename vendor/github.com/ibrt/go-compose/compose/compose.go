/*
Package compose provides a Go wrapper around Docker Compose, useful for integration testing.

	// Define Compose config.
	var composeYML =`
	test_mockserver:
	  container_name: ms
	  image: jamesdbloom/mockserver
	  ports:
	    - "10000:1080"
	    - "${SOME_ENV_VAR}" # This is replaced with the value of SOME_ENV_VAR.
	test_postgres:
	  container_name: pg
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
		c.Containers["ms"].MustGetFirstPublicPort(1080, "tcp"))

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
	fileName   string
	Containers map[string]*Container
}

var (
	logger           = log.New(os.Stdout, "go-compose: ", log.LstdFlags)
	replaceEnvRegexp = regexp.MustCompile("\\$\\{[^\\}]+\\}")
	composeUpRegexp  = regexp.MustCompile("(?m:docker start <- \\(u'(.*)'\\)$)")
)

const (
	composeProjectName = "compose"
)

// Start starts a Docker Compose configuration.
// If forcePull is true, it attempts do pull newer versions of the images.
// If rmFirst is true, it attempts to kill and delete containers before starting new ones.
func Start(dockerComposeYML string, forcePull, rmFirst bool) (*Compose, error) {
	logger.Println("initializing...")
	dockerComposeYML = replaceEnv(dockerComposeYML)

	fName, err := writeTmp(dockerComposeYML)
	if err != nil {
		return nil, err
	}

	ids, err := composeStart(fName, forcePull, rmFirst)
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
		containers[container.Name[1:]] = container
	}

	return &Compose{fileName: fName, Containers: containers}, nil
}

// MustStart is like Start, but panics on error.
func MustStart(dockerComposeYML string, forcePull, killFirst bool) *Compose {
	compose, err := Start(dockerComposeYML, forcePull, killFirst)
	if err != nil {
		panic(err)
	}
	return compose
}

// Kill kills any running containers for the current configuration.
func (c *Compose) Kill() error {
	return composeKill(c.fileName)
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

func composeStart(fName string, forcePull, rmFirst bool) ([]string, error) {
	if forcePull {
		logger.Println("pulling images...")
		if _, err := composeRun(fName, "pull"); err != nil {
			return nil, fmt.Errorf("compose: error pulling images: %v", err)
		}
	}

	if rmFirst {
		if err := composeKill(fName); err != nil {
			return nil, err
		}
		if err := composeRm(fName); err != nil {
			return nil, err
		}
	}

	logger.Println("starting containers...")
	out, err := composeRun(fName, "--verbose", "up", "-d")
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

func composeKill(fName string) error {
	logger.Println("killing stale containers...")
	_, err := composeRun(fName, "kill")
	if err != nil {
		return fmt.Errorf("compose: error killing stale containers: %v", err)
	}
	return err
}

func composeRm(fName string) error {
	logger.Println("removing stale containers...")
	_, err := composeRun(fName, "rm", "--force")
	if err != nil {
		return fmt.Errorf("compose: error removing stale containers: %v", err)
	}
	return err
}

func composeRun(fName string, otherArgs ...string) (string, error) {
	args := []string{"-f", fName, "-p", composeProjectName}
	args = append(args, otherArgs...)
	return runCmd("docker-compose", args...)
}
