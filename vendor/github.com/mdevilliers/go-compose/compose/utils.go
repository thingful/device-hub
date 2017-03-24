package compose

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
)

var dockerHostRegexp = regexp.MustCompile("://([^:]+):")

// InferDockerHost returns the current docker host based on the contents of the DOCKER_HOST environment variable.
// If DOCKER_HOST is not set, it returns "localhost".
func InferDockerHost() (string, error) {
	envHost := os.Getenv("DOCKER_HOST")
	if len(envHost) == 0 {
		return "localhost", nil
	}

	matches := dockerHostRegexp.FindAllStringSubmatch(envHost, -1)
	if len(matches) != 1 || len(matches[0]) != 2 {
		return "", fmt.Errorf("compose: cannot parse DOCKER_HOST '%v'", envHost)
	}
	return matches[0][1], nil
}

// MustInferDockerHost is like InferDockerHost, but panics on error.
func MustInferDockerHost() string {
	dockerHost, err := InferDockerHost()
	if err != nil {
		panic(err)
	}
	return dockerHost
}

func runCmd(name string, args ...string) (string, error) {
	var outBuf bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf
	err := cmd.Run()
	out := outBuf.String()
	if err != nil {
		fmt.Print(out)
	}
	return out, err
}

func writeTmp(content string) (string, error) {
	f, err := ioutil.TempFile("", "docker-compose-")

	if err != nil {
		return "", fmt.Errorf("compose: error creating temp file: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return "", fmt.Errorf("compose: error writing temp file: %v", err)
	}

	return f.Name(), nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
