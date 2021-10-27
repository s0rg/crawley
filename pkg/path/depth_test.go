package path

import (
	"reflect"
	"testing"
)

func Test_dropSpaces(t *testing.T) {
	type args struct {
		s []string
	}

	tests := []struct {
		name  string
		args  args
		wantO []string
	}{
		{"ones", args{s: []string{"", "1", ""}}, []string{"1"}},
		{"twos", args{s: []string{"2", "", "2"}}, []string{"2", "2"}},
		{"threes", args{s: []string{"", "", "3"}}, []string{"3"}},
		{"empty", args{s: []string{"", "", ""}}, []string{}},
	}

	t.Parallel()

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if gotO := dropSpaces(tc.args.s); !reflect.DeepEqual(gotO, tc.wantO) {
				t.Errorf("dropSpaces() = %v, want %v", gotO, tc.wantO)
			}
		})
	}
}

func Test_splitPath(t *testing.T) {
	type args struct {
		p string
	}

	tests := []struct {
		name  string
		args  args
		wantO []string
	}{
		{"empty", args{p: "/"}, []string{}},
		{"foo", args{p: "/foo"}, []string{"foo"}},
		{"foo-bar", args{p: "/foo/bar"}, []string{"foo", "bar"}},
		{"foo-bar-baz", args{p: "/foo/bar//baz"}, []string{"foo", "bar", "baz"}},
	}

	t.Parallel()

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if gotO := splitPath(tc.args.p); !reflect.DeepEqual(gotO, tc.wantO) {
				t.Errorf("splitPath() = %v, want %v", gotO, tc.wantO)
			}
		})
	}
}

func Test_Depth(t *testing.T) {
	type args struct {
		base string
		sub  string
	}

	tests := []struct {
		name      string
		args      args
		wantDepht int
		wantFound bool
	}{
		{"a-ok", args{base: "/", sub: "/a"}, 1, true},
		{"a-bad", args{base: "/a", sub: "/b"}, 0, false},
		{"c-bad", args{base: "/a/b", sub: "/c"}, 0, false},
		{"b-ok", args{base: "/a", sub: "/a/b"}, 1, true},
		{"c-ok", args{base: "/a", sub: "/a/b/c"}, 2, true},
	}

	t.Parallel()

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotDepht, gotFound := Depth(tc.args.base, tc.args.sub)
			if gotDepht != tc.wantDepht {
				t.Errorf("pathDepth() gotDepht = %v, want %v", gotDepht, tc.wantDepht)
			}
			if gotFound != tc.wantFound {
				t.Errorf("pathDepth() gotFound = %v, want %v", gotFound, tc.wantFound)
			}
		})
	}
}
