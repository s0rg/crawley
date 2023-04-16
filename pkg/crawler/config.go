package crawler

import (
	"fmt"
	"strings"
	"time"

	"github.com/s0rg/crawley/pkg/client"
)

const (
	minWorkers = 1
	maxWorkers = 64
	minDelay   = time.Duration(0)
	minTimeout = time.Second
	maxTimeout = time.Minute * 10
)

type config struct {
	AlowedTags []string
	Ignored    []string
	Client     client.Config
	Delay      time.Duration
	Depth      int
	Robots     RobotsPolicy
	Dirs       DirsPolicy
	Brute      bool
	NoHEAD     bool
	ScanJS     bool
}

func (c *config) validate() {
	switch {
	case c.Client.Workers < minWorkers:
		c.Client.Workers = minWorkers
	case c.Client.Workers > maxWorkers:
		c.Client.Workers = maxWorkers
	}

	switch {
	case c.Client.Timeout < minTimeout:
		c.Client.Timeout = minTimeout
	case c.Client.Timeout > maxTimeout:
		c.Client.Timeout = maxTimeout
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

	_, _ = sb.WriteString(fmt.Sprintf("workers: %d depth: %d timeout: %s", c.Client.Workers, c.Depth, c.Client.Timeout))

	if c.Brute {
		_, _ = sb.WriteString(" brute: on")
	}

	if c.ScanJS {
		_, _ = sb.WriteString(" js: on")
	}

	if c.Delay > 0 {
		_, _ = sb.WriteString(fmt.Sprintf(" delay: %s", c.Delay))
	}

	return sb.String()
}
