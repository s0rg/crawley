package robots

import (
	"testing"
)

func Test_Parser(t *testing.T) {
	const raw = `useragent: a
# some comment : with colon
disallow: /c
user-agent: b
disallow: /d
user-agent: e
sitemap: http://test.com/c
user-agent: f
disallow: /g`

}
