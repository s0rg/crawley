package robots

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type tokenKind byte

const (
	kindNone      tokenKind = 0
	kindUserAgent tokenKind = 1
	kindAllow     tokenKind = 2
	kindDisallow  tokenKind = 3
	kindSitemap   tokenKind = 4
)

const (
	defaultAgent  = "*"
	tokenSep      = ':'
	tokenComment  = '#'
	tokenAllow    = "allow"
	tokenDisallow = "disallow"
	tokenSitemap1 = "sitemap"
	tokenSitemap2 = "site-map"
	tokenUA1      = "useragent"
	tokenUA2      = "user-agent"
)

func parseTokenKind(b []byte) (k tokenKind) {
	switch {
	case bytes.EqualFold(b, []byte(tokenUA1)), bytes.EqualFold(b, []byte(tokenUA2)):
		k = kindUserAgent
	case bytes.EqualFold(b, []byte(tokenAllow)):
		k = kindAllow
	case bytes.EqualFold(b, []byte(tokenDisallow)):
		k = kindDisallow
	case bytes.EqualFold(b, []byte(tokenSitemap1)), bytes.EqualFold(b, []byte(tokenSitemap2)):
		k = kindSitemap
	}

	return k
}

func extractToken(b []byte) (k tokenKind, v string) {
	var pos int
	if pos = bytes.IndexByte(b, tokenComment); pos >= 0 {
		b = b[:pos] // cut-off comments (if any)
	}

	b = bytes.TrimSpace(b)

	if pos = bytes.IndexByte(b, tokenSep); pos == -1 {
		return
	}

	key := bytes.TrimSpace(b[:pos])
	kind := parseTokenKind(key)

	if kind == kindNone {
		return
	}

	var val []byte
	if val = bytes.TrimSpace(b[pos+1:]); len(val) == 0 {
		return
	}

	return kind, string(val)
}

func parseRobots(r io.Reader, ua string, t *TXT) (err error) {
	var (
		s    = bufio.NewScanner(r)
		deny bool
	)

	s.Split(bufio.ScanLines)

	for s.Scan() {
		switch k, v := extractToken(s.Bytes()); k {
		case kindUserAgent:
			deny = (v == defaultAgent || strings.Contains(ua, v))

		case kindDisallow:
			if deny {
				t.deny.Add(v)
			}

			fallthrough

		case kindAllow:
			t.links.Add(v)

		case kindSitemap:
			t.sitemaps.Add(v)
		}
	}

	if e := s.Err(); e != nil {
		return fmt.Errorf("scanner: %w", e)
	}

	return nil
}
