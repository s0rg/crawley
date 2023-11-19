package links

import (
	"bytes"
	"io"
	"strings"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

const (
	codeCleanChars = `"'`
	dash           = "/"
	doubleDash     = dash + dash
)

// ExtractJS extract urls from js files.
func ExtractJS(r io.Reader, h URLHandler) {
	l := js.NewLexer(parse.NewInput(r))

	for {
		tt, text := l.Next()
		switch tt {
		case js.ErrorToken:
			return
		case js.StringToken:
			if res, ok := extractJSURL(text); ok {
				h(res)
			}
		}
	}
}

func extractJSURL(v []byte) (rv string, ok bool) {
	rv = string(bytes.Trim(v, codeCleanChars))
	ok = strings.HasPrefix(rv, dash) || strings.Contains(rv, doubleDash)

	return rv, ok
}
