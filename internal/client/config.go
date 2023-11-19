package client

import "time"

type Config struct {
	UserAgent string
	Headers   []string
	Cookies   []string
	Workers   int
	Timeout   time.Duration
	SkipSSL   bool
}
