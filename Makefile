
SOURCE_VERSION = $(shell git rev-parse --short=6 HEAD)
BUILD_FLAGS = -v -ldflags "-X github.com/thingful/device-hub.SourceVersion=$(SOURCE_VERSION)"
PACKAGES := $(shell go list ./... | grep -v /vendor/ )

EXE_NAME := 'device-hub'
CLI_EXE_NAME := 'device-hub-cli'

GO_TEST = go test -covermode=atomic
GO_INTEGRATION = $(GO_TEST) -bench=. -v --tags=integration
GO_COVER = go tool cover
GO_BENCH = go test -bench=.
ARTEFACT_DIR = coverage

all: linux-arm linux-i386 linux-amd64 darwin-amd64 ## build executables for the various environments

.PHONY: all

get-build-deps: ## install build dependencies
	go get -u github.com/chespinoza/goliscan
	docker build -t thingful-device-hub-proto -f docker/Dockerfile.protobuf .

.PHONY: get-build-deps

check-license: ## check the license header in every code file
	@./scripts/check-license.sh

.PHONY: check-license

check-vendor-licenses: ## check if licenses of project dependencies meet project requirements 
	@goliscan check --direct-only -strict
	@goliscan check --indirect-only -strict
.PHONY: check-vendor-licenses

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
	docker run -v $(PWD)/proto:/go/proto thingful/device-hub-proto
	# strip `omitempty` from the json tags
	ls ./proto/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

.PHONY: proto

proto-verify: proto ## verify proto binding has been generated
	git diff --exit-code

.PHONY: proto-verify

docker-up: ## run dependencies as docker containers
	docker-compose up -d
	docker ps

.PHONY: docker_up

docker-build: linux-amd64  ## build a docker container containing the device-hub executables
	docker build -t thingful/device-hub:latest -t thingful/device-hub:$(SOURCE_VERSION) .

.PHONY: docker_build

darwin-amd64: tmp/build/darwin-amd64/$(EXE_NAME) tmp/build/darwin-amd64/$(CLI_EXE_NAME) ## build for mac amd64

linux-i386: tmp/build/linux-i386/$(EXE_NAME) tmp/build/linux-i386/$(CLI_EXE_NAME) ## build for linux i386

linux-amd64: tmp/build/linux-amd64/$(EXE_NAME) tmp/build/linux-amd64/$(CLI_EXE_NAME) ## build for linux amd64

linux-arm: tmp/build/linux-arm/$(EXE_NAME) tmp/build/linux-arm/$(CLI_EXE_NAME) ## build for linux arm (raspberry-pi)

.PHONY: darwin-amd64 linux-i386 linux-amd64 linux-arm

tmp/build/linux-i386/$(EXE_NAME):
	GOOS=linux GOARCH=386 go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub

tmp/build/linux-i386/$(CLI_EXE_NAME):
	GOOS=linux GOARCH=386 go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub-cli

tmp/build/linux-amd64/$(EXE_NAME):
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub

tmp/build/linux-amd64/$(CLI_EXE_NAME):
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub-cli

tmp/build/linux-arm/$(EXE_NAME):
	GOOS=linux GOARCH=arm go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub

tmp/build/linux-arm/$(CLI_EXE_NAME):
	GOOS=linux GOARCH=arm go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub-cli

tmp/build/darwin-amd64/$(EXE_NAME):
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub

tmp/build/darwin-amd64/$(CLI_EXE_NAME):
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(@) ./cmd/device-hub-cli


# 'help' parses the Makefile and displays the help text
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: help
