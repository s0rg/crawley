package crawler

import (
	"hash/fnv"
	"io"
	"net/url"
)

func urlHash(u *url.URL) (sum uint64) {
	c := *u         // copy original
	c.RawQuery = "" // remove any query parameters

	hash := fnv.New64()
	_, _ = io.WriteString(hash, c.String())

	return hash.Sum64()
}
