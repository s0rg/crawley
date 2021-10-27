SHELL=/bin/bash

BIN=bin/crawley
SRC=./cmd/crawley
COP=test.coverage

GIT_HASH=`git rev-parse --short HEAD`
BUILD_AT=`date +%FT%T%z`

LDFLAGS=-w -s -X main.GitHash=${GIT_HASH} -X main.BuildDate=${BUILD_AT}

.PHONY: build

export CGO_ENABLED=0
export GOARCH=amd64

build: build-linux

build-linux: lint
	GOOS=linux go build -ldflags "${LDFLAGS}" -o "${BIN}" "${SRC}"

build-windows: lint
	GOOS=windows go build -ldflags "${LDFLAGS}" -o "${BIN}".exe "${SRC}"

build-darwin: lint
	GOOS=darwin go build -ldflags "${LDFLAGS}" -o "${BIN}".osx "${SRC}"

vet:
	go vet ./...

lint: vet
	golangci-lint run

test: vet
	go test -race -count 1 -v -tags=test -coverprofile="${COP}" ./...

test-cover: test
	go tool cover -func="${COP}"

clean:
	[ -f "${BIN}" ] && rm "${BIN}"
	[ -f "${BIN}".exe ] && rm "${BIN}".exe
	[ -f "${BIN}".osx ] && rm "${BIN}".osx
	[ -f "${COP}" ] && rm "${COP}"
