package client

import (
	"log"
	"net/http"
	"strings"
)

const (
	keyvalParts     = 2
	keyvalSeparator = "="
	valuesSeparator = ";"
)

func prepareCookies(raw []string) (rv []*http.Cookie) {
	for _, r := range raw {
		for _, p := range strings.Split(r, valuesSeparator) {
			if val, ok := parseOne(p); ok {
				rv = append(rv, val)
			} else {
				log.Printf("cannot parse '%s' as cookie, expected format: 'key=value;' as in curl", r)
			}
		}
	}

	return rv
}

func parseOne(raw string) (rv *http.Cookie, ok bool) {
	pair := strings.SplitN(raw, keyvalSeparator, keyvalParts)
	if len(pair) != keyvalParts {
		return
	}

	var name, value string

	if name = strings.TrimSpace(pair[0]); name == "" {
		return
	}

	if value = strings.TrimSpace(pair[1]); value == "" {
		return
	}

	rv = &http.Cookie{
		Name:  name,
		Value: value,
	}

	return rv, true
}
