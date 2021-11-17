[![Build](https://github.com/s0rg/crawley/workflows/ci/badge.svg)](https://github.com/s0rg/crawley/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/s0rg/crawley)](https://goreportcard.com/report/github.com/s0rg/crawley)
[![Maintainability](https://api.codeclimate.com/v1/badges/6542cd90a6c665e4202e/maintainability)](https://codeclimate.com/github/s0rg/crawley/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/e1c002df2b4571e01537/test_coverage)](https://codeclimate.com/github/s0rg/crawley/test_coverage)

[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/s0rg/crawley/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/s0rg/crawley)](go.mod)
[![Release](https://img.shields.io/github/v/release/s0rg/crawley)](https://github.com/s0rg/crawley/releases/latest)
![Downloads](https://img.shields.io/github/downloads/s0rg/crawley/total.svg)

# crawley

Crawls web pages and prints any link it can find.

# features

- fast html SAX-parser (powered by `golang.org/x/net/html`)
- small (<1000 SLOC), idiomatic, 100% test covered codebase
- grabs most of useful resources urls (pics, videos, audios, forms, etc...)
- found urls are streamed to stdout and guranteed to be unique (with fragments omitted)
- scan depth (limited by starting host and path, by default - 0) can be configured
- can crawl rules and sitemaps from `robots.txt`
- `brute` mode - scan html comments for urls (this can lead to bogus results)
- make use of `HTTP_PROXY` / `HTTPS_PROXY` environment values
- directory-only scan mode (aka `fast-scan`)

# installation

- [binaries](https://github.com/s0rg/crawley/releases) for Linux, FreeBSD, macOS and Windows

## Archlinux User Repository

Crawley is available in the AUR. Linux distributions with access to it can obtain the package from [here](https://aur.archlinux.org/packages/crawley-bin/).
You can also use your favourite AUR helper to install it, e. g. `paru -S crawley-bin`.

# usage

```
crawley [flags] url

possible flags:

-brute
    scan html comments
-delay duration
    per-request delay (0 - disable) (default 150ms)
-depth int
    scan depth (-1 - unlimited)
-dirs
    policy for non-resource urls: show / hide / only (default "show")
-headless
    disable pre-flight HEAD requests
-help
    this flags (and their defaults) description
-robots string
    policy for robots.txt: ignore / crawl / respect (default "ignore")
-silent
    suppress info and error messages in stderr
-skip-ssl
    skip ssl verification
-user-agent string
    user-agent string
-version
    show version
-workers int
    number of workers (default - number of CPU cores)
```
