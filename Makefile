VERSION := $(shell git describe --tags 2>/dev/null || echo dev)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

build:
	go build -o ${PROJECTNAME} ${LDFLAGS} cmd/ddns-cloudflare/main.go

test:
	go test -v -race ./...
