
SOURCE_VERSION = $(shell git rev-parse --short=6 HEAD)
BUILD_FLAGS = -v -ldflags "-X github.com/thingful/device-hub.SourceVersion=$(SOURCE_VERSION)"
PACKAGES := $(shell go list ./... | grep -v /vendor/ )

EXE_NAME := 'device-hub'

GO_TEST = go test -covermode=atomic
GO_INTEGRATION = $(GO_TEST) -bench=. -v --tags=integration
GO_COVER = go tool cover
GO_BENCH = go test -bench=.
ARTEFACT_DIR = coverage

all: pi linux darwin ## build executables for the various environments

.PHONY: all

test: ## run tests
	$(GO_TEST) $(PACKAGES)

.PHONY: test

test_integration: ## run integration tests (SLOW)
	mkdir -p $(ARTEFACT_DIR)
	echo 'mode: atomic' > $(ARTEFACT_DIR)/cover-integration.out
	touch $(ARTEFACT_DIR)/cover.tmp
	$(foreach package, $(PACKAGES), $(GO_INTEGRATION) -coverprofile=$(ARTEFACT_DIR)/cover.tmp $(package) && tail -n +2 $(ARTEFACT_DIR)/cover.tmp >> $(ARTEFACT_DIR)/cover-integration.out || exit;)
.PHONY: test_integration

clean: ## clean up
	rm -rf tmp/
	rm -rf $(ARTEFACT_DIR)

.PHONY: clean

bench: ## run benchmark tests
	$(GO_BENCH) $(PACKAGES)

.PHONY: bench

coverage: test_integration ## generate and display coverage report
	$(GO_COVER) -func=$(ARTEFACT_DIR)/cover-integration.out

.PHONY: test_integration 

proto: ## regenerate protobuf files
	protoc --go_out=plugins=grpc:. ./proto/*.proto

.PHONY: proto

docker_up: ## run dependencies as docker containers
	docker-compose up -d
	docker ps

.PHONY: docker_up

pi : tmp/build/$(EXE_NAME)-linux-arm
darwin: tmp/build/$(EXE_NAME)-darwin-amd64
linux: tmp/build/$(EXE_NAME)-linux-amd64

tmp/build/$(EXE_NAME)-linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub

tmp/build/$(EXE_NAME)-linux-arm:
	GOOS=linux GOARCH=arm go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub

tmp/build/$(EXE_NAME)-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub

# 'help' parses the Makefile and displays the help text
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: help

