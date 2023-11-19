package links

import (
	"net/url"
	"strings"
	"testing"
)

func TestExtractSitemap(t *testing.T) {
	t.Parallel()

	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://HOST/</loc>
  </url>
  <url>
    <loc>http://HOST/tools/</loc>
    <lastmod>2015-05-07T19:13:09+09:00</lastmod>
  </url>
  <url>
    <loc>http://HOST/contribution-to-oss/</loc>
    <lastmod>2015-05-07</lastmod>
    <changefreq>monthly</changefreq>
  </url>
  <url>
    <loc>http://HOST/page-1/</loc>
    <lastmod>2015-05-07T19:13:09+09:00</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.9</priority>
  </url>
</urlset>`

	u, _ := url.Parse("http://HOST")

	l := make([]string, 0, 4)

	ExtractSitemap(strings.NewReader(xml), u, func(s string) {
		l = append(l, s)
	})

	if len(l) != 4 {
		t.Error("unexpected results count")
	}
}

func TestExtractSitemapIndex(t *testing.T) {
	t.Parallel()

	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap>
    <loc>http://www.example.com/sitemap1.xml.gz</loc>
    <lastmod>2004-10-01T18:23:17+00:00</lastmod>
  </sitemap>
  <sitemap>
    <loc>http://www.example.com/sitemap2.xml.gz</loc>
    <lastmod>2005-01-01</lastmod>
  </sitemap>
  <sitemap>
    <loc>http://www.example.com/sitemap3.xml.gz</loc>
  </sitemap>
</sitemapindex>`

	u, _ := url.Parse("http://www.example.com")

	l := make([]string, 0, 3)

	ExtractSitemap(strings.NewReader(xml), u, func(s string) {
		l = append(l, s)
	})

	if len(l) != 3 {
		t.Error("unexpected results count")
	}
}

func TestExtractSitemapTokenError(t *testing.T) {
	t.Parallel()

	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap>
    <loc>http://www.example.com/sitemap1.xml.gz</loc>
    <last
`

	u, _ := url.Parse("http://www.example.com")

	l := make([]string, 0, 1)

	ExtractSitemap(strings.NewReader(xml), u, func(s string) {
		l = append(l, s)
	})

	if len(l) != 0 {
		t.Error("unexpected results count")
	}
}

func TestExtractSitemapURLError(t *testing.T) {
	t.Parallel()

	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap>
    <loc>[%]</loc>
  </sitemap>
</sitemapindex>`

	u, _ := url.Parse("http://www.example.com")

	l := make([]string, 0, 1)

	ExtractSitemap(strings.NewReader(xml), u, func(s string) {
		l = append(l, s)
	})

	if len(l) != 0 {
		t.Error("unexpected results count")
	}
}
