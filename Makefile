SHELL := /usr/bin/env bash

OUT := crawley
ALL := ./...
BIN := ./bin/${OUT}
SRC := ./cmd/${OUT}
COP := cover.out
SELF = $(CURDIR)/$(lastword $(MAKEFILE_LIST))

GIT_TAG := `git describe --abbrev=0 2>/dev/null || echo -n "no-tag"`
GIT_REV := `git rev-parse --short HEAD 2>/dev/null || echo -n "no-git"`
BUILD_AT := `date +%FT%T%z`
LDFLAGS := -w -s \
		  -X main.GitTag=${GIT_TAG} \
		  -X main.GitHash=${GIT_REV} \
		  -X main.BuildDate=${BUILD_AT}

COMPILER := go build
CCFLAGS := ${COMPILER} -ldflags "${LDFLAGS}"

TESTER := go test
TSTFLAGS := ${TESTER} -race -count 1 -vet=off -tags=test -v

export CGO_ENABLED=0

.PHONY: $(wildcard *)

## help: Prints this help message
help:
	@echo -e "\nUsage:"
	@sed -n 's/^##//p' "${SELF}" | column -t -s ':' | sed -e 's/^/\t/'

## build: Default build action - for Linux
build: build/linux

## build/linux: Builds for Linux
build/linux: code/vet
	@echo "Building for Linux..."
	@GOOS=linux ${CCFLAGS} -o "${BIN}" "${SRC}"

## build/freebsd: Builds for FreeBSD
build/freebsd: code/vet
	@echo "Building for FreeBSD..."
	@GOOS=freebsd ${CCFLAGS} -o "${BIN}".bin "${SRC}"

## build/windows: Builds for Windows
build/windows: code/vet
	@echo "Building for Windows..."
	@GOOS=windows ${CCFLAGS} -o "${BIN}".exe "${SRC}"

## build/darwin: Builds for MacOS
build/darwin: code/vet
	@echo "Building for MacOS..."
	@GOOS=darwin ${CCFLAGS} -o "${BIN}".osx "${SRC}"

## code/vet: Performs basic linting for code
code/vet:
	@echo "Running go vet..."
	@go vet "${ALL}"

## code/lint: Performs advanced linting for code
code/lint: code/vet
	@echo "Running golangci-lint..."
	@golangci-lint run

## code/test-all: Runs tests suite
code/test-all: code/vet
	@echo "Running all tests"
	@CGO_ENABLED=1 ${TSTFLAGS} -tags test -coverprofile="${COP}" "${ALL}"

## code/test [name]: Runs specified test
code/test: code/vet
	@echo "Running single test: ${name}"
	@CGO_ENABLED=1 ${TSTFLAGS} -run ${name} "${ALL}"

## code/test-cover: Runs test-coverage
code/test-cover: code/vet
	@echo "Running test covarage"
	@go test -v -coverprofile="$(COP)" -tags test -cover ./... -coverpkg ./... -covermode=count
	@go tool cover -func="$(COP)" -o="$(COP)"

## code/benchmark: Runs code benchmarks
code/benchmark: code/test
	@echo "Running benchmarks..."
	@CGO_ENABLED=1 ${TESTER} -v -benchmem -bench=${ALL}

## code/coverage: Calculates overall test coverage
code/coverage: code/test-all
	@echo "Calculating code coverage..."
	@go tool cover -func="${COP}"

## code/clean: Performs clean-up
code/clean:
	@echo "Cleaning up..."
	[ -f "${COP}" ] && rm "${COP}"
	[ -f "${BIN}" ] && rm "${BIN}"
	[ -f "${BIN}".bin ] && rm "${BIN}".bin
	[ -f "${BIN}".exe ] && rm "${BIN}".exe
	[ -f "${BIN}".osx ] && rm "${BIN}".osx
