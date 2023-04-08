package client

import (
	"net/http"
	"reflect"
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

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("prepareCookies(): %v, want: %v", got, tt.want)
		}
	}
}
