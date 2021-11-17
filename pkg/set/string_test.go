package set

import (
	"testing"
)

func TestString(t *testing.T) {
	t.Parallel()

	s := make(String)

	const (
		val1 = "a"
		val2 = "b"
		val3 = "c"
	)

	s.Add(val1)
	s.Add(val2)
	s.Load(val1, val2)

	if len(s.List()) != 2 {
		t.Error("unexpected length")
	}

	if !s.Has(val1) {
		t.Error("no val1")
	}

	if s.Has(val3) {
		t.Error("has val3")
	}
}
