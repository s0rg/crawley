package crawler

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"mime"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html/atom"

	"github.com/s0rg/crawley/pkg/client"
	"github.com/s0rg/crawley/pkg/links"
	"github.com/s0rg/crawley/pkg/path"
	"github.com/s0rg/crawley/pkg/set"
)

const (
	nMID = 64
	nBIG = nMID * 2

	crawlTimeout = 5 * time.Second
	contentType  = "Content-Type"
	contentHTML  = "text/html"
)

type task struct {
	URI   *url.URL
	Crawl bool
	Done  bool
}

// Crawler holds crawling process config and state.
type Crawler struct {
	UserAgent string
	Delay     time.Duration
	Workers   int
	Depth     int
	SkipSSL   bool
	wg        sync.WaitGroup
	handleCh  chan string
	crawlCh   chan *url.URL
	taskCh    chan task
}

// New creates Crawler instance.
func New(ua string, workers, depth int, delay time.Duration, skipSSL bool) (c *Crawler) {
	c = &Crawler{
		UserAgent: ua,
		Workers:   workers,
		Depth:     depth,
		Delay:     delay,
		SkipSSL:   skipSSL,
	}

	return c
}

// Run starts crawling process for given base uri.
func (c *Crawler) Run(uri string, fn func(string)) (err error) {
	var base *url.URL

	if base, err = url.Parse(uri); err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	cop := (c.Workers + c.Depth + 1)
	c.handleCh = make(chan string, cop*nMID)
	c.crawlCh = make(chan *url.URL, cop*nMID)
	c.taskCh = make(chan task, cop*nBIG)

	defer c.close()

	seen := make(set.U64)

	seen.Add(urlHash(base))
	c.startCrawlers()

	go c.handler(fn)

	c.crawlCh <- base

	for w := 1; w > 0; {
		t := <-c.taskCh

		if t.Done {
			w--
		} else if seen.Add(urlHash(t.URI)) {
			if c.crawl(base, &t) {
				w++
			}

			c.handleCh <- t.URI.String()
		}
	}

	return nil
}

func (c *Crawler) crawl(b *url.URL, t *task) (yes bool) {
	if !t.Crawl {
		return
	}

	if !canCrawl(b, t.URI, c.Depth) {
		return
	}

	go func(u *url.URL) { c.crawlCh <- u }(t.URI)

	return true
}

func (c *Crawler) close() {
	c.wg.Add(c.Workers)
	close(c.crawlCh)
	c.wg.Wait() // wait for crawlers

	c.wg.Add(1)
	close(c.handleCh)
	c.wg.Wait() // wait for handler

	close(c.taskCh)
}

func (c *Crawler) startCrawlers() {
	web := client.New(c.UserAgent, c.Workers, c.SkipSSL)

	for i := 0; i < c.Workers; i++ {
		go c.crawler(web)
	}
}

func (c *Crawler) linkHandler(a atom.Atom, u *url.URL) {
	t := task{URI: u}

	switch a {
	case atom.A, atom.Iframe:
		t.Crawl = true
	}

	c.taskCh <- t
}

func (c *Crawler) crawler(web *client.HTTP) {
	defer c.wg.Done()

	for uri := range c.crawlCh {
		time.Sleep(c.Delay)

		ctx, cancel := context.WithTimeout(context.Background(), crawlTimeout)
		us := uri.String()

		var parse bool

		if hdrs, err := web.Head(ctx, us); err != nil {
			log.Printf("HEAD %s error: %v", us, err)
		} else if typ, _, perr := mime.ParseMediaType(hdrs.Get(contentType)); perr == nil {
			parse = typ == contentHTML
		}

		if parse {
			time.Sleep(c.Delay)

			if body, err := web.Get(ctx, us); err != nil {
				log.Printf("GET %s error: %v", us, err)
			} else {
				links.Extract(uri, body, c.linkHandler)
			}
		}

		cancel()

		c.taskCh <- task{Done: true}
	}
}

func (c *Crawler) handler(fn func(string)) {
	for s := range c.handleCh {
		fn(s)
	}

	c.wg.Done()
}

func canCrawl(a, b *url.URL, d int) (yes bool) {
	if a.Host != b.Host {
		return
	}

	depth, found := path.Depth(a.EscapedPath(), b.EscapedPath())
	if !found || depth > d {
		return
	}

	return true
}

func urlHash(u *url.URL) (sum uint64) {
	c := *u         // copy original
	c.RawQuery = "" // remove any query parameters

	hash := fnv.New64()
	_, _ = io.WriteString(hash, c.String())

	return hash.Sum64()
}
