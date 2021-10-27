package path

import (
	"strings"
)

// Depth calculates relative depth for `sub` to `base` resorces path.
func Depth(base, sub string) (n int, ok bool) {
	var (
		bp = splitPath(base)
		sp = splitPath(sub)
	)

	if len(sp) < len(bp) {
		return
	}

	for i := 0; i < len(bp); i++ {
		if bp[i] != sp[i] {
			return
		}
	}

	return len(sp) - len(bp), true
}

func splitPath(p string) (o []string) {
	return dropSpaces(strings.Split(p, "/"))
}

func dropSpaces(s []string) (o []string) {
	o = make([]string, 0, len(s))

	for _, v := range s {
		if v != "" {
			o = append(o, v)
		}
	}

	return o
}
