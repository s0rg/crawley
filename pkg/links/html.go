package links

import (
	"bufio"
	"bytes"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	keySRC    = "src"
	keySRCS   = "srcset"
	keyHREF   = "href"
	keyDATA   = "data"
	keyACTION = "action"
	keyPOSTER = "poster"
)

// HTMLHandler is a callback for found links.
type HTMLHandler func(atom.Atom, string)

// TokenFilter is a callback for token filtration.
type TokenFilter func(html.Token) bool

// HTMLParams holds config for ExtractHTML.
type HTMLParams struct {
	Filter  TokenFilter
	Handler HTMLHandler
	Brute   bool
}

// AllowALL - stub that implements TokenFilter, it allows all tokens.
func AllowALL(_ html.Token) bool { return true }

// ExtractHTML extract urls from html.
func ExtractHTML(r io.Reader, b *url.URL, p HTMLParams) {
	var (
		tkns = html.NewTokenizer(r)
		key  = keySRC
		t    html.Token
	)

	for {
		switch tkns.Next() {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			if t = tkns.Token(); p.Filter(t) {
				extractToken(b, t, &key, p.Handler)
			}
		case html.CommentToken:
			if p.Brute {
				extractComment(tkns.Token().Data, p.Handler)
			}
		}
	}
}

func extractComment(s string, h HTMLHandler) {
	ss := bufio.NewScanner(strings.NewReader(s))
	ss.Split(bufio.ScanWords)

	const (
		prefixHTTP  = "http://"
		prefixHTTPS = "https://"
		endBytes    = `<(')>"`
	)

	var (
		buf, low []byte
		pos, end int
	)

	for ss.Scan() {
		buf = ss.Bytes()
		low = bytes.ToLower(buf)

		if pos = bytes.Index(low, []byte(prefixHTTP)); pos == -1 {
			pos = bytes.Index(low, []byte(prefixHTTPS))
		}

		if pos == -1 {
			continue
		}

		if end = bytes.IndexAny(low[pos:], endBytes); end > -1 {
			buf = buf[:pos+end]
		}

		if uri := bytes.TrimSpace(buf[pos:]); len(uri) > 0 {
			suri := string(uri)
			if u, err := url.Parse(suri); err == nil && u.Host != "" {
				h(atom.A, suri)
			}
		}
	}
}

func extractToken(b *url.URL, t html.Token, k *string, h HTMLHandler) {
	var (
		poster string
		uri    string
	)

	switch t.DataAtom {
	case atom.A:
		uri = extractTag(b, &t, keyHREF)

	case atom.Img, atom.Image, atom.Iframe, atom.Script, atom.Track:
		uri = extractTag(b, &t, keySRC)

	case atom.Form:
		uri = extractTag(b, &t, keyACTION)

	case atom.Object:
		uri = extractTag(b, &t, keyDATA)

	case atom.Video:
		poster = extractTag(b, &t, keyPOSTER)

		fallthrough

	case atom.Audio:
		*k = keySRC
		uri = extractTag(b, &t, keySRC)

	case atom.Picture:
		*k = keySRCS

	case atom.Source:
		uri = extractTag(b, &t, *k)
	}

	callHandler(h, t.DataAtom, uri)
	callHandler(h, t.DataAtom, poster)
}

func callHandler(h HTMLHandler, a atom.Atom, s string) {
	if s != "" {
		h(a, s)
	}
}

func extractTag(base *url.URL, token *html.Token, key string) (rv string) {
	for i := 0; i < len(token.Attr); i++ {
		if a := &token.Attr[i]; a.Key == key {
			if res, ok := cleanURL(base, a.Val); ok {
				return res
			}
		}
	}

	return rv
}
