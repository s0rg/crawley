package robots

import (
	"io"
	"net/url"

	"github.com/s0rg/crawley/pkg/set"
)

const path = "/robots.txt"

type TXT struct {
	links    set.String
	deny     set.String
	sitemaps set.String
	mode     accessMode
}

func AllowALL() *TXT {
	return &TXT{}
}

func DenyALL() *TXT {
	return &TXT{mode: denyAll}
}

func FromReader(ua string, r io.Reader) (t *TXT, err error) {
	t = &TXT{
		mode:     gotRules,
		links:    make(set.String),
		deny:     make(set.String),
		sitemaps: make(set.String),
	}

	if err = parseRobots(r, ua, t); err != nil {
		return
	}

	return t, nil
}

func URL(u *url.URL) (rv string) {
	var t url.URL

	t.Scheme = u.Scheme
	t.Host = u.Host
	t.Path = path

	return t.String()
}

func (t *TXT) Forbidden(path string) (yes bool) {
	switch t.mode {
	case gotRules:
		_, yes = t.deny[path]
	case denyAll:
		yes = true
	case allowAll:
	}

	return yes
}

func (t *TXT) Sitemaps() (rv []string) {
	return t.sitemaps.List()
}

func (t *TXT) Links() (rv []string) {
	return t.links.List()
}
