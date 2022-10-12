package values

import (
	"testing"
)

func TestSimpleSet(t *testing.T) {
	t.Parallel()

	var (
		l   Simple
		err error
	)

	if err = l.Set("a"); err != nil {
		t.Fatalf("set a - unexpected error: %v", err)
	}

	if len(l.Values) != 1 {
		t.Fatalf("len a - unexpected length: %d", len(l.Values))
	}

	if l.Values[0] != "a" {
		t.Fatalf("res a - unexpected value: %v", l.Values[0])
	}

	if err = l.Set("b,c"); err != nil {
		t.Fatalf("set b,c - unexpected error: %v", err)
	}

	if len(l.Values) != 3 {
		t.Fatalf("len b,c - unexpected length: %d", len(l.Values))
	}

	if l.Values[1] != "b" || l.Values[2] != "c" {
		t.Fatalf("res b,c - unexpected value: %v", l.Values[1])
	}
}

func TestSimpleString(t *testing.T) {
	t.Parallel()

	var l Simple

	if l.String() != "" {
		t.Fatal("non-empty result")
	}

	_ = l.Set("a")

	if l.String() != "a" {
		t.Fatal("expected a")
	}

	_ = l.Set("b")

	if l.String() != "a,b" {
		t.Fatal("expected a,b")
	}
}
