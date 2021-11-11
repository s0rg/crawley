package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	ua = "test-ua"
)

func TestHTTPGetOK(t *testing.T) {
	t.Parallel()

	c := New(ua, 1, false)

	const body = "test-body"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Error("method")
		}

		if r.UserAgent() != ua {
			t.Error("agent")
		}

		_, _ = io.WriteString(w, body)
	}))

	defer ts.Close()

	res, err := c.Get(context.Background(), ts.URL)
	if err != nil {
		t.Fatal("get:", err)
	}

	defer res.Close()

	buf, err := io.ReadAll(res)
	if err != nil {
		t.Fatal("read:", err)
	}

	if string(buf) != body {
		t.Error("body")
	}
}

func TestHTTPGetERR(t *testing.T) {
	t.Parallel()

	c := New("", 1, false)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer ts.Close()

	if _, err := c.Get(context.Background(), "["); err == nil {
		t.Error("url - err is nil")
	}

	var ctx context.Context

	if _, err := c.Get(ctx, ts.URL); err == nil {
		t.Error("ctx - err is nil")
	}

	if _, err := c.Get(context.Background(), ts.URL); err == nil {
		t.Error("status - err is nil")
	}
}

func TestHTTPHeadOK(t *testing.T) {
	t.Parallel()

	c := New(ua, 1, false)

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

	hdr, err := c.Head(context.Background(), ts.URL)
	if err != nil {
		t.Fatal("head:", err)
	}

	if hdr.Get(key) != val {
		t.Error("bad key")
	}
}

func TestHTTPHeadERR(t *testing.T) {
	t.Parallel()

	c := New("", 1, false)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer ts.Close()

	if _, err := c.Head(context.Background(), "]"); err == nil {
		t.Error("url - err is nil")
	}

	var ctx context.Context

	if _, err := c.Head(ctx, ts.URL); err == nil {
		t.Error("ctx - err is nil")
	}

	if _, err := c.Head(context.Background(), ts.URL); err == nil {
		t.Error("status - err is nil")
	}
}
