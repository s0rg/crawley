package links

import (
	"net/url"
	"strings"
	"testing"
)

func TestExtractJS(t *testing.T) {
	t.Parallel()

	const (
		js = `function() {
 		urls = [
			"http://example.com",
			"smb://example.com",
    		"https://www.example.co.us",
			"/path/to/file",
			"/user/create.action?user=Test"
			"/api/create.php?user=test&pass=test#home",
		    "api/create.php",
			"api/create.php?user=test"
		    "api/create.php?user=test&pass=test",
			"api/create.php?user=test#home",
			"user/create.action?user=Test",
			"user/create.notaext?user=Test",
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
		]
		}`
		count = 23
	)

	u, _ := url.Parse("http://HOST")

	var c int

	ExtractJS(strings.NewReader(js), u, func(_ string) {
		c++
	})

	if c != count {
		t.Fatalf("unexpected result got: %d", c)
	}
}
