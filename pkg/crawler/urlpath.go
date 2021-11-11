package crawler

import (
	"net/url"
	"path"
	"strings"
)

func isPathSep(r rune) (yes bool) {
	const pathSep = '/'

	return r == pathSep
}

func canCrawl(a, b *url.URL, d int) (yes bool) {
	if a.Host != b.Host {
		return
	}

	var apath, bpath string

	if apath = a.Path; apath == "" {
		apath = "/"
	}

	if bpath = b.Path; bpath == "" {
		bpath = "/"
	}

	depth, found := relativeDepth(apath, bpath)
	if !found {
		return
	}

	if d >= 0 && depth > d {
		return
	}

	return true
}

func relativeDepth(base, sub string) (n int, ok bool) {
	var (
		bn = path.Clean(base)
		sn = path.Clean(sub)
	)

	if len(sn) <= len(bn) {
		return
	}

	if !strings.HasPrefix(sn, bn) {
		return
	}

	fields := strings.FieldsFunc(sn[len(bn):], isPathSep)

	for i := 0; i < len(fields); i++ {
		if fields[i] != "" {
			n++
		}
	}

	return n, true
}

func isResorce(v string) (yes bool) {
	_, tmp := path.Split(v)
	if tmp == "" {
		return
	}

	if tmp = path.Ext(tmp); tmp == "" {
		return
	}

	return true
}
