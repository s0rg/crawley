package crawler

import (
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestPrepareFilter(t *testing.T) {
	t.Parallel()

	f := prepareFilter([]string{"a", "form", "foo", "div"})
	d := []html.Token{
		{DataAtom: atom.A},
		{DataAtom: atom.Form},
		{DataAtom: atom.Div},
	}

	for _, k := range d {
		if !f(k) {
			t.Fatalf("not allowed: %v", k)
		}
	}

	v := html.Token{DataAtom: atom.Video}

	if f(v) {
		t.Fatalf("allowed: %v", v)
	}
}
