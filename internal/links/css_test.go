package links

import (
	"strings"
	"testing"
)

func TestExtractCSS(t *testing.T) {
	t.Parallel()

	const css = `
.background {
  overground: url();
  foreground: yellow;
  background: url("test.png");
}
`

	var c int

	ExtractCSS(strings.NewReader(css), func(_ string) {
		c++
	})

	if c != 1 {
		t.Fail()
	}
}
