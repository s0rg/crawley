//go:build !test

package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/s0rg/crawley/pkg/crawler"
)

const (
	appName      = "Crawley"
	appURL       = "https://github.com/s0rg/crawley"
	minWorkers   = 1
	maxWorkers   = 64
	minDepth     = 0
	minDelay     = 50 * time.Millisecond
	defaultDelay = 250 * time.Millisecond
)

var (
	gitHash      string
	gitVersion   string
	buildDate    string
	defaultAgent = "Mozilla/5.0 (compatible; Win64; x64) Mr. " + appName + "/" + gitVersion + "-" + gitHash
	fVersion     = flag.Bool("version", false, "show version")
	fSkipSSL     = flag.Bool("skip-ssl", false, "skip ssl verification")
	fDepth       = flag.Int("depth", 0, "scan depth")
	fWorkers     = flag.Int("workers", runtime.NumCPU(), "number of workers")
	fDelay       = flag.Duration("delay", defaultDelay, "per-request delay")
	fUA          = flag.String("user-agent", defaultAgent, "user-agent string")
	fRobots      = flag.String("robots", "ignore", "action for robots.txt: ignore/crawl/respect")
)

func actionFromString(s string) (a crawler.RobotsAction, err error) {
	switch s {
	case "ignore":
		a = crawler.RobotsIgnore
	case "crawl":
		a = crawler.RobotsCrawl
	case "respect":
		a = crawler.RobotsRespect
	default:
		err = fmt.Errorf("unknown: '%s'", s)

		return
	}

	return a, nil
}

func sanitize(workers, depth int, delay time.Duration) (wrk, dep int, dur time.Duration) {
	wrk = minWorkers
	if workers > wrk {
		wrk = maxWorkers
		if workers < maxWorkers {
			wrk = workers
		}
	}

	if dep = minDepth; dep < depth {
		dep = depth
	}

	if dur = minDelay; dur < delay {
		dur = delay
	}

	return wrk, dep, dur
}

func printer(s string) {
	fmt.Println(s)
}

func crawl(
	uri, ua, action string,
	workers, depth int,
	delay time.Duration,
	skipSSL bool,
) error {
	act, err := actionFromString(action)
	if err != nil {
		return fmt.Errorf("action: %w", err)
	}

	c := crawler.New(ua, workers, depth, delay, skipSSL, act)

	log.Printf("%s started for: %s", appName, uri)
	log.Printf("workers: %d depth: %d delay: %s", workers, depth, delay)

	if err := c.Run(uri, printer); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	log.Printf("complete for: %s", uri)

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

	wnum, depth, delay := sanitize(*fWorkers, *fDepth, *fDelay)

	if err := crawl(flag.Arg(0), *fUA, *fRobots, wnum, depth, delay, *fSkipSSL); err != nil {
		log.Fatal("crawler:", err)
	}
}
