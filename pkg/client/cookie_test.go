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
		{"1",
			[]string{"NAME1=VALUE1; NAME2=VALUE2", "NAME3=VALUE3"},
			[]*http.Cookie{
				{Name: "NAME1", Value: "VALUE1"},
				{Name: "NAME2", Value: "VALUE2"},
				{Name: "NAME3", Value: "VALUE3"},
			}},
		{"2",
			[]string{"", "NAME=", "=VALUE", ";;", "===", " VALID = COOKIE "},
			[]*http.Cookie{
				{Name: "VALID", Value: "COOKIE"},
			}},
	}

	for _, tt := range tests {
		if got := prepareCookies(tt.args); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("prepareCookies() = %v, want %v", got, tt.want)
		}
	}
}
