-include .env

COMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags ${COMMIT})
PROJECTNAME := $(shell basename "$(PWD)")
IMG ?= xerac/paper-soccer-gui:${VERSION}


# Go related variables.
GOBASE := $(shell pwd)
GOFILES := $(shell find $(GOBASE) -type f -name "*.go")
GOMAIN := .
BINDIR := $(GOBASE)/bin

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags '-X=main.Version=$(VERSION) -X=main.Build=$(COMMIT) -extldflags="-static"'

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: fmt dep get test clean build run exec

fmt:
	@echo ">>>  Formatting project"
	go fmt ./...

dep:
	@echo ">>>  Add missing and remove unused modules..."
	go mod tidy

gen:
	@echo ">>>  Generating assets.go"
	go run gen.go

get: dep
	@echo ">>>  Checking if there is any missing dependencies..."
	go get -u ./...

test: build clean
	@echo ">>>  Testing..."
	go test ./...

clean:
	@echo ">>>  Cleaning build cache"
	-rm -r $(BINDIR) 2> /dev/null
	go clean ./...

build: gen
	@echo ">>>  Building binary..."
	mkdir -p $(BINDIR) 2> /dev/null
	go build $(LDFLAGS) -o $(BINDIR)/$(PROJECTNAME) $(GOMAIN)

run:
	@echo ">>>  Running..."
	go run $(GOMAIN)

exec: build
	@echo ">>>  Executing binary..."
	@$(BINDIR)/$(PROJECTNAME)

docker-build: build
	@echo ">>>  Building docker image..."
	docker build --no-cache -t $(IMG) .

docker-push: docker-build
	@echo ">>>  Pushing docker image..."
	docker push $(IMG)