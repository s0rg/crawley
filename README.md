[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/s0rg/crawley/blob/main/LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fs0rg%2Fcrawley.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fs0rg%2Fcrawley?ref=badge_shield)
[![Go Version](https://img.shields.io/github/go-mod/go-version/s0rg/crawley)](go.mod)
[![Release](https://img.shields.io/github/v/release/s0rg/crawley)](https://github.com/s0rg/crawley/releases/latest)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
![Downloads](https://img.shields.io/github/downloads/s0rg/crawley/total.svg)

[![CI](https://github.com/s0rg/crawley/workflows/ci/badge.svg)](https://github.com/s0rg/crawley/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/s0rg/crawley)](https://goreportcard.com/report/github.com/s0rg/crawley)
[![Maintainability](https://qlty.sh/badges/f6bbc710-32a4-430b-ba73-51ae05fd0916/maintainability.svg)](https://qlty.sh/gh/s0rg/projects/crawley)
[![Code Coverage](https://qlty.sh/badges/f6bbc710-32a4-430b-ba73-51ae05fd0916/test_coverage.svg)](https://qlty.sh/gh/s0rg/projects/crawley)
[![libraries.io](https://img.shields.io/librariesio/github/s0rg/crawley)](https://libraries.io/github/s0rg/crawley)
![Issues](https://img.shields.io/github/issues/s0rg/crawley)

# crawley

Crawls web pages and prints any link it can find.


# features

- fast html SAX-parser (powered by [x/net/html](https://golang.org/x/net/html))
- js/css lexical parsers (powered by [tdewolff/parse](https://github.com/tdewolff/parse)) - extract api endpoints from js code and `url()` properties
- small (below 1500 SLOC), idiomatic, 100% test covered codebase
- grabs most of useful resources urls (pics, videos, audios, forms, etc...)
- found urls are streamed to stdout and guranteed to be unique (with fragments omitted)
- scan depth (limited by starting host and path, by default - 0) can be configured
- can be polite - crawl rules and sitemaps from `robots.txt`
- `brute` mode - scan html comments for urls (this can lead to bogus results)
- make use of `HTTP_PROXY` / `HTTPS_PROXY` environment values + handles proxy auth (use `HTTP_PROXY="socks5://127.0.0.1:1080/" crawley` for socks5)
- directory-only scan mode (aka `fast-scan`)
- user-defined cookies, in curl-compatible format (i.e. `-cookie "ONE=1; TWO=2" -cookie "ITS=ME" -cookie @cookie-file`)
- user-defined headers, same as curl: `-header "ONE: 1" -header "TWO: 2" -header @headers-file`
- tag filter - allow to specify tags to crawl for (single: `-tag a -tag form`, multiple: `-tag a,form`, or mixed)
- url ignore - allow to ignore urls with matched substrings from crawling (i.e.: `-ignore logout`)
- subdomains support - allow depth crawling for subdomains as well (e.g. `crawley http://some-test.site` will be able to crawl `http://www.some-test.site`)


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

- [binaries / deb / rpm](https://github.com/s0rg/crawley/releases) for Linux, FreeBSD, macOS and Windows.
- [archlinux](https://aur.archlinux.org/packages/crawley-bin/) you can use your favourite AUR helper to install it, e. g. `paru -S crawley-bin`.


# usage

```
crawley [flags] url

possible flags with default values:

-all
    scan all known sources (js/css/...)
-brute
    scan html comments
-cookie value
    extra cookies for request, can be used multiple times, accept files with '@'-prefix
-css
    scan css for urls
-delay duration
    per-request delay (0 - disable) (default 150ms)
-depth int
    scan depth (set -1 for unlimited)
-dirs string
    policy for non-resource urls: show / hide / only (default "show")
-header value
    extra headers for request, can be used multiple times, accept files with '@'-prefix
-headless
    disable pre-flight HEAD requests
-ignore value
    patterns (in urls) to be ignored in crawl process
-js
    scan js code for endpoints
-proxy-auth string
    credentials for proxy: user:password
-robots string
    policy for robots.txt: ignore / crawl / respect (default "ignore")
-silent
    suppress info and error messages in stderr
-skip-ssl
    skip ssl verification
-subdomains
    support subdomains (e.g. if www.domain.com found, recurse over it)
-tag value
    tags filter, single or comma-separated tag names
-timeout duration
    request timeout (min: 1 second, max: 10 minutes) (default 5s)
-user-agent string
    user-agent string
-version
    show version
-workers int
      number of workers (default - number of CPU cores)
-ignore-query
    ignore query parameters in URL
```


# flags autocompletion

Crawley can handle flags autocompletion in bash and zsh via `complete`:

```bash
complete -C "/full-path-to/bin/crawley" crawley
```


# license
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fs0rg%2Fcrawley.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fs0rg%2Fcrawley?ref=badge_large)
