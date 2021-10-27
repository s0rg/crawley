[![Build](https://github.com/s0rg/crawley/workflows/ci/badge.svg)](https://github.com/s0rg/crawley/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/s0rg/crawley)](https://goreportcard.com/report/github.com/s0rg/crawley)
[![Maintainability](https://api.codeclimate.com/v1/badges/6542cd90a6c665e4202e/maintainability)](https://codeclimate.com/github/s0rg/crawley/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/e1c002df2b4571e01537/test_coverage)](https://codeclimate.com/github/s0rg/crawley/test_coverage)
[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/s0rg/crawley/blob/main/LICENSE)

# crawley

    The unix-way web crawler.
    It crawls web pages and prints any link it can found within.
    Scan depth (by default - 0) can be configured via cli flags.

# features

 - fast SAX-parser (powered by golang.org/x/net/html)
 - small (<1000 SLOC), idiomatic, 100% test covered codebase
 - grabs most of useful resources links (pics, videos, audios, etc...)
 - links are streamed to stdout

# usage
```
crawley [flags] url

possible flags:

-delay duration
    per-request delay
-depth int
    scan depth
-skip-ssl
    skip ssl verification
-user-agent string
    user-agent string
-version
    show version
-workers int
    number of workers
```
