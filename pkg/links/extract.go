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
	jsScheme  = "javascript"
)

// Handler is a callback for found links.
type Handler func(a atom.Atom, u *url.URL)

// Extract run `handler` for every link found inside html from `r`, rebasing them to `b` (if need).
func Extract(b *url.URL, r io.ReadCloser, brute bool, handler Handler) {
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
		case html.CommentToken:
			if brute {
				extractComment(tkns.Token().Data, handler)
			}
		}
	}
}

func extractComment(s string, h Handler) {
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
			if u, err := url.Parse(string(uri)); err == nil && u.Host != "" {
				h(atom.A, u)
			}
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

	case atom.Form:
		u, ok = extractTag(b, &t, keyACTION)

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
