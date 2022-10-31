package crawler

import (
	"hash/fnv"
	"io"
	"log"
	"mime"
	"net/url"
	"path"
	"path/filepath"
	"strings"

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

var parsableExts = make(set.Set[string]).Load(
	".asp",
	".aspx",
	".cgi",
	".jsp",
	".html",
	".pl",
	".php",
	".xhtml",
	".xml",
	".js",
)

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

func webExt(v string) (rv string) {
	const maxCut = 2

	p := strings.SplitN(v, "?", maxCut)

	return filepath.Ext(p[0])
}

func canParse(v string) (yes bool) {
	_, tmp := path.Split(v)
	if tmp == "" {
		return true
	}

	if tmp = webExt(tmp); tmp == "" {
		return true
	}

	return parsableExts.Has(strings.ToLower(tmp))
}

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

func isResorce(v string) (yes bool) {
	_, tmp := path.Split(v)
	if tmp == "" {
		return
	}

	if tmp = webExt(tmp); tmp == "" {
		return
	}

	return true
}

func isHTML(v string) (yes bool) {
	typ, _, err := mime.ParseMediaType(v)
	if err != nil {
		return
	}

	return typ == contentHTML
}

func isJS(v, n string) (yes bool) {
	typ, _, err := mime.ParseMediaType(v)
	if err == nil && typ == contentJS {
		return true
	}

	return webExt(n) == fileExtJS
}

func urlhash(s string) (rv uint64) {
	hash := fnv.New64()
	_, _ = io.WriteString(hash, strings.ToLower(s))

	return hash.Sum64()
}
