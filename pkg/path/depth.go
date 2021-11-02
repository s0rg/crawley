package path

import (
	"path"
	"strings"
)

const pathSep = '/'

func isPathSep(r rune) (yes bool) {
	return r == pathSep
}

// Depth calculates relative depth for `sub` to `base` resorces path.
func Depth(base, sub string) (n int, ok bool) {
	var (
		bn = path.Clean(base)
		sn = path.Clean(sub)
	)

	if len(sn) <= len(bn) {
		return
	}

	if !strings.HasPrefix(sn, bn) {
		return
	}

	fields := strings.FieldsFunc(sn[len(bn):], isPathSep)

	for i := 0; i < len(fields); i++ {
		if fields[i] != "" {
			n++
		}
	}

	return n, true
}
