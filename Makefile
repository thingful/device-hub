
SOURCE_VERSION = $(shell git rev-parse --short=6 HEAD)
BUILD_FLAGS = -v -ldflags "-X main.SourceVersion=$(SOURCE_VERSION)"
PACKAGES := $(shell go list ./... | grep -v /vendor/ )

EXE_NAME := 'device-hub'

all: pi linux darwin ## build executables for the various environments

.PHONY: all

test: ## run the tests
	go test -v $(PACKAGES)

.PHONY: test

clean: ## clean up
	rm -rf tmp/

.PHONY: clean

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

