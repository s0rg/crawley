SHELL=/usr/bin/env bash

BIN=bin/crawley
SRC=./cmd/crawley
COP=cover.out
SELF=$(CURDIR)/$(lastword $(MAKEFILE_LIST))

GIT_TAG=`git describe --abbrev=0 2>/dev/null || echo -n "no-tag"`
GIT_HASH=`git rev-parse --short HEAD 2>/dev/null || echo -n "no-git"`
BUILD_AT=`date +%FT%T%z`

LDFLAGS=-w -s \
		-X main.GitTag=${GIT_TAG} \
		-X main.GitHash=${GIT_HASH} \
		-X main.BuildDate=${BUILD_AT}

export CGO_ENABLED=0

.PHONY: $(wildcard *)

## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' "${SELF}" | column -t -s ':'

## build: Default build action, for Linux
build: build/linux

## build/linux: Builds for Linux
build/linux: vet
	@GOOS=linux go build -ldflags "${LDFLAGS}" -o "${BIN}" "${SRC}"

## build/freebsd: Builds for FreeBSD
build/freebsd: vet
	@GOOS=freebsd go build -ldflags "${LDFLAGS}" -o "${BIN}".bin "${SRC}"

## build/windows: Builds for Windows
build/windows: vet
	@GOOS=windows go build -ldflags "${LDFLAGS}" -o "${BIN}".exe "${SRC}"

## build/darwin: Builds for MacOS
build/darwin: vet
	@GOOS=darwin go build -ldflags "${LDFLAGS}" -o "${BIN}".osx "${SRC}"

## code/vet: Performs basic linting for code
code/vet:
	@go vet ./...

## code/lint: Performs advanced linting for code
code/lint: code/vet
	@golangci-lint run

## code/test-all: Runs tests suite
code/test-all: code/vet
	@CGO_ENABLED=1 go test -race -count 1 -vet=off -tags=test -coverprofile="${COP}" -v ./...

## code/test [name]: Runs specified test
code/test: code/vet
	@CGO_ENABLED=1 go test -race -count 1 -vet=off -tags=test -v -run ${name} ./...

## code/coverage: Calculates overall test coverage
code/coverage: code/test-all
	@go tool cover -func="${COP}"

## code/clean: Performs clean-up
code/clean:
	[ -f "${COP}" ] && rm "${COP}"
	[ -f "${BIN}" ] && rm "${BIN}"
	[ -f "${BIN}".bin ] && rm "${BIN}".bin
	[ -f "${BIN}".exe ] && rm "${BIN}".exe
	[ -f "${BIN}".osx ] && rm "${BIN}".osx
