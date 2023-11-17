package crawler

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/s0rg/crawley/pkg/client"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	c := &config{}
	c.validate()

	if c.Client.Workers != minWorkers {
		t.Error("empty - bad workers")
	}

	c.Delay = time.Duration(-100)
	c.Depth = -5
	c.Client.Workers = 1000000
	c.Client.Timeout = time.Hour

	c.validate()

	if c.Client.Workers != maxWorkers {
		t.Error("non empty - bad workers")
	}

	if c.Delay != minDelay {
		t.Error("non empty - bad delay")
	}

	if c.Depth != -1 {
		t.Error("non empty - bad depth")
	}

	if c.Client.Timeout != maxTimeout {
		t.Error("non empty - bad timeout")
	}
}

func TestOptions(t *testing.T) {
	const (
		ua      = "foo"
		rp      = RobotsRespect
		dp      = DirsOnly
		delay   = time.Hour
		timeout = time.Minute * 5
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
		WithTagsFilter([]string{"a", "form"}),
		WithScanJS(fbool),
		WithIgnored([]string{"logout"}),
		WithTimeout(timeout),
	}

	c := &config{}

	for _, o := range opts {
		o(c)
	}

	c.validate()

	if c.Client.UserAgent != ua {
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

	if c.Client.Workers != workers {
		t.Error("bad workers")
	}

	if c.Brute != fbool {
		t.Error("bad brute mode")
	}

	if c.Client.SkipSSL != fbool {
		t.Error("bad skip-ssl")
	}

	if c.NoHEAD != fbool {
		t.Error("bad no-head")
	}

	if c.Dirs != dp {
		t.Error("bad dirs policy")
	}

	if !reflect.DeepEqual(c.Client.Headers, extHeaders) {
		t.Error("bad extra headers")
	}

	if !reflect.DeepEqual(c.Client.Cookies, extCookies) {
		t.Error("bad extra cookies")
	}

	if len(c.AlowedTags) != 2 {
		t.Error("unexpected filter size")
	}

	if c.ScanJS != fbool {
		t.Error("bad scan-js")
	}

	if len(c.Ignored) != 1 {
		t.Error("unexpected ignored size")
	}

	if c.Client.Timeout != timeout {
		t.Error("bad timeout")
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	c := &config{
		Client:  client.Config{Workers: 13},
		Depth:   666,
		Brute:   true,
		ScanJS:  true,
		ScanCSS: true,
	}

	c.validate()

	v := c.String()

	if !strings.Contains(v, "13") {
		t.Error("1 - bad workers")
	}

	if !strings.Contains(v, "666") {
		t.Error("1 - bad depth")
	}

	if !strings.Contains(v, "brute: on") {
		t.Error("1 - bad brute mode")
	}

	if !strings.Contains(v, "+js") {
		t.Error("1 - bad js mode")
	}

	if !strings.Contains(v, "+css") {
		t.Error("1 - bad css mode")
	}

	if strings.Contains(v, "delay") {
		t.Error("1 - delay found")
	}

	c = &config{
		Delay: time.Millisecond * 100,
	}

	c.validate()

	v = c.String()

	if strings.Contains(v, "brute") {
		t.Error("2 - bad brute mode")
	}

	if !strings.Contains(v, "100") {
		t.Error("2 - bad delay")
	}
}

func TestProxyAuth(t *testing.T) {
	t.Parallel()

	const creds = "user:pass"

	var (
		c       = &config{}
		opts    = []Option{WithProxyAuth(creds)}
		headers = []string{proxyAuthHeader(creds)}
	)

	for _, o := range opts {
		o(c)
	}

	c.validate()

	if !reflect.DeepEqual(c.Client.Headers, headers) {
		t.Fatalf("bad extra headers: %v", c.Client.Headers)
	}
}
