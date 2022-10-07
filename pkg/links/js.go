package links

import (
	"bufio"
	"bytes"
	"io"
	"net/url"
	"regexp"
)

var reJSURL = regexp.MustCompile(`\"/?[a-zA-Z0-9_/?=&.+ ]*\"`)

const (
	codeSplitChars = "\n;"
	codeCleanChars = `"'`
)

// ExtractJS extract urls from js files.
func ExtractJS(r io.Reader, b *url.URL, h URLHandler) {
	s := bufio.NewScanner(r)
	s.Split(splitCode)

	var res [][]byte

	for s.Scan() {
		res = reJSURL.FindAll(s.Bytes(), -1)
		for i := 0; i < len(res); i++ {
			if uri, ok := clean(b, cleanResult(res[i])); ok {
				h(uri)
			}
		}
	}
}

func splitCode(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexAny(data, codeSplitChars); i >= 0 {
		return i + 1, cleanSplit(data[0:i]), nil
	}

	if atEOF {
		return len(data), cleanSplit(data), nil
	}

	return 0, nil, nil
}

func cleanSplit(s []byte) (rv []byte) {
	return bytes.Trim(s, codeSplitChars)
}

func cleanResult(s []byte) (rv string) {
	return string(bytes.Trim(s, codeCleanChars))
}
