package links

import (
	"encoding/xml"
	"io"
	"net/url"
)

type (
	// URLHandler is a callback for links.
	URLHandler func(string)

	entry struct {
		Loc string `xml:"loc"`
	}
)

// ExtractSitemap extract urls from sitemap*.xml.
func ExtractSitemap(r io.Reader, b *url.URL, h URLHandler) {
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

		if uri, ok = clean(b, e.Loc); !ok {
			continue
		}

		h(uri)
	}
}
