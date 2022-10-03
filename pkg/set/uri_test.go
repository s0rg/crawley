package set

import (
	"testing"
)

func TestURI(t *testing.T) {
	t.Parallel()

	s := make(URI)

	const (
		val1 = "http://test/1"
		val2 = "http://test/2"
	)

	if !s.Add(val1) {
		t.Errorf("add val1 - step 1 failure")
	}

	if !s.Add(val2) {
		t.Errorf("add val2 - step 1 failure")
	}

	if s.Add(val1) {
		t.Errorf("add val1 - step 2 failure")
	}

	if s.Add(val2) {
		t.Errorf("add val2 - step 2 failure")
	}
}

func TestHash(t *testing.T) {
	t.Parallel()

	const val = "http://test/some/path?foo"

	h1, h2 := hash(val), hash(val)

	if h1 != h2 {
		t.Error("hashes mismatch")
	}
}
