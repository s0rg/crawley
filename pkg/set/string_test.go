package set

import (
	"testing"
)

func Test_String(t *testing.T) {
	t.Parallel()

	s := make(String)

	const (
		val1 = "a"
		val2 = "b"
	)

	s.Add(val1)
	s.Add(val2)
	s.Add(val1)
	s.Add(val2)

	if len(s.List()) != 2 {
		t.Error("unexpected length")
	}
}
