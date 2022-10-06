package links

import (
	"io"
	"net/url"
	"regexp"
)

var reJSURL = regexp.MustCompile(`
  (?:"|')
  (
    ((?:[a-zA-Z]{1,10}:/|/)
    [^"'/]{1,}\.
    [a-zA-Z]{2,}[^"']{0,})
    |
    ((?:/|\.\./|\./)
    [^"'><,;| *()(%$^/\\[\]]
    [^"'><,;|()]{1,})
    |
    ([a-zA-Z0-9_\-/]{1,}/
    [a-zA-Z0-9_\-/]{1,}
    \.(?:[a-zA-Z]{1,4}|action)
    (?:[\?|#][^"|']{0,}|))
    |
    ([a-zA-Z0-9_\-/]{1,}/
    [a-zA-Z0-9_\-/]{3,}
    (?:[\?|#][^"|']{0,}|))
    |
    ([a-zA-Z0-9_\-]{1,}
    \.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)
    (?:[\?|#][^"|']{0,}|))
  )
  (?:"|')`)

// ExtractJS extract urls from js files.
func ExtractJS(r io.Reader, b *url.URL, h URLHandler) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return
	}

	res := reJSURL.FindAll(buf, -1)
	for i := 0; i < len(res); i++ {
		//if uri, ok := clean(b, string(res[i])); ok {
		//	h(uri)
		//}
		h(string(res[i]))
	}
}
