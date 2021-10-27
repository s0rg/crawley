package crawler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

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

	c := New("", 1, 1, time.Millisecond*50, false)

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

	c := New("", 1, 1, time.Millisecond*50, false)

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

	c := New("", 1, 1, time.Millisecond*50, false)

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

	c := New("", 1, 1, time.Millisecond*50, false)

	if err := c.Run(ts.URL, nil); err != nil {
		t.Error("run - error")
	}
}
