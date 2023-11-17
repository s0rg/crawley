package crawler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/s0rg/set"
	"golang.org/x/net/html/atom"

	"github.com/s0rg/crawley/pkg/client"
	"github.com/s0rg/crawley/pkg/links"
	"github.com/s0rg/crawley/pkg/robots"
)

type crawlClient interface {
	Get(context.Context, string) (io.ReadCloser, http.Header, error)
	Head(context.Context, string) (http.Header, error)
}

const (
	chMult    = 256
	chTimeout = 100 * time.Millisecond
)

type taskFlag byte

const (
	// TaskDefault marks result for printing only.
	TaskDefault taskFlag = iota
	// TaskCrawl marks result as to-be-crawled.
	TaskCrawl
	// TaskDone marks result as final - crawling ends here.
	TaskDone
)

type crawlResult struct {
	URI  string
	Hash uint64
	Flag taskFlag
}

// Crawler holds crawling process config and state.
type Crawler struct {
	cfg      *config
	handleCh chan string
	crawlCh  chan *url.URL
	resultCh chan crawlResult
	robots   *robots.TXT
	filter   links.TokenFilter
	wg       sync.WaitGroup
}

// New creates Crawler instance.
func New(opts ...Option) (c *Crawler) {
	cfg := &config{}

	for _, o := range opts {
		o(cfg)
	}

	cfg.validate()

	return &Crawler{
		cfg:    cfg,
		robots: robots.AllowALL(),
		filter: prepareFilter(cfg.AlowedTags),
	}
}

// Run starts crawling process for given base uri.
func (c *Crawler) Run(uri string, urlcb func(string)) (err error) {
	var base *url.URL

	if base, err = url.Parse(uri); err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	workers := c.cfg.Client.Workers

	n := (workers + 1)
	c.handleCh = make(chan string, n*chMult)
	c.crawlCh = make(chan *url.URL, n*chMult)
	c.resultCh = make(chan crawlResult, n*chMult)

	defer c.close()

	seen := make(set.Unordered[uint64])
	seen.Add(urlhash(uri))

	web := client.New(&c.cfg.Client)
	c.initRobots(base, web)

	for i := 0; i < workers; i++ {
		go c.worker(web)
	}

	c.wg.Add(workers)

	go func() {
		for s := range c.handleCh {
			urlcb(s)
		}

		c.wg.Done()
	}()

	c.crawlCh <- base

	var t crawlResult

	for w := 1; w > 0; {
		t = <-c.resultCh

		switch {
		case t.Flag == TaskDone:
			w--
		case seen.Add(t.Hash):
			if t.Flag == TaskCrawl && c.tryEnqueue(base, &t) {
				w++
			}

			c.tryHandle(t.URI)
		}
	}

	return nil
}

// DumpConfig returns internal config representation.
func (c *Crawler) DumpConfig() string {
	return c.cfg.String()
}

func (c *Crawler) tryHandle(u string) {
	show := true

	idx := strings.LastIndexByte(u, '/')
	if idx == -1 {
		return
	}

	switch c.cfg.Dirs {
	case DirsHide:
		show = isResorce(u[idx:])
	case DirsOnly:
		show = !isResorce(u[idx:])
	}

	if !show {
		return
	}

	t := time.NewTimer(chTimeout)
	defer t.Stop()

	select {
	case c.handleCh <- u:
	case <-t.C:
	}
}

func (c *Crawler) tryEnqueue(base *url.URL, r *crawlResult) (yes bool) {
	u, err := url.Parse(r.URI)
	if err != nil {
		return
	}

	if !canCrawl(base, u, c.cfg.Depth) ||
		c.robots.Forbidden(u.Path) ||
		(c.cfg.Dirs == DirsOnly && isResorce(u.Path)) {
		return
	}

	t := time.NewTimer(chTimeout)
	defer t.Stop()

	select {
	case c.crawlCh <- u:
		return true
	case <-t.C:
	}

	return false
}

func (c *Crawler) close() {
	close(c.crawlCh)
	c.wg.Wait() // wait for crawlers

	c.wg.Add(1) // for handler`s Done()
	close(c.handleCh)
	c.wg.Wait() // wait for handler

	close(c.resultCh)
}

func (c *Crawler) initRobots(host *url.URL, web crawlClient) {
	if c.cfg.Robots == RobotsIgnore {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Client.Timeout)
	defer cancel()

	body, _, err := web.Get(ctx, robots.URL(host))
	if err != nil {
		var herr client.HTTPError

		if !errors.As(err, &herr) {
			log.Println("[-] GET /robots.txt:", err)

			return
		}

		if herr.Code() == http.StatusInternalServerError {
			c.robots = robots.DenyALL()
		}

		return
	}

	defer body.Close()

	rbt, err := robots.FromReader(c.cfg.Client.UserAgent, body)
	if err != nil {
		log.Println("[-] parse robots.txt:", err)

		return
	}

	c.robots = rbt

	c.crawlRobots(host)
}

func (c *Crawler) crawlRobots(host *url.URL) {
	base := *host
	base.Fragment = ""
	base.RawQuery = ""

	for _, u := range c.robots.Links() {
		t := base
		t.Path = u

		c.linkHandler(atom.A, t.String())
	}

	for _, u := range c.robots.Sitemaps() {
		if _, e := url.Parse(u); e == nil {
			c.crawlHandler(u)
		}
	}
}

func (c *Crawler) isIgnored(v string) (yes bool) {
	if len(c.cfg.Ignored) == 0 {
		return
	}

	return slices.ContainsFunc(c.cfg.Ignored, func(s string) bool {
		return strings.Contains(v, s)
	})
}

func (c *Crawler) linkHandler(a atom.Atom, s string) {
	r := crawlResult{
		URI:  s,
		Hash: urlhash(s),
	}

	fetch := (a == atom.A || a == atom.Iframe) ||
		(c.cfg.ScanJS && a == atom.Script)

	if fetch && !c.isIgnored(s) {
		r.Flag = TaskCrawl
	}

	t := time.NewTimer(chTimeout)
	defer t.Stop()

	select {
	case c.resultCh <- r:
	case <-t.C:
	}
}

func (c *Crawler) staticHandler(s string) {
	c.linkHandler(atom.Link, s)
}

func (c *Crawler) crawlHandler(s string) {
	c.linkHandler(atom.A, s)
}

func (c *Crawler) process(
	ctx context.Context,
	web crawlClient,
	base *url.URL,
	uri string,
) {
	body, hdrs, err := web.Get(ctx, uri)
	if err != nil {
		var herr client.HTTPError

		// ignore any http errors, just parse body (if any)
		if !errors.As(err, &herr) {
			log.Printf("[-] GET %s: %v", uri, err)

			return
		}
	}

	content := hdrs.Get(contentType)

	switch {
	case isHTML(content):
		links.ExtractHTML(body, base, links.HTMLParams{
			Brute:        c.cfg.Brute,
			ScanJS:       c.cfg.ScanJS,
			ScanCSS:      c.cfg.ScanCSS,
			Filter:       c.filter,
			HandleHTML:   c.linkHandler,
			HandleStatic: c.staticHandler,
		})
	case isSitemap(uri):
		links.ExtractSitemap(body, base, c.crawlHandler)
	case c.cfg.ScanJS && isJS(content, uri):
		links.ExtractJS(body, c.staticHandler)
	case c.cfg.ScanCSS && false:
		links.ExtractCSS(body, c.staticHandler)
	}

	client.Discard(body)
}

func (c *Crawler) worker(web crawlClient) {
	defer c.wg.Done()

	for uri := range c.crawlCh {
		if c.cfg.Delay > 0 {
			time.Sleep(c.cfg.Delay)
		}

		ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Client.Timeout)
		us := uri.String()

		var canProcess bool

		if c.cfg.NoHEAD {
			canProcess = canParse(uri.Path)
		} else {
			if hdrs, err := web.Head(ctx, us); err != nil {
				log.Printf("[-] HEAD %s: %v", us, err)
			} else {
				ct := hdrs.Get(contentType)

				canProcess = isHTML(ct) || isSitemap(us) || (c.cfg.ScanJS && isJS(ct, us))
			}
		}

		if canProcess {
			c.process(ctx, web, uri, us)
		}

		cancel()

		c.resultCh <- crawlResult{Flag: TaskDone}
	}
}
