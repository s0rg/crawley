package links

import (
	"net/url"
	"reflect"
	"testing"
)

func TestClean(t *testing.T) {
	t.Parallel()

	type args struct {
		b *url.URL
		r string
	}

	const (
		testRes2 = "http://result/"
		testRes3 = "http://test/?foo=bar"
		testRes4 = "http://test/api/v1/user"
	)

	tests := []struct {
		name   string
		args   args
		wantU  string
		wantOk bool
	}{
		{"bad-uri", args{b: testBase, r: "[%]"}, "", false},
		{"empty-uri", args{b: testBase, r: "http://"}, "", false},
		{"js-scheme", args{b: testBase, r: "javascript://result"}, "", false},
		{"result-ok", args{b: testBase, r: "result"}, testRes1, true},
		{"result-no-scheme", args{b: testBase, r: "//result"}, testRes2, true},
		{"result-params", args{b: testBase, r: "/?foo=bar"}, testRes3, true},
		{"result-api", args{b: testBase, r: "/api/v1/user"}, testRes4, true},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotU, gotOk := cleanURL(tc.args.b, tc.args.r)
			if gotOk != tc.wantOk {
				t.Errorf("clean() gotOk = %v, want %v", gotOk, tc.wantOk)
			}

			if gotOk {
				if !reflect.DeepEqual(gotU, tc.wantU) {
					t.Errorf("clean() gotU = %v, want %v", gotU, tc.wantU)
				}
			}
		})
	}
}
