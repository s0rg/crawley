package links

import (
	"bytes"
	"io"
	"regexp"
	"strings"
)

const (
	regexAPI = `(?:"|')` +
		`(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|` +
		`((?:/|\.\./|\./)[^"'><,;| *()%$^/\\\[\]][^"'><,;|()]{1,})|` +
		`([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|` +
		`([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|` +
		`([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml|cgi)(?:[\?|#][^"|']{0,}|)))(?:"|')`
	mimeAppPrefix  = "application/"
	mimeTxtPrefix  = "text/"
	codeCleanChars = `"'`
)

var reJSAPI = regexp.MustCompile(regexAPI)

// ExtractJS extract urls from js files.
func ExtractJS(r io.Reader, h URLHandler) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return
	}

	res := reJSAPI.FindAll(buf, -1)
	for i := 0; i < len(res); i++ {
		if uri, ok := cleanResult(res[i]); ok {
			h(uri)
		}
	}
}

func cleanResult(s []byte) (rv string, ok bool) {
	rv = string(bytes.Trim(s, codeCleanChars))

	if strings.HasPrefix(rv, mimeAppPrefix) || strings.HasPrefix(rv, mimeTxtPrefix) {
		return "", false
	}

	return rv, true
}
