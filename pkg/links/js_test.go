package links

import (
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

	ExtractJS(strings.NewReader(js), func(uri string) {
		c++
		t.Logf("found: %s", uri)
	})

	if c != count {
		t.Fatalf("unexpected result got: %d", c)
	}
}
