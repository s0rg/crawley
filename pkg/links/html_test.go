package links

import (
	"bytes"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const testRes1 = "http://test/result"

var testBase, _ = url.Parse("http://test/")

func TestClean(t *testing.T) {
	t.Parallel()

	type args struct {
		b *url.URL
		r string
	}

	testRes2 := "http://result/"
	testRes3 := "http://test/?foo=bar"

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

func TestExtractTag(t *testing.T) {
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
		name  string
		args  args
		wantU string
	}{
		{"ok", args{b: testBase, t: &tOK, k: "b"}, testRes1},
		{"bad", args{b: testBase, t: &tBAD, k: "a"}, ""},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotU := extractTag(tc.args.b, tc.args.t, tc.args.k)

			if !reflect.DeepEqual(gotU, tc.wantU) {
				t.Errorf("extractTag() gotU = %v, want %v", gotU, tc.wantU)
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
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
		wantURL  string
	}{
		{"no-link", tbad, "", "", ""},
		{"img-ok", html.Token{DataAtom: atom.Img, Attr: attrs}, "", "", testRes1},
		{"image-ok", html.Token{DataAtom: atom.Image, Attr: attrs}, "", "", testRes1},
		{"video-ok", html.Token{DataAtom: atom.Video, Attr: attrs}, "", keySRC, testRes1},
		{"audio-ok", html.Token{DataAtom: atom.Audio, Attr: attrs}, "", keySRC, testRes1},
		{"script-ok", html.Token{DataAtom: atom.Script, Attr: attrs}, "", "", testRes1},
		{"track-ok", html.Token{DataAtom: atom.Track, Attr: attrs}, "", "", testRes1},
		{"object-ok", html.Token{DataAtom: atom.Object, Attr: attrs}, "", "", testRes1},
		{"a-ok", html.Token{DataAtom: atom.A, Attr: attrs}, "", "", testRes1},
		{"iframe-ok", html.Token{DataAtom: atom.Iframe, Attr: attrs}, "", "", testRes1},
		{"video-empty", html.Token{DataAtom: atom.Video}, "", keySRC, ""},
		{"audio-empty", html.Token{DataAtom: atom.Audio}, "", keySRC, ""},
		{"picture-empty", html.Token{DataAtom: atom.Picture}, "", keySRCS, ""},
		{"source-src-ok", html.Token{DataAtom: atom.Source, Attr: attrs}, keySRC, keySRC, testRes1},
		{"form-action-ok", html.Token{DataAtom: atom.Form, Attr: attrs}, "", "", testRes1},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			key := tc.keyStart

			var res string

			extractToken(testBase, tc.token, &key, func(_ atom.Atom, s string) {
				res = s
			})

			if key != tc.keyWant {
				t.Errorf("extractToken() key gotU = %v, want %v", key, tc.keyWant)
			}

			if !reflect.DeepEqual(res, tc.wantURL) {
				t.Errorf("extractToken() link gotU = %v, want %v", res, tc.wantURL)
			}
		})
	}
}

func TestExtractURLS(t *testing.T) {
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
		lnk     string
	}{
		{"ok-1", raw1, true, testRes1},
		{"ok-4", raw4, true, testRes1},
		{"comment", raw3, true, testRes1},
		{"bad", raw2, false, ""},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBufferString(tc.raw)

			var res string

			ExtractHTML(buf, testBase, HTMLParams{
				Brute:  true,
				Filter: AllowALL,
				Handler: func(_ atom.Atom, s string) {
					res = s
				},
			})

			if tc.hasLink {
				if !reflect.DeepEqual(res, tc.lnk) {
					t.Errorf("extractToken() link gotU = %v, want %v", res, tc.lnk)
				}
			}
		})
	}
}

func TestExtractComment(t *testing.T) {
	const comment = `
loremipsumhTTp://foo fdfdfs HttPs://bar
       http://
 https://baz  http://boo"`

	t.Parallel()

	res := []string{}
	want := []string{"http://foo", "https://bar", "https://baz", "http://boo"}

	handler := func(_ atom.Atom, s string) {
		res = append(res, strings.ToLower(s))
	}

	extractComment(comment, handler)

	if len(res) != 4 {
		t.Error("unexpected len")
	}

	if !reflect.DeepEqual(res, want) {
		t.Error("unexpected result")
	}
}

func TestExtractAllowed(t *testing.T) {
	t.Parallel()

	const raw1 = `<html><a href="result-a">here</a><form action="result-form"></form></html>`

	buf := bytes.NewBufferString(raw1)

	var res []string

	filter := func(tkn html.Token) bool {
		return tkn.DataAtom == atom.A
	}

	ExtractHTML(buf, testBase, HTMLParams{
		Brute:  true,
		Filter: filter,
		Handler: func(_ atom.Atom, s string) {
			res = append(res, s)
		},
	})

	if len(res) != 1 {
		t.Fatal("unexpected result len")
	}

	if !strings.HasSuffix(res[0], "result-a") {
		t.Fatalf("unexpected result: %v", res)
	}
}
