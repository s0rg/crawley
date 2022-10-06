package crawler

import (
	"log"
	"mime"
	"path/filepath"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/s0rg/crawley/pkg/links"
	"github.com/s0rg/crawley/pkg/set"
)

const (
	proxyAuthHdr = "Proxy-Authorization"
	proxyAuthTyp = "Basic"
	contentType  = "Content-Type"
	contentHTML  = "text/html"
	contentJS    = "application/javascript"
	fileExtJS    = ".js"
)

func isHTML(v string) (yes bool) {
	typ, _, err := mime.ParseMediaType(v)
	if err != nil {
		return
	}

	return typ == contentHTML
}

func isJS(v, n string) (yes bool) {
	typ, _, err := mime.ParseMediaType(v)
	if err != nil {
		return
	}

	if typ == contentJS {
		return true
	}

	return filepath.Ext(n) == fileExtJS
}

func prepareFilter(tags []string) links.TokenFilter {
	if len(tags) == 0 {
		return links.AllowALL
	}

	atoms := make(set.Set[atom.Atom])

	var a atom.Atom

	for _, t := range tags {
		if a = atom.Lookup([]byte(t)); a != 0 {
			atoms.Add(a)
		} else {
			log.Printf("[!] invalid tag: `%s` skipping...", t)
		}
	}

	return func(t html.Token) (ok bool) {
		return atoms.Has(t.DataAtom)
	}
}
