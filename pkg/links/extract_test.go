package links

import (
	"bytes"
	"io"
	"net/url"
	"reflect"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	testBase, _ = url.Parse("http://test/")
	testRes1, _ = url.Parse("http://test/result")
)

func Test_clean(t *testing.T) {
	t.Parallel()

	type args struct {
		b *url.URL
		r string
	}

	testRes2, _ := url.Parse("http://result/")
	testRes3, _ := url.Parse("http://test/?foo=bar")

	tests := []struct {
		name   string
		args   args
		wantU  *url.URL
		wantOk bool
	}{
		{"bad-uri", args{b: testBase, r: "[%]"}, nil, false},
		{"empty-uri", args{b: testBase, r: "http://"}, nil, false},
		{"js-scheme", args{b: testBase, r: "javascript://result"}, nil, false},
		{"result-ok", args{b: testBase, r: "result"}, testRes1, true},
		{"result-no-scheme", args{b: testBase, r: "//result"}, testRes2, true},
		{"result-params", args{b: testBase, r: "/?foo=bar"}, testRes3, true},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotU, gotOk := clean(tc.args.b, tc.args.r)
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

func Test_extractTag(t *testing.T) {
	t.Parallel()

	type args struct {
		b *url.URL
		t *html.Token
		k string
	}

	tBAD := html.Token{}
	tOK := html.Token{
		Attr: []html.Attribute{
			{Key: "a", Val: "key"},
			{Key: "b", Val: "result"},
		},
	}

	tests := []struct {
		name   string
		args   args
		wantU  *url.URL
		wantOk bool
	}{
		{"ok", args{b: testBase, t: &tOK, k: "b"}, testRes1, true},
		{"bad", args{b: testBase, t: &tBAD, k: "a"}, nil, false},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotU, gotOk := extractTag(tc.args.b, tc.args.t, tc.args.k)
			if gotOk != tc.wantOk {
				t.Errorf("extractTag() gotOk = %v, want %v", gotOk, tc.wantOk)
			}

			if gotOk {
				if !reflect.DeepEqual(gotU, tc.wantU) {
					t.Errorf("extractTag() gotU = %v, want %v", gotU, tc.wantU)
				}
			}
		})
	}
}

func Test_extractToken(t *testing.T) {
	t.Parallel()

	tbad := html.Token{}
	attrs := []html.Attribute{
		{Key: keySRC, Val: "result"},
		{Key: keySRCS, Val: "result"},
		{Key: keyHREF, Val: "result"},
		{Key: keyDATA, Val: "result"},
		{Key: keyACTION, Val: "result"},
	}

	tests := []struct {
		name     string
		token    html.Token
		keyStart string
		keyWant  string
		hasURL   bool
		wantURL  *url.URL
	}{
		{"no-link", tbad, "", "", false, nil},
		{"img-ok", html.Token{DataAtom: atom.Img, Attr: attrs}, "", "", true, testRes1},
		{"image-ok", html.Token{DataAtom: atom.Image, Attr: attrs}, "", "", true, testRes1},
		{"video-ok", html.Token{DataAtom: atom.Video, Attr: attrs}, "", keySRC, true, testRes1},
		{"audio-ok", html.Token{DataAtom: atom.Audio, Attr: attrs}, "", keySRC, true, testRes1},
		{"script-ok", html.Token{DataAtom: atom.Script, Attr: attrs}, "", "", true, testRes1},
		{"track-ok", html.Token{DataAtom: atom.Track, Attr: attrs}, "", "", true, testRes1},
		{"object-ok", html.Token{DataAtom: atom.Object, Attr: attrs}, "", "", true, testRes1},
		{"a-ok", html.Token{DataAtom: atom.A, Attr: attrs}, "", "", true, testRes1},
		{"iframe-ok", html.Token{DataAtom: atom.Iframe, Attr: attrs}, "", "", true, testRes1},
		{"video-empty", html.Token{DataAtom: atom.Video}, "", keySRC, false, testRes1},
		{"audio-empty", html.Token{DataAtom: atom.Audio}, "", keySRC, false, testRes1},
		{"picture-empty", html.Token{DataAtom: atom.Picture}, "", keySRCS, false, testRes1},
		{"source-src-ok", html.Token{DataAtom: atom.Source, Attr: attrs}, keySRC, keySRC, true, testRes1},
		{"form-action-ok", html.Token{DataAtom: atom.Form, Attr: attrs}, "", "", true, testRes1},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			key := tc.keyStart

			var res *url.URL

			extractToken(testBase, tc.token, &key, func(_ atom.Atom, u *url.URL) {
				res = u
			})

			if key != tc.keyWant {
				t.Errorf("extractToken() key gotU = %v, want %v", key, tc.keyWant)
			}

			if tc.hasURL {
				if !reflect.DeepEqual(res, tc.wantURL) {
					t.Errorf("extractToken() link gotU = %v, want %v", res, tc.wantURL)
				}
			}
		})
	}
}

func Test_extractURLS(t *testing.T) {
	t.Parallel()

	const (
		raw1 = `<html><a href="result">here</a></html>`
		raw2 = `<html><video></video></html>`
		raw3 = `<html><!-- http://test/result --></html>`
		raw4 = `<html><form action="result"></form></html>`
	)

	tests := []struct {
		name    string
		raw     string
		hasLink bool
		lnk     *url.URL
	}{
		{"ok-1", raw1, true, testRes1},
		{"ok-4", raw4, true, testRes1},
		{"comment", raw3, true, testRes1},
		{"bad", raw2, false, nil},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBufferString(tc.raw)

			var res *url.URL

			Extract(testBase, io.NopCloser(buf), true, func(_ atom.Atom, u *url.URL) {
				res = u
			})

			if tc.hasLink {
				if !reflect.DeepEqual(res, tc.lnk) {
					t.Errorf("extractToken() link gotU = %v, want %v", res, tc.lnk)
				}
			}
		})
	}
}

func Test_extractComment(t *testing.T) {
	const comment = `
loremipsumhTTp://foo fdfdfs HttPs://bar
       http://
 https://baz  http://boo"`

	t.Parallel()

	res := []string{}
	want := []string{"foo", "bar", "baz", "boo"}

	handler := func(_ atom.Atom, u *url.URL) {
		res = append(res, u.Host)
	}

	extractComment(comment, handler)

	if len(res) != 4 {
		t.Error("unexpected len")
	}

	if !reflect.DeepEqual(res, want) {
		t.Error("unexpected result")
	}
}
