package crawler

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	c := &config{}
	c.validate()

	if c.Workers != minWorkers {
		t.Error("empty - bad workers")
	}

	c.Workers = 1000000
	c.Delay = time.Duration(-100)
	c.Depth = -5

	c.validate()

	if c.Workers != maxWorkers {
		t.Error("non empty - bad workers")
	}

	if c.Delay != minDelay {
		t.Error("non empty - bad delay")
	}

	if c.Depth != -1 {
		t.Error("non empty - bad depth")
	}
}

func TestOptions(t *testing.T) {
	const (
		ua      = "foo"
		rp      = RobotsRespect
		dp      = DirsOnly
		delay   = time.Hour
		workers = 13
		depth   = 666
		fbool   = true
	)

	var (
		extHeaders = []string{"foo: bar"}
		extCookies = []string{"name=val"}
	)

	t.Parallel()

	opts := []Option{
		WithUserAgent(ua),
		WithRobotsPolicy(rp),
		WithDirsPolicy(dp),
		WithDelay(delay),
		WithMaxCrawlDepth(depth),
		WithWorkersCount(workers),
		WithBruteMode(fbool),
		WithSkipSSL(fbool),
		WithoutHeads(fbool),
		WithExtraHeaders(extHeaders),
		WithExtraCookies(extCookies),
	}

	c := &config{}

	for _, o := range opts {
		o(c)
	}

	c.validate()

	if c.UserAgent != ua {
		t.Error("bad user-agent")
	}

	if c.Robots != RobotsRespect {
		t.Error("bad policy")
	}

	if c.Delay != delay {
		t.Error("bad delay")
	}

	if c.Depth != depth {
		t.Error("bad depth")
	}

	if c.Workers != workers {
		t.Error("bad workers")
	}

	if c.Brute != fbool {
		t.Error("bad brute mode")
	}

	if c.SkipSSL != fbool {
		t.Error("bad skip-ssl")
	}

	if c.NoHEAD != fbool {
		t.Error("bad no-head")
	}

	if c.Dirs != dp {
		t.Error("bad dirs policy")
	}

	if !reflect.DeepEqual(c.Headers, extHeaders) {
		t.Error("bad extra headers")
	}

	if !reflect.DeepEqual(c.Cookies, extCookies) {
		t.Error("bad extra cookies")
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	c := &config{
		Workers: 13,
		Depth:   666,
		Brute:   true,
	}

	c.validate()

	v := c.String()

	if !strings.Contains(v, "13") {
		t.Error("1 - bad workers")
	}

	if !strings.Contains(v, "666") {
		t.Error("1 - bad depth")
	}

	if !strings.Contains(v, "true") {
		t.Error("1 - bad brute mode")
	}

	if strings.Contains(v, "delay") {
		t.Error("1 - delay found")
	}

	c = &config{
		Delay: time.Millisecond * 100,
	}

	c.validate()

	v = c.String()

	if !strings.Contains(v, "false") {
		t.Error("2 - bad brute mode")
	}

	if !strings.Contains(v, "100") {
		t.Error("2 - bad delay")
	}
}
