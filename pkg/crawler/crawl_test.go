package crawler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/s0rg/crawley/pkg/set"
)

const robotsEP = "/robots.txt"

func Test_canCrawl(t *testing.T) {
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

func Test_urlHash(t *testing.T) {
	t.Parallel()

	one, _ := url.Parse("http://test/some/path?foo=bar")
	two, _ := url.Parse("http://test/some/path?other")

	h1 := urlHash(one)
	h2 := urlHash(two)

	if one.RawQuery == "" {
		t.Error("one: modified")
	}

	if two.RawQuery == "" {
		t.Error("two: modified")
	}

	if h1 != h2 {
		t.Error("hashes mismatch")
	}
}

func Test_Crawler(t *testing.T) {
	t.Parallel()

	var (
		gotHEAD bool
		gotGET  bool
	)

	const body = `
<html>
<a href="result">here</a>
<img src="http://other.host/image.bmp"/>
<iframe src="some/deep/path"/>
</html>`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			gotHEAD = true
			w.Header().Add(contentType, contentHTML)

		case http.MethodGet:
			if gotHEAD {
				gotGET = true
			}
			_, _ = io.WriteString(w, body)
		}
	}))

	defer ts.Close()

	results := make([]string, 0, 3)

	handler := func(s string) {
		results = append(results, s)
	}

	c := New("", 1, 1, time.Millisecond*50, false, RobotsIgnore)

	if err := c.Run(ts.URL, handler); err != nil {
		t.Errorf("run: %v", err)
	}

	if !gotHEAD {
		t.Error("no head")
	}

	if !gotGET {
		t.Error("no get")
	}

	if len(results) != 3 {
		t.Log(results)
		t.Error("results: less than expected")
	}

	if results[1] != "http://other.host/image.bmp" {
		t.Error("results: no image")
	}

	if !strings.HasSuffix(results[0], "/result") {
		t.Error("results: bad item at 0")
	}

	if !strings.HasSuffix(results[2], "/deep/path") {
		t.Error("results: bad item at 2")
	}
}

func Test_CrawlerBadLink(t *testing.T) {
	t.Parallel()

	c := New("", 1, 1, time.Millisecond*50, false, RobotsIgnore)

	if err := c.Run("%", nil); err == nil {
		t.Error("run - no error")
	}
}

func Test_CrawlerBadHead(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer ts.Close()

	c := New("", 1, 1, time.Millisecond*50, false, RobotsIgnore)

	if err := c.Run(ts.URL, nil); err != nil {
		t.Error("run - error")
	}
}

func Test_CrawlerBadGet(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			w.Header().Add(contentType, contentHTML)

		case http.MethodGet:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	defer ts.Close()

	c := New("", 1, 1, time.Millisecond*50, false, RobotsIgnore)

	if err := c.Run(ts.URL, nil); err != nil {
		t.Error("run - error")
	}
}

func Test_CrawlerRobots(t *testing.T) {
	t.Parallel()

	const (
		body  = `<html><a href="/a">a</a><a href="/b">b</a><a href="/c">c</a></html>`
		bodyA = `<html><a href="http://a">a</a></html>`
		bodyB = `<html><a href="http://b">b</a></html>`
		bodyC = `<html><a href="http://c">c</a></html>`
		robot = `useragent: a
disallow: /a
disallow: /c
user-agent: b
disallow: /b
sitemap: http://other.host/sitemap.xml`

		resSitemap = "http://other.host/sitemap.xml"
		resHostA   = "http://a/"
		resHostB   = "http://b/"
		resHostC   = "http://c/"
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case robotsEP:
			_, _ = io.WriteString(w, robot)

		case "/a":
			switch r.Method {
			case http.MethodHead:
				w.Header().Add(contentType, contentHTML)

			case http.MethodGet:
				_, _ = io.WriteString(w, bodyA)
			}

		case "/b":
			switch r.Method {
			case http.MethodHead:
				w.Header().Add(contentType, contentHTML)

			case http.MethodGet:
				_, _ = io.WriteString(w, bodyB)
			}

		case "/c":
			switch r.Method {
			case http.MethodHead:
				w.Header().Add(contentType, contentHTML)

			case http.MethodGet:
				_, _ = io.WriteString(w, bodyC)
			}

		default:
			_, _ = io.WriteString(w, body)
		}
	}))

	defer ts.Close()

	// case A

	resA := make(set.String)

	handlerA := func(s string) {
		resA.Add(s)
	}

	cA := New("a", 1, 1, time.Millisecond*50, false, RobotsRespect)

	if err := cA.Run(ts.URL, handlerA); err != nil {
		t.Error("run A - error:", err)
	}

	if len(resA) != 5 {
		t.Fatal("unexpected len for A")
	}

	if !resA.Has(resSitemap) {
		t.Error("no sitemap for A")
	}

	if !resA.Has(resHostB) {
		t.Error("no hostB for A")
	}

	if resA.Has(resHostA) {
		t.Error("hostA for A")
	}

	if resA.Has(resHostC) {
		t.Error("hostC for A")
	}

	// case B

	resB := make(set.String)

	handlerB := func(s string) {
		resB.Add(s)
	}

	cB := New("b", 1, 1, time.Millisecond*50, false, RobotsRespect)

	if err := cB.Run(ts.URL, handlerB); err != nil {
		t.Error("run B - error:", err)
	}

	if len(resB) != 6 {
		t.Fatal("unexpected len for B")
	}

	if !resB.Has(resSitemap) {
		t.Error("no sitemap for B")
	}

	if resB.Has(resHostB) {
		t.Error("hostB for B")
	}

	if !resB.Has(resHostA) {
		t.Error("no hostA for B")
	}

	if !resB.Has(resHostC) {
		t.Error("no hostC for B")
	}
}

func Test_CrawlerRobots500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case robotsEP:
			w.WriteHeader(http.StatusInternalServerError)

		default:
			_, _ = io.WriteString(w, "")
		}
	}))

	defer ts.Close()

	res := []string{}

	handler := func(s string) {
		res = append(res, s)
	}

	c := New("", 1, 1, time.Millisecond*50, false, RobotsRespect)

	if err := c.Run(ts.URL, handler); err != nil {
		t.Error("run error:", err)
	}

	if len(res) != 0 {
		t.Error("unexpected len")
	}

	if !c.robots.Forbidden("/some") {
		t.Error("not forbidden")
	}
}

func Test_CrawlerRobots400(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case robotsEP:
			w.WriteHeader(http.StatusForbidden)

		default:
			_, _ = io.WriteString(w, "")
		}
	}))

	defer ts.Close()

	res := []string{}

	handler := func(s string) {
		res = append(res, s)
	}

	c := New("", 1, 1, time.Millisecond*50, false, RobotsRespect)

	if err := c.Run(ts.URL, handler); err != nil {
		t.Error("run error:", err)
	}

	if len(res) != 0 {
		t.Error("unexpected len")
	}

	if c.robots.Forbidden("/some") {
		t.Error("forbidden")
	}
}

type testClient struct {
	err    error
	bodyIO io.ReadCloser
}

func (tc *testClient) Get(_ context.Context, _ string) (body io.ReadCloser, err error) {
	return tc.bodyIO, tc.err
}

func (tc *testClient) Head(_ context.Context, _ string) (h http.Header, err error) {
	return
}

type errReader struct {
	err error
}

func (er *errReader) Read(_ []byte) (n int, err error) {
	return 0, er.err
}

func Test_CrawlerRobotsRequestErr(t *testing.T) {
	t.Parallel()

	var (
		base, _ = url.Parse("http://test/")
		genErr  = errors.New("generic error")
		tc      = testClient{err: genErr}
		c       = New("", 1, 1, time.Millisecond*50, false, RobotsRespect)
	)

	c.initRobots(base, &tc)

	if c.robots.Forbidden("/some") {
		t.Error("forbidden")
	}
}

func Test_CrawlerRobotsBodytErr(t *testing.T) {
	t.Parallel()

	var (
		base, _ = url.Parse("http://test/")
		genErr  = errors.New("generic error")
		tc      = testClient{err: nil, bodyIO: io.NopCloser(&errReader{err: genErr})}
		c       = New("", 1, 1, time.Millisecond*50, false, RobotsRespect)
	)

	c.initRobots(base, &tc)

	if c.robots.Forbidden("/some") {
		t.Error("forbidden")
	}
}
