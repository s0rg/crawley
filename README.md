[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/s0rg/crawley/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/s0rg/crawley)](go.mod)
[![Release](https://img.shields.io/github/v/release/s0rg/crawley)](https://github.com/s0rg/crawley/releases/latest)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
![Downloads](https://img.shields.io/github/downloads/s0rg/crawley/total.svg)

[![CI](https://github.com/s0rg/crawley/workflows/ci/badge.svg)](https://github.com/s0rg/crawley/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/s0rg/crawley)](https://goreportcard.com/report/github.com/s0rg/crawley)
[![Maintainability](https://api.codeclimate.com/v1/badges/6542cd90a6c665e4202e/maintainability)](https://codeclimate.com/github/s0rg/crawley/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/e1c002df2b4571e01537/test_coverage)](https://codeclimate.com/github/s0rg/crawley/test_coverage)
[![libraries.io](https://img.shields.io/librariesio/github/s0rg/crawley)](https://libraries.io/github/s0rg/crawley)
![Issues](https://img.shields.io/github/issues/s0rg/crawley)

# crawley

Crawls web pages and prints any link it can find.

# features

- fast html SAX-parser (powered by `golang.org/x/net/html`)
- small (<1500 SLOC), idiomatic, 100% test covered codebase
- grabs most of useful resources urls (pics, videos, audios, forms, etc...)
- found urls are streamed to stdout and guranteed to be unique (with fragments omitted)
- scan depth (limited by starting host and path, by default - 0) can be configured
- can crawl rules and sitemaps from `robots.txt`
- `brute` mode - scan html comments for urls (this can lead to bogus results)
- make use of `HTTP_PROXY` / `HTTPS_PROXY` environment values + handles proxy auth
- directory-only scan mode (aka `fast-scan`)
- user-defined cookies, in curl-compatible format (i.e. `-cookie "ONE=1; TWO=2" -cookie "ITS=ME" -cookie @cookie-file`)
- user-defined headers, same as curl: `-header "ONE: 1" -header "TWO: 2" -header @headers-file`
- tag filter - allow to specify tags to crawl for (single: `-tag a -tag form`, multiple: `-tag a,form`, or mixed)
- url ignore - allow to ignore urls with matched substrings from crawling (i.e.: '-ignore logout')
- js parser - extract api endpoints from js files, this done by regexp, so results can be messy

# examples
```sh
# print all links from first page:
crawley http://some-test.site

# print all js files and api endpoints:
crawley -depth -1 -tag script -js http://some-test.site

# print all endpoints from js:
crawley -js http://some-test.site/app.js

# download all png images from site:
crawley -depth -1 -tag img http://some-test.site | grep '\.png$' | wget -i -

# fast directory traversal:
crawley -headless -delay 0 -depth -1 -dirs only http://some-test.site
```

# installation

- [binaries](https://github.com/s0rg/crawley/releases) for Linux, FreeBSD, macOS and Windows, just download and run.
- [archlinux](https://aur.archlinux.org/packages/crawley-bin/) you can use your favourite AUR helper to install it, e. g. `paru -S crawley-bin`.

# usage

```
crawley [flags] url

possible flags:

-brute
    scan html comments
-cookie value
    extra cookies for request, can be used multiple times, accept files with '@'-prefix
-delay duration
    per-request delay (0 - disable) (default 150ms)
-depth int
    scan depth (-1 - unlimited)
-dirs string
    policy for non-resource urls: show / hide / only (default "show")
-header value
    extra headers for request, can be used multiple times, accept files with '@'-prefix
-headless
    disable pre-flight HEAD requests
-help
    this flags (and their defaults) description
-ignore value
    patterns (in urls) to be ignored in crawl process
-js
    scan js files for endpoints
-proxy-auth string
    credentials for proxy: user:password
-robots string
    policy for robots.txt: ignore / crawl / respect (default "ignore")
-silent
    suppress info and error messages in stderr
-skip-ssl
    skip ssl verification
-tag value
    tags filter, single or comma-separated tag names allowed
-user-agent string
    user-agent string
-version
    show version
-workers int
    number of workers (default - number of CPU cores)
```
