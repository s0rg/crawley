package crawler

import (
	"net/url"
	"path"
	"strings"

	"github.com/s0rg/crawley/pkg/set"
)

var parsableExts = make(set.Set[string]).Load(
	"asp",
	"aspx",
	"cgi",
	"jsp",
	"html",
	"pl",
	"php",
	"xhtml",
	"xml",
)

func isSitemap(s string) (yes bool) {
	const (
		sitemapXML = "sitemap.xml"
		sitemapIDX = "sitemap-index.xml"
	)

	switch {
	case strings.HasSuffix(s, sitemapXML), strings.HasSuffix(s, sitemapIDX):
		yes = true
	}

	return
}

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

func canParse(v string) (yes bool) {
	_, tmp := path.Split(v)
	if tmp == "" {
		return true
	}

	if tmp = path.Ext(tmp); tmp == "" {
		return true
	}

	tmp = strings.ToLower(tmp[1:])

	return parsableExts.Has(tmp)
}
