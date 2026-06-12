package client

import (
	"net/http"
	"testing"
)

func Test_prepareCookies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want []*http.Cookie
	}{
		{
			"1",
			[]string{"NAME1=VALUE1;NAME2=ENCODED%20VALUE;", "NAME3=VALUE3"},
			[]*http.Cookie{
				{Name: "NAME1", Value: "VALUE1"},
				{Name: "NAME2", Value: "ENCODED%20VALUE"},
				{Name: "NAME3", Value: "VALUE3"},
			},
		},
		{
			"2",
			[]string{"", "NAME=", "=VALUE", ";;", "===", " VALID = COOKIE "},
			[]*http.Cookie{
				{Name: "NAME", Value: ""},
				{Name: "VALID", Value: "COOKIE"},
			},
		},
		{
			"3",
			[]string{"some_file.txt"},
			[]*http.Cookie{},
		},
	}

	for _, tt := range tests {
		got := prepareCookies(tt.args)

		if len(got) != len(tt.want) {
			t.Errorf("prepareCookies() invalid result count for: %v", tt.want)
		}

		if len(got) == 0 {
			continue
		}

		cookiesEqual(t, tt.want, got)
	}
}

func cookiesEqual(t *testing.T, a, b []*http.Cookie) {
	t.Helper()

	if len(b) > len(a) {
		a, b = b, a
	}

	index := make(map[string]string)

	for _, v := range b {
		index[v.Name] = v.Value
	}

	for _, v := range a {
		bv, ok := index[v.Name]
		if !ok {
			continue
		}

		if bv != v.Value {
			t.Fail()
		}
	}
}
