package links

import (
	"encoding/xml"
	"io"
	"net/url"
)

type SitemapHandler func(string)

type entry struct {
	Loc string `xml:"loc"`
}

func ExtractSitemap(b *url.URL, r io.ReadCloser, handler SitemapHandler) {
	dec := xml.NewDecoder(r)

	var (
		t   xml.Token
		e   entry
		se  xml.StartElement
		uri *url.URL
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

		if uri, ok = clean(b, e.Loc); !ok {
			continue
		}

		handler(uri.String())
	}
}
