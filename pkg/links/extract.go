package links

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	keySRC   = "src"
	keySRCS  = "srcset"
	keyHREF  = "href"
	keyDATA  = "data"
	jsScheme = "javascript"
)

// Handler is a callback for found links.
type Handler func(a atom.Atom, u *url.URL)

// Extract run `handler` for every link found inside html from `r`, rebasing them to `b` (if need).
func Extract(b *url.URL, r io.ReadCloser, handler Handler) {
	defer r.Close()

	var (
		tkns = html.NewTokenizer(r)
		key  = keySRC
	)

	for {
		switch tkns.Next() {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			extractToken(b, tkns.Token(), &key, handler)
		}
	}
}

func extractToken(b *url.URL, t html.Token, k *string, h Handler) {
	var (
		u  *url.URL
		ok bool
	)

	switch t.DataAtom {
	case atom.A:
		u, ok = extractTag(b, &t, keyHREF)

	case atom.Img, atom.Image, atom.Iframe, atom.Script, atom.Track:
		u, ok = extractTag(b, &t, keySRC)

	case atom.Object:
		u, ok = extractTag(b, &t, keyDATA)

	case atom.Video, atom.Audio:
		*k = keySRC
		u, ok = extractTag(b, &t, keySRC)

	case atom.Picture:
		*k = keySRCS

	case atom.Source:
		u, ok = extractTag(b, &t, *k)
	}

	if ok {
		h(t.DataAtom, u)
	}
}

func extractTag(b *url.URL, t *html.Token, k string) (u *url.URL, ok bool) {
	for i := 0; i < len(t.Attr); i++ {
		if a := &t.Attr[i]; a.Key == k {
			return clean(b, a.Val)
		}
	}

	return nil, false
}

func clean(b *url.URL, r string) (u *url.URL, ok bool) {
	u, err := url.Parse(r)
	if err != nil {
		return
	}

	if u.Host == "" {
		if u = b.ResolveReference(u); u.Host == "" {
			return
		}
	}

	switch u.Scheme {
	case jsScheme:
		return
	case "":
		u.Scheme = b.Scheme
	}

	if u.Path == "" {
		u.Path = "/"
	}

	u.Fragment = ""

	return u, true
}
