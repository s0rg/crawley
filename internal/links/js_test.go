package links

import (
	"errors"
	"strings"
	"testing"
)

var errReader = errors.New("reader error")

type badReader struct {
	err error
}

func (r *badReader) Read(_ []byte) (n int, err error) {
	return 0, r.err
}

func TestExtractJSError(t *testing.T) {
	t.Parallel()

	var (
		r = badReader{err: errReader}
		c int
	)

	ExtractJS(&r, func(_ string) {
		c++
	})

	if c != 0 {
		t.Fatalf("unexpected result got: %d", c)
	}
}

func TestExtractJS(t *testing.T) {
	t.Parallel()

	const (
		js = `function() {
 		urls = [
			// invalid ones
			"user/create.notaext?user=Test",
			"text/html",
			"text/plain",
			"application/json",
			"api/create.php?user=test#home",
		    "api/create.php",
			"api/create.php?user=test"
		    "api/create.php?user=test&pass=test",
			"user/create.action?user=Test",
		    "api/user",
		    "test_1.json",
    		"v1/create",
    		"api/v1/user/2",
			"api/v1/search?text=Test Hello",
			"test2.aspx?arg1=tmp1+tmp2&arg2=tmp3",
   			"addUser.action",
    		"main.js",
    		"index.html",
    		"robots.txt",
    		"users.xml"
			// valid ones
			"smb://example.com",
			"http://example.com",
			"https://www.example.co.us",
			"/api/create.php?user=test&pass=test#home",
			"/path/to/file",
			"/user/create.action?user=Test"
		];
		}`
		count = 6
	)

	var c int

	ExtractJS(strings.NewReader(js), func(_ string) {
		c++
	})

	if c != count {
		t.Fatalf("unexpected result got: %d", c)
	}
}
