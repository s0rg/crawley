package crawler

import (
	"testing"
)

func TestUrlHash(t *testing.T) {
	t.Parallel()

	one := "http://test/some/path?foo"
	two := "http://test/some/path?foo"

	h1 := urlHash(one)
	h2 := urlHash(two)

	if h1 != h2 {
		t.Error("hashes mismatch")
	}
}
