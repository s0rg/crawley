package path

import (
	"testing"
)

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
		{"d-bad", args{base: "/a/b/c", sub: "/d/b/c/a"}, 0, false},
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

func Benchmark_Depth(b *testing.B) {
	const (
		x = "/some/rather/long/path"
		y = "/some/rather/long/path/but/longer"
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Depth(x, y)
	}
}
