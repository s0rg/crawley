package client

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"
)

const transportTimeout = 10 * time.Second

// HTTP holds pre-configured http.Client.
type HTTP struct {
	c       *http.Client
	ua      string
	cookies []*http.Cookie
	headers []*header
}

// New creates and configure client for later use.
func New(
	ua string,
	conns int,
	skipSSL bool,
	headers, cookies []string,
	timeout time.Duration,
) (h *HTTP) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout: transportTimeout,
		}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSL,
		},
		IdleConnTimeout:     transportTimeout,
		TLSHandshakeTimeout: transportTimeout,
		MaxConnsPerHost:     conns,
		MaxIdleConns:        conns,
		MaxIdleConnsPerHost: conns,
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	return &HTTP{
		ua:      ua,
		c:       client,
		headers: prepareHeaders(headers),
		cookies: prepareCookies(cookies),
	}
}

// Get sends http GET request, returns non-closed body or error.
func (h *HTTP) Get(ctx context.Context, url string) (body io.ReadCloser, hdrs http.Header, err error) {
	var req *http.Request

	if req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		http.NoBody,
	); err != nil {
		return
	}

	if body, hdrs, err = h.request(req); err != nil {
		return
	}

	return body, hdrs, nil
}

// Head sends http HEAD request, return response headers or error.
func (h *HTTP) Head(ctx context.Context, url string) (hdrs http.Header, err error) {
	var req *http.Request

	if req, err = http.NewRequestWithContext(
		ctx,
		http.MethodHead,
		url,
		http.NoBody,
	); err != nil {
		return
	}

	var body io.ReadCloser

	if body, hdrs, err = h.request(req); err != nil {
		return
	}

	Discard(body)

	return hdrs, nil
}

// Discard read all contents from ReaderCloser, closing it afterwards.
func Discard(rc io.ReadCloser) {
	_, _ = io.Copy(io.Discard, rc)
	_ = rc.Close()
}

func (h *HTTP) request(req *http.Request) (body io.ReadCloser, hdrs http.Header, err error) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml;q=0.9,*/*;q=0.5")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", h.ua)
	req.Header.Set("Referer", req.URL.String())

	h.enrich(req)

	var resp *http.Response

	if resp, err = h.c.Do(req); err != nil {
		return
	}

	if http.StatusOK < resp.StatusCode && resp.StatusCode >= http.StatusMultipleChoices {
		err = ErrFromResp(resp)
	}

	return resp.Body, resp.Header, err
}

func (h *HTTP) enrich(req *http.Request) {
	for _, hdr := range h.headers {
		req.Header.Set(hdr.Key, hdr.Val)
	}

	for _, cck := range h.cookies {
		req.AddCookie(cck)
	}
}
