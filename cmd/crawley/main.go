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

	extCookies, extHeaders values.Smart
	tags                   values.Simple

	fVersion      = flag.Bool("version", false, "show version")
	fBrute        = flag.Bool("brute", false, "scan html comments")
	fSkipSSL      = flag.Bool("skip-ssl", false, "skip ssl verification")
	fSilent       = flag.Bool("silent", false, "suppress info and error messages in stderr")
	fNoHeads      = flag.Bool("headless", false, "disable pre-flight HEAD requests")
	fDepth        = flag.Int("depth", 0, "scan depth (-1 - unlimited)")
	fWorkers      = flag.Int("workers", runtime.NumCPU(), "number of workers")
	fDelay        = flag.Duration("delay", defaultDelay, "per-request delay (0 - disable)")
	fUA           = flag.String("user-agent", defaultUA, "user-agent string")
	fRobotsPolicy = flag.String("robots", "ignore", "policy for robots.txt: ignore / crawl / respect")
	fDirsPolicy   = flag.String("dirs", "show", "policy for non-resource urls: show / hide / only")
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

func initValues() (
	headers, cookies []string,
	err error,
) {
	var wd string

	if wd, err = os.Getwd(); err != nil {
		err = fmt.Errorf("work dir: %w", err)

		return
	}

	fs := os.DirFS(wd)

	if headers, err = extHeaders.Load(fs); err != nil {
		err = fmt.Errorf("headers: %w", err)

		return
	}

	if cookies, err = extCookies.Load(fs); err != nil {
		err = fmt.Errorf("cookies: %w", err)

		return
	}

	return headers, cookies, nil
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

	headers, cookies, err := initValues()
	if err != nil {
		err = fmt.Errorf("values: %w", err)

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
		crawler.WithExtraHeaders(headers),
		crawler.WithExtraCookies(cookies),
		crawler.WithTagsFilter(tags.Values),
	}

	return rv, nil
}

func main() {
	flag.Var(
		&extHeaders,
		"header",
		"extra headers for request, can be used multiple times, accept files with '@'-prefix",
	)
	flag.Var(
		&extCookies,
		"cookie",
		"extra cookies for request, can be used multiple times, accept files with '@'-prefix",
	)
	flag.Var(
		&tags,
		"tag",
		"tags filter, single or comma-separated tag names",
	)

	flag.Parse()

	if *fVersion {
		fmt.Println(version())

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
