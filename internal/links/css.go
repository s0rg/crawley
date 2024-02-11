package links

import (
	"bytes"
	"io"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

const (
	braceOpen  = '('
	braceClose = ')'
)

// ExtractCSS extract urls from css files.
func ExtractCSS(r io.Reader, h URLHandler) {
	l := css.NewLexer(parse.NewInput(r))

	for {
		switch tt, text := l.Next(); tt {
		case css.ErrorToken:
			return
		case css.URLToken:
			if res, ok := extractCSSURL(text); ok {
				h(res)
			}
		}
	}
}

func extractCSSURL(v []byte) (rv string, ok bool) {
	o := bytes.IndexByte(v, braceOpen)
	c := bytes.LastIndexByte(v, braceClose)

	rv = string(bytes.Trim(v[o+1:c], codeCleanChars))

	return rv, rv != ""
}
