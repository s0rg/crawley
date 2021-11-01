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

type errReader struct {
	err error
}

func (er *errReader) Read(_ []byte) (n int, err error) {
	return 0, er.err
}

func Benchmark_FromReader(b *testing.B) {
	buf := bytes.NewBufferString(rawRobots)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = FromReader("b", buf)
	}
}

func Test_FromReader(t *testing.T) {
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

func Test_NewAllowAll(t *testing.T) {
	t.Parallel()

	txt := AllowALL()

	if txt.Forbidden("/a") {
		t.Error("cannot access /a")
	}
}

func Test_NewDenyAll(t *testing.T) {
	t.Parallel()

	txt := DenyALL()

	if !txt.Forbidden("/a") {
		t.Error("can access /a")
	}
}

func Test_ErrReader(t *testing.T) {
	var (
		genErr = errors.New("generic error")
		errIO  = errReader{err: genErr}
	)

	t.Parallel()

	_, err := FromReader("", &errIO)
	if err == nil {
		t.Fatal("no error")
	}

	if !errors.Is(err, genErr) {
		t.Error("bad error")
	}
}

func Test_URL(t *testing.T) {
	testCases := []string{
		"http://example.com/",
		"http://example.com/some/path",
		"http://example.com/some/path?with=query",
	}

	const wantURL = "http://example.com/robots.txt"

	t.Parallel()

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
