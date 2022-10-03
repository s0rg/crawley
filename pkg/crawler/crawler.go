package crawler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html/atom"

	"github.com/s0rg/crawley/pkg/client"
	"github.com/s0rg/crawley/pkg/links"
	"github.com/s0rg/crawley/pkg/robots"
	"github.com/s0rg/crawley/pkg/set"
)

type crawlClient interface {
	Get(context.Context, string) (io.ReadCloser, http.Header, error)
	Head(context.Context, string) (http.Header, error)
}

const (
	nMID = 64
	nBIG = nMID * 2

	crawlTimeout  = 5 * time.Second
	robotsTimeout = 3 * time.Second
)

type taskFlag byte

const (
	// TaskDefault marks result for printing only.
	TaskDefault taskFlag = iota
	// TaskCrawl marks result as can-be-crawled.
	TaskCrawl
	// TaskDone marks result as final - crawling end here.
	TaskDone
)

type crawlResult struct {
	URI  string
	Flag taskFlag
}

// Crawler holds crawling process config and state.
type Crawler struct {
	cfg      *config
	wg       sync.WaitGroup
	handleCh chan string
	crawlCh  chan *url.URL
	resultCh chan crawlResult
	robots   *robots.TXT
	filter   links.TokenFilter
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
		filter: prepareFilter(cfg.AlowedTags),
	}
}

// Run starts crawling process for given base uri.
func (c *Crawler) Run(uri string, fn func(string)) (err error) {
	var base *url.URL

	if base, err = url.Parse(uri); err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	n := (c.cfg.Workers + c.cfg.Depth + 1)
	c.handleCh = make(chan string, n*nMID)
	c.crawlCh = make(chan *url.URL, n*nMID)
	c.resultCh = make(chan crawlResult, n*nBIG)

	defer c.close()

	seen := make(set.URI)
	seen.Add(uri)

	web := client.New(
		c.cfg.UserAgent,
		c.cfg.Workers,
		c.cfg.SkipSSL,
		c.cfg.Headers,
		c.cfg.Cookies,
	)
	c.initRobots(base, web)

	for i := 0; i < c.cfg.Workers; i++ {
		go c.crawler(web)
	}

	c.wg.Add(c.cfg.Workers)

	go c.handler(fn)

	c.crawlCh <- base

	var t crawlResult

	for w := 1; w > 0; {
		t = <-c.resultCh

		switch {
		case t.Flag == TaskDone:
			w--
		case seen.Add(t.URI):
			if t.Flag == TaskCrawl && c.crawl(base, &t) {
				w++
			}

			c.emit(t.URI)
		}
	}

	return nil
}

// DumpConfig returns internal config representation.
func (c *Crawler) DumpConfig() string {
	return c.cfg.String()
}

func (c *Crawler) emit(u string) {
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

	c.handleCh <- u
}

func (c *Crawler) crawl(base *url.URL, t *crawlResult) (yes bool) {
	u, err := url.Parse(t.URI)
	if err != nil {
		return
	}

	switch {
	case !canCrawl(base, u, c.cfg.Depth), c.robots.Forbidden(u.Path), c.cfg.Dirs == DirsOnly && isResorce(u.Path):
		return
	default:
		go func(r *url.URL) { c.crawlCh <- r }(u)
	}

	return true
}

func (c *Crawler) close() {
	close(c.crawlCh)
	c.wg.Wait() // wait for crawlers

	c.wg.Add(1)
	close(c.handleCh)
	c.wg.Wait() // wait for handler

	close(c.resultCh)
}

func (c *Crawler) initRobots(host *url.URL, web crawlClient) {
	c.robots = robots.AllowALL()

	if c.cfg.Robots == RobotsIgnore {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), robotsTimeout)
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

	rbt, err := robots.FromReader(c.cfg.UserAgent, body)
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
			c.linkHandler(atom.A, u)
		}
	}
}

func (c *Crawler) sitemapHandler(s string) {
	c.linkHandler(atom.A, s)
}

func (c *Crawler) linkHandler(a atom.Atom, s string) {
	t := crawlResult{URI: s}

	switch a {
	case atom.A, atom.Iframe:
		t.Flag = TaskCrawl
	}

	c.resultCh <- t
}

func (c *Crawler) fetch(
	ctx context.Context,
	web crawlClient,
	base *url.URL,
	uri string,
) {
	body, hdrs, err := web.Get(ctx, uri)
	if err != nil {
		var herr client.HTTPError

		if !errors.As(err, &herr) {
			log.Printf("[-] GET %s: %v", uri, err)

			return
		}
	}

	switch {
	case isHTML(hdrs.Get(contentType)):
		links.ExtractHTML(base, body, c.cfg.Brute, c.filter, c.linkHandler)
	case isSitemap(uri):
		links.ExtractSitemap(base, body, c.sitemapHandler)
	}

	client.Discard(body)
}

func (c *Crawler) crawler(web crawlClient) {
	defer c.wg.Done()

	for uri := range c.crawlCh {
		if c.cfg.Delay > 0 {
			time.Sleep(c.cfg.Delay)
		}

		ctx, cancel := context.WithTimeout(context.Background(), crawlTimeout)
		us := uri.String()

		var parse bool

		if c.cfg.NoHEAD {
			parse = canParse(uri.Path)
		} else {
			if hdrs, err := web.Head(ctx, us); err != nil {
				log.Printf("[-] HEAD %s: %v", us, err)
			} else {
				parse = isHTML(hdrs.Get(contentType)) || isSitemap(us)
			}
		}

		if parse {
			c.fetch(ctx, web, uri, us)
		}

		cancel()

		c.resultCh <- crawlResult{Flag: TaskDone}
	}
}

func (c *Crawler) handler(fn func(string)) {
	for s := range c.handleCh {
		fn(s)
	}

	c.wg.Done()
}
