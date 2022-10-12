package robots

import (
	"bytes"
	"errors"
	"net/url"
	"testing"
)

const rawRobots = `useragent: a
# some comment : with colon
disallow: /c
allow: /
user-agent: b
disallow: /d
: broken

broken
user-agent: e
sitemap: http://test.com/c
useragent: f
disallow: /g
user-agent: *
disallow:
unknown: ha-ha`

var errGeneric = errors.New("generic error")

type errReader struct {
	err error
}

func (r *errReader) Read(_ []byte) (n int, err error) {
	return 0, r.err
}

func BenchmarkFromReader(b *testing.B) {
	buf := bytes.NewBufferString(rawRobots)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = FromReader("b", buf)
	}
}

func TestFromReader(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBufferString(rawRobots)

	txt, err := FromReader("b", buf)
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	lnk := txt.Links()
	if len(lnk) != 4 {
		t.Error("links: less than expected")
	}

	smp := txt.Sitemaps()
	if len(smp) != 1 {
		t.Error("sitemaps: less than expected")
	}

	if txt.Forbidden("/a") {
		t.Error("cannot access /a")
	}

	if !txt.Forbidden("/d") {
		t.Error("can access /d")
	}
}

func TestNewAllowAll(t *testing.T) {
	t.Parallel()

	txt := AllowALL()

	if txt.Forbidden("/a") {
		t.Error("cannot access /a")
	}
}

func TestNewDenyAll(t *testing.T) {
	t.Parallel()

	txt := DenyALL()

	if !txt.Forbidden("/a") {
		t.Error("can access /a")
	}
}

func TestErrReader(t *testing.T) {
	t.Parallel()

	errIO := errReader{err: errGeneric}

	_, err := FromReader("", &errIO)
	if err == nil {
		t.Fatal("no error")
	}

	if !errors.Is(err, errGeneric) {
		t.Error("bad error")
	}
}

func TestURL(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"http://example.com/",
		"http://example.com/some/path",
		"http://example.com/some/path?with=query",
	}

	const wantURL = "http://example.com/robots.txt"

	for _, c := range testCases {
		u, err := url.Parse(c)
		if err != nil {
			t.Fatalf("cannot parse: %s: %v", c, err)
		}

		if URL(u) != wantURL {
			t.Errorf("fail for: %s", c)
		}
	}
}
