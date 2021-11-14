SHELL=/bin/bash

BIN=bin/crawley
SRC=./cmd/crawley
COP=test.coverage

GIT_TAG=`git describe --abbrev=0 2>/dev/null || echo -n "no-tag"`
GIT_HASH=`git rev-parse --short HEAD 2>/dev/null || echo -n "no-git"`
BUILD_AT=`date +%FT%T%z`

LDFLAGS=-w -s -X main.gitHash=${GIT_HASH} -X main.buildDate=${BUILD_AT} -X main.gitVersion=${GIT_TAG}

export CGO_ENABLED=0

.PHONY: build

build: build-linux

build-linux: vet
	GOOS=linux go build -ldflags "${LDFLAGS}" -o "${BIN}" "${SRC}"

build-freebsd: vet
	GOOS=freebsd go build -ldflags "${LDFLAGS}" -o "${BIN}".bin "${SRC}"

build-windows: vet
	GOOS=windows go build -ldflags "${LDFLAGS}" -o "${BIN}".exe "${SRC}"

build-darwin: vet
	GOOS=darwin go build -ldflags "${LDFLAGS}" -o "${BIN}".osx "${SRC}"

vet:
	go vet ./...

lint: vet
	golangci-lint run

test: vet
	CGO_ENABLED=1 go test -race -count 1 -v -tags=test -coverprofile="${COP}" ./...

test-cover: test
	go tool cover -func="${COP}"

clean:
	[ -f "${COP}" ] && rm "${COP}"
	[ -f "${BIN}" ] && rm "${BIN}"
	[ -f "${BIN}".bin ] && rm "${BIN}".bin
	[ -f "${BIN}".exe ] && rm "${BIN}".exe
	[ -f "${BIN}".osx ] && rm "${BIN}".osx
