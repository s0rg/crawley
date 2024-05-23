package crawler

import (
	"encoding/base64"
	"hash/fnv"
	"io"
	"log"
	"mime"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/s0rg/set"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/s0rg/crawley/internal/links"
)

const (
	proxyAuthKey   = "Proxy-Authorization"
	proxyAuthBasic = "Basic"

	contentType = "Content-Type"
	contentHTML = "text/html"
	contentCSS  = "text/css"
	contentJS   = "application/javascript"
	fileExtJS   = ".js"
	fileExtCSS  = ".css"
)

var parsableExts = set.Load(make(set.Unordered[string]),
	".asp",
	".aspx",
	".cgi",
	".htm",
	".html",
	".jsp",
	".php",
	".pl",
	".xhtml",
	".xml",
	fileExtJS,
	fileExtCSS,
)

func proxyAuthHeader(v string) (rv string) {
	return proxyAuthKey + ": " + proxyAuthBasic + " " + base64.StdEncoding.EncodeToString([]byte(v))
}

func prepareFilter(tags []string) links.TokenFilter {
	if len(tags) == 0 {
		return links.AllowALL
	}

	atoms := make(set.Unordered[atom.Atom])

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

func canCrawl(a, b *url.URL, d int, subdomains bool) (yes bool) {
	if a.Host != b.Host {
		if subdomains{
			domainA := strings.Split(a.Host, ".")
			domainB := strings.Split(b.Host, ".")
			if len(domainA) >= len(domainB){
				// The base domain must be shorter than the found domain
				return
			}
			j := len(domainB) - 1
			for i := len(domainA) - 1; i >= 0 && j >= 0; i-- {
				// Traverse each domain from the end, to check if their top-level domain are the same
				if domainA[i] != domainB[j] {
					// not the same top-level host
					return
				}
				j--
			}
		} else{
			return
		}
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

	if len(sn) < len(bn) {
		return
	}

	if !strings.HasPrefix(sn, bn) {
		return
	}

	const pathSep = '/'

	fields := strings.FieldsFunc(sn[len(bn):], func(r rune) bool {
		return r == pathSep
	})

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

	return strings.HasSuffix(s, sitemapXML) || strings.HasSuffix(s, sitemapIDX)
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

func isCSS(v, n string) (yes bool) {
	typ, _, err := mime.ParseMediaType(v)
	if err == nil && typ == contentCSS {
		return true
	}

	return webExt(n) == fileExtCSS
}

func urlhash(s string) (rv uint64) {
	hash := fnv.New64()
	_, _ = io.WriteString(hash, strings.ToLower(s))

	return hash.Sum64()
}
