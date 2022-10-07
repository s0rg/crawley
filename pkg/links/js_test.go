package links

import (
	"errors"
	"strings"
	"testing"
)

func TestExtractJS(t *testing.T) {
	t.Parallel()

	const (
		js = `function() {
 		urls = [
			// invalid ones
			"user/create.notaext?user=Test",
			// valid ones
			"smb://example.com",
			"http://example.com",
			"https://www.example.co.us",
			"/api/create.php?user=test&pass=test#home",
			"api/create.php?user=test#home",
			"/path/to/file",
			"/user/create.action?user=Test"
		    "api/create.php",
			"api/create.php?user=test"
		    "api/create.php?user=test&pass=test",
			"user/create.action?user=Test",
		    "api/user",
    		"v1/create",
    		"api/v1/user/2",
			"api/v1/search?text=Test Hello",
		    "test_1.json",
			"test2.aspx?arg1=tmp1+tmp2&arg2=tmp3",
   			"addUser.action",
    		"main.js",
    		"index.html",
    		"robots.txt",
    		"users.xml"
		];
		}`
		count = 22
	)

	var c int

	ExtractJS(strings.NewReader(js), func(_ string) {
		c++
	})

	if c != count {
		t.Fatalf("unexpected result got: %d", c)
	}
}

type badReader struct{}

func (r *badReader) Read(b []byte) (n int, err error) {
	return 0, errors.New("reader error")
}

func TestExtractJSError(t *testing.T) {
	t.Parallel()

	var (
		r badReader
		c int
	)

	ExtractJS(&r, func(_ string) {
		c++
	})

	if c != 0 {
		t.Fatalf("unexpected result got: %d", c)
	}
}
