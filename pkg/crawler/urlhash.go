package crawler

import (
	"hash/fnv"
	"io"
	"strings"
)

func urlHash(s string) (sum uint64) {
	hash := fnv.New64()
	_, _ = io.WriteString(hash, strings.ToLower(s))

	return hash.Sum64()
}
