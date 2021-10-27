package set

import "testing"

func Test_U64(t *testing.T) {
	t.Parallel()

	s := make(U64)

	const (
		val1 = uint64(1)
		val2 = uint64(2)
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
