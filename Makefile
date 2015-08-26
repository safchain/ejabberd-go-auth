.PHONY: all clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install

build:
	@go build $(GOFLAGS) ./...

install: build
	@go get $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...
