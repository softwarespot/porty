# If the "VERSION" environment variable is not set, then use this value instead
VERSION?=1.0.0
TIME=$(shell date +%FT%T%z)
GOVERSION=$(shell go version | awk '{print $$3}' | sed s/go//)

LDFLAGS=-ldflags "\
	-X github.com/softwarespot/porty/internal/version.Version=${VERSION} \
	-X github.com/softwarespot/porty/internal/version.Time=${TIME} \
	-X github.com/softwarespot/porty/internal/version.User=${USER} \
	-X github.com/softwarespot/porty/internal/version.GoVersion=${GOVERSION} \
	-s \
	-w \
"

build:
	@echo building to bin/porty
	@go build $(LDFLAGS) -o ./bin/porty

install:
	@echo copying bin/porty to $(HOME)/bin/porty
	@mv ./bin/porty $(HOME)/bin/porty

test:
	@go test -cover -v ./...

.PHONY: build install test
