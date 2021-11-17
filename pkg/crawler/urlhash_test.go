package crawler

import (
	"net/url"
	"testing"
)

func TestUrlHash(t *testing.T) {
	t.Parallel()

	one, _ := url.Parse("http://test/some/path?foo")
	two, _ := url.Parse("http://test/some/path?foo#bar")

	h1 := urlHash(one)
	h2 := urlHash(two)

	if one.RawQuery == "" {
		t.Error("one: modified")
	}

	if two.RawQuery == "" {
		t.Error("two: modified")
	}

	if h1 != h2 {
		t.Error("hashes mismatch")
	}
}
