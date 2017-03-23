package compose

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

var goodYML = `
test_mockserver:
  container_name: ms
  image: jamesdbloom/mockserver
  ports:
    - "10000:1080"
    - "1090"
test_postgres:
  container_name: pg
  image: postgres
  ports:
    - "5432"
`

var badYML = `
bad
`

func TestGoodYML(t *testing.T) {
	compose := MustStart(goodYML, true, true)
	defer compose.MustKill()

	if compose.Containers["ms"].Name != "/ms" {
		t.Fatalf("found name '%v', expected '/ms", compose.Containers["ms"].Name)
	}
	if compose.Containers["pg"].Name != "/pg" {
		t.Fatalf("found name '%v', expected '/pg", compose.Containers["pg"].Name)
	}
	if port := compose.Containers["ms"].MustGetFirstPublicPort(1080, "tcp"); port != 10000 {
		t.Fatalf("found port %v, expected 10000", port)
	}
}

func TestRestartGoodYML(t *testing.T) {
	TestGoodYML(t)
}

func TestBadYML(t *testing.T) {
	compose, err := Start(badYML, true, true)
	if err == nil {
		defer compose.MustKill()
		t.Error("expected error")
	}
}

func TestMustInferDockerHost(t *testing.T) {
	envHost := os.Getenv("DOCKER_HOST")
	defer os.Setenv("DOCKER_HOST", envHost)

	os.Setenv("DOCKER_HOST", "")
	if host := MustInferDockerHost(); host != "localhost" {
		t.Errorf("found '%v', expected 'localhost'", host)
	}
	os.Setenv("DOCKER_HOST", "tcp://192.168.99.100:2376")
	if host := MustInferDockerHost(); host != "192.168.99.100" {
		t.Errorf("found '%v', expected '192.168.99.100'", host)
	}
}

func TestMustConnectWithDefaults(t *testing.T) {
	compose := MustStart(goodYML, true, true)
	defer compose.MustKill()

	mockServerURL := fmt.Sprintf("http://%v:%v", MustInferDockerHost(), compose.Containers["ms"].MustGetFirstPublicPort(1080, "tcp"))

	MustConnectWithDefaults(func() error {
		logger.Print("attempting to connect to mockserver...")
		_, err := http.Get(mockServerURL)
		if err == nil {
			logger.Print("connected to mockserver")
		}
		return err
	})
}

func TestInspectUnknownContainer(t *testing.T) {
	_, err := Inspect("bad")
	if err == nil {
		t.Error("expected error")
	}
}

func TestMustInspect(t *testing.T) {
	compose := MustStart(goodYML, true, true)
	defer compose.MustKill()

	ms := MustInspect(compose.Containers["ms"].ID)
	if ms.Name != "/ms" {
		t.Errorf("found '%v', expected '/ms", ms.Name)
	}
}
