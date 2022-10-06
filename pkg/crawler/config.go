package crawler

import (
	"fmt"
	"strings"
	"time"
)

const (
	minWorkers = 1
	maxWorkers = 64
	minDelay   = time.Duration(0)
)

type config struct {
	Headers    []string
	Cookies    []string
	AlowedTags []string
	Ignored    []string
	UserAgent  string
	Delay      time.Duration
	Workers    int
	Depth      int
	Robots     RobotsPolicy
	Dirs       DirsPolicy
	SkipSSL    bool
	Brute      bool
	NoHEAD     bool
	ScanJS     bool
}

func (c *config) validate() {
	switch {
	case c.Workers < minWorkers:
		c.Workers = minWorkers
	case c.Workers > maxWorkers:
		c.Workers = maxWorkers
	}

	if c.Delay < minDelay {
		c.Delay = minDelay
	}

	if c.Depth < 0 {
		c.Depth = -1
	}
}

func (c *config) String() (rv string) {
	var sb strings.Builder

	_, _ = sb.WriteString(fmt.Sprintf("workers: %d depth: %d brute: %t", c.Workers, c.Depth, c.Brute))

	if c.Delay > 0 {
		_, _ = sb.WriteString(fmt.Sprintf(" delay: %s", c.Delay))
	}

	return sb.String()
}
