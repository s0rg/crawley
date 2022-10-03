package set

import (
	"hash/fnv"
	"io"
	"strings"
)

// URI holds set of uint64 hashes.
type URI map[uint64]stub

func (u URI) Add(v string) (ok bool) {
	h := hash(v)

	if _, ok = u[h]; ok {
		return false
	}

	u[h] = stub{}

	return true
}

func hash(s string) (rv uint64) {
	hash := fnv.New64()
	_, _ = io.WriteString(hash, strings.ToLower(s))

	return hash.Sum64()
}
