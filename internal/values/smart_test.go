package values

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestSmartSet(t *testing.T) {
	t.Parallel()

	var (
		l   Smart
		err error
		res []string
	)

	if err = l.Set("a"); err != nil {
		t.Fatalf("set a - unexpected error: %v", err)
	}

	if res, err = l.Load(nil); err != nil {
		t.Fatalf("load a - unexpected error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("len a - unexpected length: %d", len(res))
	}

	if res[0] != "a" {
		t.Fatalf("res a - unexpected value: %v", res[0])
	}

	if err = l.Set("b"); err != nil {
		t.Fatalf("set b - unexpected error: %v", err)
	}

	if res, err = l.Load(nil); err != nil {
		t.Fatalf("load b - unexpected error: %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("len b - unexpected length: %d", len(res))
	}

	if res[1] != "b" {
		t.Fatalf("res b - unexpected value: %v", res[1])
	}
}

func TestSmartString(t *testing.T) {
	t.Parallel()

	var l Smart

	if l.String() != "" {
		t.Fatal("non-empty result")
	}
}

func TestSmartLoad(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		"foo": {Data: []byte("foo1\nfoo2")},
		"bar": {Data: []byte("bar1\nbar2")},
	}

	var l Smart

	_ = l.Set("foo0")
	_ = l.Set("@foo")

	res, err := l.Load(fsys)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 3 {
		t.Fatal("unexpexted length", len(res))
	}

	if res[0] != "foo0" {
		t.Fatal("unexpexted value 0")
	}

	if res[1] != "foo1" {
		t.Fatal("unexpexted value 1")
	}

	if res[2] != "foo2" {
		t.Fatal("unexpexted value 2")
	}
}

func TestSmartLoadFSErrorDir(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		"foo": {Mode: 0o777 | fs.ModeDir},
	}

	var l Smart

	_ = l.Set("foo0")
	_ = l.Set("@foo")

	_, err := l.Load(fsys)
	if err == nil {
		t.Fatal("unexepected nil-error")
	}
}
