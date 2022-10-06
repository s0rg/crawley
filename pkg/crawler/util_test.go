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

func TestIsJS(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Type string
		Name string
		Want bool
	}

	const (
		nameHTML = "test.html"
		nameJS   = "test.js"
	)

	cases := []testCase{
		{Type: contentHTML, Name: nameHTML, Want: false},
		{Type: contentJS, Name: nameJS, Want: true},
		{Type: contentJS, Name: nameHTML, Want: true},
		{Type: contentHTML, Name: nameJS, Want: true},
		{Type: "", Name: nameHTML, Want: false},
		{Type: "", Name: nameJS, Want: true},
	}

	for i, tc := range cases {
		if isJS(tc.Type, tc.Name) != tc.Want {
			t.Fatalf("case[%d] fail", i)
		}
	}
}
