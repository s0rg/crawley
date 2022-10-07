//go:build !test

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/s0rg/crawley/pkg/crawler"
	"github.com/s0rg/crawley/pkg/values"
)

const (
	appName      = "Crawley"
	appSite      = "https://github.com/s0rg/crawley"
	defaultDelay = 150 * time.Millisecond
)

var (
	GitHash   string
	GitTag    string
	BuildDate string
	defaultUA = "Mozilla/5.0 (compatible; Win64; x64) Mr." + appName + "/" + GitTag + "-" + GitHash

	cookies, headers values.Smart
	tags, ignored    values.Simple

	fDepth        = flag.Int("depth", 0, "scan depth (set -1 for unlimited)")
	fWorkers      = flag.Int("workers", runtime.NumCPU(), "number of workers")
	fBrute        = flag.Bool("brute", false, "scan html comments")
	fNoHeads      = flag.Bool("headless", false, "disable pre-flight HEAD requests")
	fScanJS       = flag.Bool("js", false, "scan js files for endpoints")
	fSkipSSL      = flag.Bool("skip-ssl", false, "skip ssl verification")
	fSilent       = flag.Bool("silent", false, "suppress info and error messages in stderr")
	fVersion      = flag.Bool("version", false, "show version")
	fDirsPolicy   = flag.String("dirs", "show", "policy for non-resource urls: show / hide / only")
	fProxyAuth    = flag.String("proxy-auth", "", "credentials for proxy: user:password")
	fRobotsPolicy = flag.String("robots", "ignore", "policy for robots.txt: ignore / crawl / respect")
	fUA           = flag.String("user-agent", defaultUA, "user-agent string")
	fDelay        = flag.Duration("delay", defaultDelay, "per-request delay (0 - disable)")
)

func version() string {
	return fmt.Sprintf("%s %s-%s build at: %s with %s site: %s",
		appName,
		GitTag,
		GitHash,
		BuildDate,
		runtime.Version(),
		appSite,
	)
}

func puts(s string) {
	_, _ = os.Stdout.WriteString(s + "\n")
}

func crawl(uri string, opts ...crawler.Option) error {
	c := crawler.New(opts...)

	log.Printf("[*] config: %s", c.DumpConfig())
	log.Printf("[*] crawling url: %s", uri)

	if err := c.Run(uri, puts); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	log.Printf("[*] complete")

	return nil
}

func loadSmart() (h, c []string, err error) {
	var wd string

	if wd, err = os.Getwd(); err != nil {
		err = fmt.Errorf("work dir: %w", err)

		return
	}

	fs := os.DirFS(wd)

	if h, err = headers.Load(fs); err != nil {
		err = fmt.Errorf("headers: %w", err)

		return
	}

	if c, err = cookies.Load(fs); err != nil {
		err = fmt.Errorf("cookies: %w", err)

		return
	}

	return h, c, nil
}

func initOptions() (rv []crawler.Option, err error) {
	robots, err := crawler.ParseRobotsPolicy(*fRobotsPolicy)
	if err != nil {
		err = fmt.Errorf("robots policy: %w", err)

		return
	}

	dirs, err := crawler.ParseDirsPolicy(*fDirsPolicy)
	if err != nil {
		err = fmt.Errorf("dirs policy: %w", err)

		return
	}

	h, c, err := loadSmart()
	if err != nil {
		err = fmt.Errorf("load: %w", err)

		return
	}

	rv = []crawler.Option{
		crawler.WithUserAgent(*fUA),
		crawler.WithDelay(*fDelay),
		crawler.WithMaxCrawlDepth(*fDepth),
		crawler.WithWorkersCount(*fWorkers),
		crawler.WithSkipSSL(*fSkipSSL),
		crawler.WithBruteMode(*fBrute),
		crawler.WithDirsPolicy(dirs),
		crawler.WithRobotsPolicy(robots),
		crawler.WithoutHeads(*fNoHeads),
		crawler.WithScanJS(*fScanJS),
		crawler.WithExtraHeaders(h),
		crawler.WithExtraCookies(c),
		crawler.WithTagsFilter(tags.Values),
		crawler.WithIgnored(ignored.Values),
		crawler.WithProxyAuth(*fProxyAuth),
	}

	return rv, nil
}

func main() {
	flag.Var(
		&headers,
		"header",
		"extra headers for request, can be used multiple times, accept files with '@'-prefix",
	)
	flag.Var(
		&cookies,
		"cookie",
		"extra cookies for request, can be used multiple times, accept files with '@'-prefix",
	)
	flag.Var(
		&tags,
		"tag",
		"tags filter, single or comma-separated tag names",
	)
	flag.Var(
		&ignored,
		"ignore",
		"patterns (in urls) to be ignored in crawl process",
	)

	flag.Parse()

	if *fVersion {
		puts(version())

		return
	}

	if flag.NArg() != 1 {
		flag.Usage()

		return
	}

	opts, err := initOptions()
	if err != nil {
		log.Fatal("[-] options:", err)
	}

	if *fSilent {
		log.SetOutput(io.Discard)
	}

	if err := crawl(flag.Arg(0), opts...); err != nil {
		// forcing back stderr in case of errors, otherwise
		// if 'silent' is on - no one will know what happened.
		log.SetOutput(os.Stderr)
		log.Fatal("[-] crawler:", err)
	}
}
