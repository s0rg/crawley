package links

import (
	"encoding/xml"
	"io"
	"net/url"
)

type (
	// SitemapHandler is a callback for links found in sitemap.
	SitemapHandler func(string)

	entry struct {
		Loc string `xml:"loc"`
	}
)

// ExtractSitemap run `handler` for every link found inside sitemap from `r`, rebasing them to `b` (if need).
func ExtractSitemap(base *url.URL, r io.Reader, handler SitemapHandler) {
	var (
		dec = xml.NewDecoder(r)
		t   xml.Token
		e   entry
		se  xml.StartElement
		uri string
		err error
		ok  bool
	)

	for {
		if t, err = dec.Token(); err != nil {
			break
		}

		if se, ok = t.(xml.StartElement); !ok {
			continue
		}

		switch se.Name.Local {
		default:
			continue
		case "url", "sitemap":
		}

		if err := dec.DecodeElement(&e, &se); err != nil {
			continue
		}

		if uri, ok = clean(base, e.Loc); !ok {
			continue
		}

		handler(uri)
	}
}
