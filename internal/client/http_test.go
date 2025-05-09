package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	ua = "test-ua"
)

var cfg = Config{
	UserAgent: ua,
	Workers:   1,
	SkipSSL:   false,
	Timeout:   time.Second,
}

func TestHTTPGetOK(t *testing.T) {
	t.Parallel()

	tc := cfg
	tc.Headers = []string{"FOO: BAR"}
	tc.Cookies = []string{"NAME=VALUE"}

	c := New(&tc)

	const (
		body   = "test-body"
		cookie = "VALUE"
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Error("method")
		}

		if r.UserAgent() != ua {
			t.Error("agent")
		}

		if r.Header.Get("FOO") != "BAR" {
			t.Error("extra headers")
		}

		c, err := r.Cookie("NAME")
		if err != nil {
			t.Error("extra cookies - retrieve")
		}

		if c.Value != cookie {
			t.Error("extra cookies - value")
		}

		_, _ = io.WriteString(w, body)
	}))

	defer ts.Close()

	res, _, err := c.Get(t.Context(), ts.URL)
	if err != nil {
		t.Fatal("get:", err)
	}

	defer res.Close()

	buf, err := io.ReadAll(res)
	if err != nil {
		t.Fatal("read:", err)
	}

	Discard(res)

	if string(buf) != body {
		t.Error("body")
	}
}

func TestHTTPGetERR(t *testing.T) {
	t.Parallel()

	c := New(&cfg)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer ts.Close()

	if _, _, err := c.Get(t.Context(), "["); err == nil {
		t.Error("url - err is nil")
	}

	var ctx context.Context

	if _, _, err := c.Get(ctx, ts.URL); err == nil {
		t.Error("ctx - err is nil")
	}

	if _, _, err := c.Get(t.Context(), ts.URL); err == nil {
		t.Error("status - err is nil")
	}
}

func TestHTTPHeadOK(t *testing.T) {
	t.Parallel()

	c := New(&cfg)

	const (
		key = "x-some-key"
		val = "some-val"
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Error("method")
		}

		if r.UserAgent() != ua {
			t.Error("agent")
		}

		w.Header().Add(key, val)
		w.WriteHeader(http.StatusNoContent)
	}))

	defer ts.Close()

	hdr, err := c.Head(t.Context(), ts.URL)
	if err != nil {
		t.Fatal("head:", err)
	}

	if hdr.Get(key) != val {
		t.Error("bad key")
	}
}

func TestHTTPHeadERR(t *testing.T) {
	t.Parallel()

	c := New(&cfg)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer ts.Close()

	if _, err := c.Head(t.Context(), "]"); err == nil {
		t.Error("url - err is nil")
	}

	var ctx context.Context

	if _, err := c.Head(ctx, ts.URL); err == nil {
		t.Error("ctx - err is nil")
	}

	if _, err := c.Head(t.Context(), ts.URL); err == nil {
		t.Error("status - err is nil")
	}
}
