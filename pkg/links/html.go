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
	jsScheme  = "javascript"
)

// Handler is a callback for found links.
type Handler func(atom.Atom, string)

// ExtractHTML run `handler` for every link found inside html from `r`, rebasing them to `b` (if need).
func ExtractHTML(b *url.URL, r io.ReadCloser, brute bool, handler Handler) {
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
			base(b, tkns.Token(), &key, handler)
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
			suri := string(uri)
			if u, err := url.Parse(suri); err == nil && u.Host != "" {
				h(atom.A, suri)
			}
		}
	}
}

func base(b *url.URL, t html.Token, k *string, h Handler) {
	var (
		res = make([]string, 0, 2)
		uri string
		ok  bool
	)

	switch t.DataAtom {
	case atom.A:
		uri, ok = extractTag(b, &t, keyHREF)

	case atom.Img, atom.Image, atom.Iframe, atom.Script, atom.Track:
		uri, ok = extractTag(b, &t, keySRC)

	case atom.Form:
		uri, ok = extractTag(b, &t, keyACTION)

	case atom.Object:
		uri, ok = extractTag(b, &t, keyDATA)

	case atom.Video:
		if uri, ok = extractTag(b, &t, keyPOSTER); ok {
			res = append(res, uri)
		}

		fallthrough

	case atom.Audio:
		*k = keySRC
		uri, ok = extractTag(b, &t, keySRC)

	case atom.Picture:
		*k = keySRCS

	case atom.Source:
		uri, ok = extractTag(b, &t, *k)
	}

	if ok {
		res = append(res, uri)
	}

	for i := 0; i < len(res); i++ {
		h(t.DataAtom, res[i])
	}
}

func extractTag(base *url.URL, token *html.Token, key string) (rv string, ok bool) {
	for i := 0; i < len(token.Attr); i++ {
		if a := &token.Attr[i]; a.Key == key {
			return clean(base, a.Val)
		}
	}

	return rv, false
}

func clean(base *url.URL, link string) (rv string, ok bool) {
	u, err := url.Parse(link)
	if err != nil {
		return
	}

	if u.Host == "" {
		if u = base.ResolveReference(u); u.Host == "" {
			return
		}
	}

	switch u.Scheme {
	case jsScheme:
		return
	case "":
		u.Scheme = base.Scheme
	}

	if u.Path == "" {
		u.Path = "/"
	}

	u.Fragment = ""

	return u.String(), true
}
