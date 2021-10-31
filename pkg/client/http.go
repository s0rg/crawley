package client

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	dialTimeout = 5 * time.Second
	reqTimeout  = 10 * time.Second
)

// HTTP holds pre-configured http.Client.
type HTTP struct {
	ua string
	c  *http.Client
}

// New creates and configure client for later use.
func New(ua string, conns int, skipSSL bool) (h *HTTP) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout: dialTimeout,
		}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSL,
		},
		TLSHandshakeTimeout: dialTimeout,
		MaxConnsPerHost:     conns,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 0,
	}

	client := &http.Client{
		Timeout:   reqTimeout,
		Transport: transport,
	}

	return &HTTP{
		c:  client,
		ua: ua,
	}
}

// Get sends http GET request, returns non-closed body or error.
func (h *HTTP) Get(ctx context.Context, url string) (body io.ReadCloser, err error) {
	var req *http.Request

	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
		return
	}

	if body, _, err = h.request(req); err != nil {
		return
	}

	return body, nil
}

// Head sends http HEAD request, return response headers or error.
func (h *HTTP) Head(ctx context.Context, url string) (hdrs http.Header, err error) {
	var req *http.Request

	if req, err = http.NewRequestWithContext(ctx, http.MethodHead, url, nil); err != nil {
		return
	}

	var body io.ReadCloser

	if body, hdrs, err = h.request(req); err != nil {
		return
	}

	dropContents(body)

	return hdrs, nil
}

func (h *HTTP) request(req *http.Request) (body io.ReadCloser, hdrs http.Header, err error) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml;q=0.9,*/*;q=0.5")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", h.ua)

	var resp *http.Response

	if resp, err = h.c.Do(req); err != nil {
		return
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, resp.Header, nil
	}

	dropContents(resp.Body)

	err = ErrFromResp(resp)

	return
}

func dropContents(rc io.ReadCloser) {
	_, _ = io.Copy(io.Discard, rc)
	_ = rc.Close()
}
