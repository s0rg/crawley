package client

import (
	"reflect"
	"testing"
)

func Test_prepareHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want []*header
	}{
		{
			"1",
			[]string{"foo: bar", "test: me"},
			[]*header{
				{"foo", "bar"},
				{"test", "me"},
			},
		},
		{
			"2",
			[]string{"   one: 1", "junk-key:", "two   : 2  ", ":junk-val"},
			[]*header{
				{"one", "1"},
				{"two", "2"},
			},
		},
	}

	for _, tt := range tests {
		if got := prepareHeaders(tt.args); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("prepareHeaders() = %v, want %v", got, tt.want)
		}
	}
}
