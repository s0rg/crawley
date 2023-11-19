package client

import "strings"

const (
	headerParts     = 2
	headerSeparator = ":"
)

type header struct {
	Key string
	Val string
}

func prepareHeaders(raw []string) (rv []*header) {
	rv = make([]*header, 0, len(raw))

	var (
		pair     []string
		key, val string
	)

	for _, h := range raw {
		pair = strings.SplitN(h, headerSeparator, headerParts)

		if key = strings.TrimSpace(pair[0]); key == "" {
			continue
		}

		if val = strings.TrimSpace(pair[1]); val == "" {
			continue
		}

		rv = append(rv, &header{Key: key, Val: val})
	}

	return rv
}
