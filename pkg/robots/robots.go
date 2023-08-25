package robots

import (
	"io"
	"net/url"

	"github.com/s0rg/set"
)

const path = "/robots.txt"

type accessMode byte

const (
	allowAll accessMode = 0
	gotRules accessMode = 1
	denyAll  accessMode = 2
)

// TXT holds parsed robots.txt contents and/or access mode.
type TXT struct {
	links    set.Set[string]
	deny     set.Set[string]
	sitemaps set.Set[string]
	mode     accessMode
}

// AllowALL returns empty TXT that allows all.
func AllowALL() *TXT {
	return &TXT{}
}

// DenyALL returns empty TXT that denies all.
func DenyALL() *TXT {
	return &TXT{mode: denyAll}
}

// FromReader parse robots.txt body from given reader.
func FromReader(ua string, r io.Reader) (t *TXT, err error) {
	t = &TXT{
		mode:     gotRules,
		links:    make(set.Unordered[string]),
		deny:     make(set.Unordered[string]),
		sitemaps: make(set.Unordered[string]),
	}

	if err = parseRobots(r, ua, t); err != nil {
		return
	}

	return t, nil
}

// URL construct robots.txt url for given host.
func URL(u *url.URL) (rv string) {
	var t url.URL

	t.Scheme = u.Scheme
	t.Host = u.Host
	t.Path = path

	return t.String()
}

// Forbidden checks if path is forbidden by given rules and mode.
func (t *TXT) Forbidden(path string) (yes bool) {
	switch t.mode {
	case gotRules:
		yes = t.deny.Has(path)
	case denyAll:
		yes = true
	case allowAll:
	}

	return yes
}

// Sitemaps returns list of parsed sitemaps urls.
func (t *TXT) Sitemaps() (rv []string) {
	return set.ToSlice(t.sitemaps)
}

// Links returns list of parsed rules urls (both allowed and denided, for all useragents).
func (t *TXT) Links() (rv []string) {
	return set.ToSlice(t.links)
}
