package robots

import (
	"io"

	"github.com/s0rg/crawley/pkg/set"
)

type accessMode byte

const (
	allowAll accessMode = 0
	gotRules accessMode = 1
	denyAll  accessMode = 2
)

type TXT struct {
	links    set.String
	deny     set.String
	sitemaps set.String
	ua       string
	mode     accessMode
}

func NewAllowAll() *TXT {
	return &TXT{}
}

func NewDenyAll() *TXT {
	return &TXT{mode: denyAll}
}

func NewFromReader(ua string, r io.Reader) (t *TXT, err error) {
	t = &TXT{
		mode:  gotRules,
		ua:    ua,
		links: make(set.String),
		deny:  make(set.String),
	}

	if err = parseRobots(r, t); err != nil {
		return
	}

	return t, nil
}

func (t *TXT) CanAccess(path string) (yes bool) {
	switch t.mode {
	case gotRules:
		if _, ok := t.deny[path]; !ok {
			yes = true
		}
	case allowAll:
		yes = true
	case denyAll:
	}

	return yes
}

func (t *TXT) Sitemaps() (rv []string) {
	return t.sitemaps.List()
}

func (t *TXT) Links() (rv []string) {
	return t.links.List()
}
