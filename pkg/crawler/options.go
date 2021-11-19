package crawler

import "time"

// Option is a configuration func.
type Option func(*config)

// WithUserAgent sets User-Agent string.
func WithUserAgent(v string) Option {
	return func(c *config) {
		c.UserAgent = v
	}
}

// WithDelay sets crawl delay.
func WithDelay(v time.Duration) Option {
	return func(c *config) {
		c.Delay = v
	}
}

// WithMaxCrawlDepth sets maximum depth to crawl.
func WithMaxCrawlDepth(v int) Option {
	return func(c *config) {
		c.Depth = v
	}
}

// WithWorkersCount sets maximum workers.
func WithWorkersCount(v int) Option {
	return func(c *config) {
		c.Workers = v
	}
}

// WithRobotsPolicy sets RobotsPolicy for crawler.
func WithRobotsPolicy(v RobotsPolicy) Option {
	return func(c *config) {
		c.Robots = v
	}
}

// WithDirsPolicy sets DirsPolicy for crawler.
func WithDirsPolicy(v DirsPolicy) Option {
	return func(c *config) {
		c.Dirs = v
	}
}

// WithSkipSSL tells crawley to skip any ssl handshake errors.
func WithSkipSSL(v bool) Option {
	return func(c *config) {
		c.SkipSSL = v
	}
}

// WithBruteMode enables "brute-mode" - html comments scan.
func WithBruteMode(v bool) Option {
	return func(c *config) {
		c.Brute = v
	}
}

// WithoutHeads disables pre-flight HEAD requests.
func WithoutHeads(v bool) Option {
	return func(c *config) {
		c.NoHEAD = v
	}
}
