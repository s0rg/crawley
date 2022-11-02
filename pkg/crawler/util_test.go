package crawler

import (
	"net/url"
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
		{Type: "", Name: nameJS + "?v=1", Want: true},
	}

	for i, tc := range cases {
		if isJS(tc.Type, tc.Name) != tc.Want {
			t.Fatalf("case[%d] fail", i)
		}
	}
}

func TestRelativeDepth(t *testing.T) {
	type args struct {
		base string
		sub  string
	}

	tests := []struct {
		name      string
		args      args
		wantDepht int
		wantFound bool
	}{
		{"a-ok", args{base: "/", sub: "/a"}, 1, true},
		{"a-bad", args{base: "/a", sub: "/b"}, 0, false},
		{"c-bad", args{base: "/a/b", sub: "/c"}, 0, false},
		{"b-ok", args{base: "/a", sub: "/a/b"}, 1, true},
		{"c-ok", args{base: "/a", sub: "/a/b/c"}, 2, true},
		{"d-bad", args{base: "/a/b/c", sub: "/d/b/c/a"}, 0, false},
	}

	t.Parallel()

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotDepht, gotFound := relativeDepth(tc.args.base, tc.args.sub)
			if gotDepht != tc.wantDepht {
				t.Errorf("pathDepth() gotDepht = %v, want %v", gotDepht, tc.wantDepht)
			}
			if gotFound != tc.wantFound {
				t.Errorf("pathDepth() gotFound = %v, want %v", gotFound, tc.wantFound)
			}
		})
	}
}

func BenchmarkDepth(b *testing.B) {
	const (
		x = "/some/rather/long/path"
		y = "/some/rather/long/path/but/longer"
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = relativeDepth(x, y)
	}
}

func TestCanCrawl(t *testing.T) {
	t.Parallel()

	type args struct {
		b *url.URL
		u *url.URL
		d int
	}

	base, _ := url.Parse("http://test/some/path")
	badh, _ := url.Parse("http://other/path")
	url0, _ := url.Parse("http://test/some")
	url1, _ := url.Parse("http://test/some/path/even")
	url2, _ := url.Parse("http://test/some/path/even/more")
	url3, _ := url.Parse("http://test")

	tests := []struct {
		name    string
		args    args
		wantYes bool
	}{
		{"url0-1", args{b: base, u: url0, d: 1}, false},
		{"url1-0", args{b: base, u: url1, d: 0}, false},
		{"url1-1", args{b: base, u: url1, d: 1}, true},
		{"url2-0", args{b: base, u: url2, d: 0}, false},
		{"url2-1", args{b: base, u: url2, d: 1}, false},
		{"url2-2", args{b: base, u: url2, d: 2}, true},
		{"url2-3", args{b: base, u: url2, d: 3}, true},
		{"badh-1", args{b: base, u: badh, d: 1}, false},
		{"url2-0-1", args{b: base, u: url0, d: -1}, false},
		{"url2-1-1", args{b: base, u: url1, d: -1}, true},
		{"url2-2-1", args{b: base, u: url2, d: -1}, true},
		{"url3-3", args{b: base, u: url3, d: 0}, false},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if gotYes := canCrawl(tc.args.b, tc.args.u, tc.args.d); gotYes != tc.wantYes {
				t.Errorf("canCrawl() = %v, want %v", gotYes, tc.wantYes)
			}
		})
	}
}

func TestIsResorce(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Val  string
		Want bool
	}

	cases := []testCase{
		{"/some/path", false},
		{"/some/other/path/", false},
		{"/path/to/some/resource.zip", true},
	}

	for _, tc := range cases {
		if got := isResorce(tc.Val); got != tc.Want {
			t.Errorf("failed for: '%s' got: %t", tc.Val, got)
		}
	}
}

func TestCanParse(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Val  string
		Want bool
	}

	cases := []testCase{
		{"/some/path", true},
		{"/some/other/path/", true},
		{"/some/resource.html", true},
		{"/path/to/some/resource.zip", false},
	}

	for _, tc := range cases {
		if got := canParse(tc.Val); got != tc.Want {
			t.Errorf("failed for: '%s' got: %t", tc.Val, got)
		}
	}
}

func TestIsSitemap(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Val  string
		Want bool
	}

	cases := []testCase{
		{"/some/path", false},
		{"/some/other/path/sitemap.xml", true},
		{"/some/resource.html", false},
		{"/path/to/some/sitemap-index.xml", true},
	}

	for _, tc := range cases {
		if got := isSitemap(tc.Val); got != tc.Want {
			t.Errorf("failed for: '%s' got: %t", tc.Val, got)
		}
	}
}

func TestUrlHash(t *testing.T) {
	t.Parallel()

	const val = "http://test/some/path?foo"

	h1, h2 := urlhash(val), urlhash(val)

	if h1 != h2 {
		t.Error("hashes mismatch")
	}
}

func TestProxyAuthHeader(t *testing.T) {
	t.Parallel()

	const (
		got  = "user:pass"
		want = "Proxy-Authorization: Basic dXNlcjpwYXNz"
	)

	if rv := proxyAuthHeader(got); rv != want {
		t.Errorf("invalid header want: '%s' got: '%s'", want, rv)
	}
}
