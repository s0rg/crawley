package crawler

import (
	"net/url"
	"testing"
)

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
