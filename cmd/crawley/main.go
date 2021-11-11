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
)

const (
	appName      = "Crawley"
	appURL       = "https://github.com/s0rg/crawley"
	defaultDelay = 150 * time.Millisecond
)

var (
	gitHash       string
	gitVersion    string
	buildDate     string
	fVersion      = flag.Bool("version", false, "show version")
	fBrute        = flag.Bool("brute", false, "scan html comments")
	fSkipSSL      = flag.Bool("skip-ssl", false, "skip ssl verification")
	fSkipDirs     = flag.Bool("skip-dirs", false, "skip non-resource urls in output")
	fSilent       = flag.Bool("silent", false, "suppress info and error messages in stderr")
	fDepth        = flag.Int("depth", 0, "scan depth (-1 - unlimited)")
	fWorkers      = flag.Int("workers", runtime.NumCPU(), "number of workers")
	fDelay        = flag.Duration("delay", defaultDelay, "per-request delay (0 - disable)")
	fUA           = flag.String("user-agent", defaultAgent, "user-agent string")
	fRobotsPolicy = flag.String("robots", "ignore", "policy for robots.txt: ignore / crawl / respect")
	defaultAgent  = "Mozilla/5.0 (compatible; Win64; x64) Mr." + appName + "/" + gitVersion + "-" + gitHash
)

func callback(s string) {
	_, _ = os.Stdout.WriteString(s + "\n")
}

func crawl(uri string, opts ...crawler.Option) error {
	c := crawler.New(opts...)

	log.Printf("[*] config: %s", c.DumpConfig())
	log.Printf("[*] crawling url: %s", uri)

	if err := c.Run(uri, callback); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	log.Printf("[*] complete")

	return nil
}

func main() {
	flag.Parse()

	if *fVersion {
		fmt.Printf("%s %s git: %s build: %s site: %s\n", appName, gitVersion, gitHash, buildDate, appURL)

		return
	}

	if flag.NArg() != 1 {
		flag.Usage()

		return
	}

	policy, err := crawler.ParseRobotsPolicy(*fRobotsPolicy)
	if err != nil {
		log.Fatal("crawler:", err)
	}

	if *fSilent {
		log.SetOutput(io.Discard)
	}

	if err := crawl(
		flag.Arg(0),
		crawler.WithUserAgent(*fUA),
		crawler.WithRobotsPolicy(policy),
		crawler.WithDelay(*fDelay),
		crawler.WithMaxCrawlDepth(*fDepth),
		crawler.WithWorkersCount(*fWorkers),
		crawler.WithBruteMode(*fBrute),
		crawler.WithSkipSSL(*fSkipSSL),
		crawler.WithSkipDirs(*fSkipDirs),
	); err != nil {
		// forcing back stderr in case of errors, otherwise if 'silent' is on -
		// no one will know what happened.
		log.SetOutput(os.Stderr)
		log.Fatal("crawler:", err)
	}
}
