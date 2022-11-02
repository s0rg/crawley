//go:build !test

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/s0rg/crawley/pkg/crawler"
	"github.com/s0rg/crawley/pkg/values"
)

const (
	appName      = "Crawley"
	appHelp      = "the unix-way web crawler"
	appSite      = "https://github.com/s0rg/crawley"
	defaultDelay = 150 * time.Millisecond
)

// build-time values.
var (
	GitTag    string
	GitHash   string
	BuildDate string
	defaultUA = "Mozilla/5.0 (compatible; Win64; x64) Mr." + appName + "/" + GitTag + "-" + GitHash
)

// command-line flags.
var (
	fDepth, fWorkers        int
	fSilent, fVersion       bool
	fBrute, fNoHeads        bool
	fSkipSSL, fScanJS       bool
	fDirsPolicy, fProxyAuth string
	fRobotsPolicy, fUA      string
	fDelay                  time.Duration
	cookies, headers        values.Smart
	tags, ignored           values.List
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

func usage() {
	var sb strings.Builder

	const twoCR = "\n\n"

	sb.WriteString(appName)
	sb.WriteString(" - ")
	sb.WriteString(appHelp)
	sb.WriteString(", usage:")
	sb.WriteString(twoCR)

	sb.WriteString(filepath.Base(os.Args[0]))
	sb.WriteString(" [flags] url")
	sb.WriteString(twoCR)

	sb.WriteString("possible flags with default values:")
	sb.WriteString(twoCR)

	_, _ = os.Stderr.WriteString(sb.String())

	flag.PrintDefaults()
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
	robots, err := crawler.ParseRobotsPolicy(fRobotsPolicy)
	if err != nil {
		err = fmt.Errorf("robots policy: %w", err)

		return
	}

	dirs, err := crawler.ParseDirsPolicy(fDirsPolicy)
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
		crawler.WithUserAgent(fUA),
		crawler.WithDelay(fDelay),
		crawler.WithMaxCrawlDepth(fDepth),
		crawler.WithWorkersCount(fWorkers),
		crawler.WithSkipSSL(fSkipSSL),
		crawler.WithBruteMode(fBrute),
		crawler.WithDirsPolicy(dirs),
		crawler.WithRobotsPolicy(robots),
		crawler.WithoutHeads(fNoHeads),
		crawler.WithScanJS(fScanJS),
		crawler.WithExtraHeaders(h),
		crawler.WithExtraCookies(c),
		crawler.WithTagsFilter(tags.Values),
		crawler.WithIgnored(ignored.Values),
		crawler.WithProxyAuth(fProxyAuth),
	}

	return rv, nil
}

func setupFlags() {
	flag.Var(&headers, "header",
		"extra headers for request, can be used multiple times, accept files with '@'-prefix",
	)
	flag.Var(&cookies, "cookie",
		"extra cookies for request, can be used multiple times, accept files with '@'-prefix",
	)

	flag.Var(&tags, "tag", "tags filter, single or comma-separated tag names")
	flag.Var(&ignored, "ignore", "patterns (in urls) to be ignored in crawl process")

	flag.IntVar(&fDepth, "depth", 0, "scan depth (set -1 for unlimited)")
	flag.IntVar(&fWorkers, "workers", runtime.NumCPU(), "number of workers")

	flag.BoolVar(&fBrute, "brute", false, "scan html comments")
	flag.BoolVar(&fNoHeads, "headless", false, "disable pre-flight HEAD requests")
	flag.BoolVar(&fScanJS, "js", false, "scan js files for endpoints")
	flag.BoolVar(&fSkipSSL, "skip-ssl", false, "skip ssl verification")
	flag.BoolVar(&fSilent, "silent", false, "suppress info and error messages in stderr")
	flag.BoolVar(&fVersion, "version", false, "show version")

	flag.StringVar(&fDirsPolicy, "dirs", crawler.DefaultDirsPolicy,
		"policy for non-resource urls: show / hide / only")
	flag.StringVar(&fRobotsPolicy, "robots", crawler.DefaultRobotsPolicy,
		"policy for robots.txt: ignore / crawl / respect")
	flag.StringVar(&fUA, "user-agent", defaultUA, "user-agent string")
	flag.StringVar(&fProxyAuth, "proxy-auth", "", "credentials for proxy: user:password")

	flag.DurationVar(&fDelay, "delay", defaultDelay, "per-request delay (0 - disable)")

	flag.Usage = usage
}

func main() {
	setupFlags()

	flag.Parse()

	if fVersion {
		puts(version())

		return
	}

	if flag.NArg() != 1 {
		usage()

		return
	}

	opts, err := initOptions()
	if err != nil {
		log.Fatal("[-] options:", err)
	}

	if fSilent {
		log.SetOutput(io.Discard)
	}

	if err := crawl(flag.Arg(0), opts...); err != nil {
		// forcing back stderr in case of errors, otherwise, if 'silent' is on - no one will know what happened.
		log.SetOutput(os.Stderr)
		log.Fatal("[-] crawler:", err)
	}
}
