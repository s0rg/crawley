package links

import (
	"bytes"
	"net/url"
	"slices"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const testRes1 = "http://test/result"

var testBase, _ = url.Parse("http://test/")

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
		{name: "ok", args: args{b: testBase, t: &tOK, k: "b"}, wantU: testRes1},
		{name: "bad", args: args{b: testBase, t: &tBAD, k: "a"}, wantU: ""},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotU := extractTag(tc.args.b, tc.args.t, tc.args.k)

			if gotU != tc.wantU {
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
		keyStart string
		keyWant  string
		wantURL  string
		token    html.Token
	}{
		{
			name:     "no-link",
			token:    tbad,
			keyStart: "",
			keyWant:  "",
			wantURL:  "",
		},
		{
			name:     "img-ok",
			token:    html.Token{DataAtom: atom.Img, Attr: attrs},
			keyStart: "",
			keyWant:  "",
			wantURL:  testRes1,
		},
		{
			name:     "image-ok",
			token:    html.Token{DataAtom: atom.Image, Attr: attrs},
			keyStart: "",
			keyWant:  "",
			wantURL:  testRes1,
		},
		{
			name:     "video-ok",
			token:    html.Token{DataAtom: atom.Video, Attr: attrs},
			keyStart: "",
			keyWant:  keySRC,
			wantURL:  testRes1,
		},
		{
			name:     "audio-ok",
			token:    html.Token{DataAtom: atom.Audio, Attr: attrs},
			keyStart: "",
			keyWant:  keySRC,
			wantURL:  testRes1,
		},
		{
			name:     "script-ok",
			token:    html.Token{DataAtom: atom.Script, Attr: attrs},
			keyStart: "",
			keyWant:  "",
			wantURL:  testRes1,
		},
		{
			name:     "track-ok",
			token:    html.Token{DataAtom: atom.Track, Attr: attrs},
			keyStart: "",
			keyWant:  "",
			wantURL:  testRes1,
		},
		{
			name:     "object-ok",
			token:    html.Token{DataAtom: atom.Object, Attr: attrs},
			keyStart: "",
			keyWant:  "",
			wantURL:  testRes1,
		},
		{
			name:     "a-ok",
			token:    html.Token{DataAtom: atom.A, Attr: attrs},
			keyStart: "",
			keyWant:  "",
			wantURL:  testRes1,
		},
		{
			name:     "iframe-ok",
			token:    html.Token{DataAtom: atom.Iframe, Attr: attrs},
			keyStart: "",
			keyWant:  "",
			wantURL:  testRes1,
		},
		{
			name:     "audio-empty",
			token:    html.Token{DataAtom: atom.Audio},
			keyStart: "",
			keyWant:  keySRC,
			wantURL:  "",
		},
		{
			name:     "audio-empty",
			token:    html.Token{DataAtom: atom.Audio},
			keyStart: "",
			keyWant:  keySRC,
			wantURL:  "",
		},
		{
			name:     "picture-empty",
			token:    html.Token{DataAtom: atom.Picture},
			keyStart: "",
			keyWant:  keySRCS,
			wantURL:  "",
		},
		{
			name:     "source-src-ok",
			token:    html.Token{DataAtom: atom.Source, Attr: attrs},
			keyStart: keySRC,
			keyWant:  keySRC,
			wantURL:  testRes1,
		},
		{
			name:     "form-action-ok",
			token:    html.Token{DataAtom: atom.Form, Attr: attrs},
			keyStart: "",
			keyWant:  "",
			wantURL:  testRes1,
		},
		{
			name:     "css-link-ok",
			token:    html.Token{DataAtom: atom.Link, Attr: attrs},
			keyStart: keySRC,
			keyWant:  keySRC,
			wantURL:  testRes1,
		},
		{
			name:     "css-style-empty",
			token:    html.Token{DataAtom: atom.Style},
			keyStart: "",
			keyWant:  "",
			wantURL:  "",
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			key := tc.keyStart

			var res string

			extractToken(
				testBase,
				tc.token,
				&key,
				func(_ atom.Atom, s string) {
					res = s
				},
			)

			if key != tc.keyWant {
				t.Errorf("extractToken() key gotU = %v, want %v", key, tc.keyWant)
			}

			if res != tc.wantURL {
				t.Errorf("extractToken() link gotU = %v, want %v", res, tc.wantURL)
			}
		})
	}
}

func TestExtractTokenJS(t *testing.T) {
	t.Parallel()

	const raw = `<html><script>var url = "http://example.com";</script></html>`

	var res string

	ExtractHTML(
		bytes.NewBufferString(raw),
		testBase,
		HTMLParams{
			Filter: AllowALL,
			ScanJS: true,
			HandleStatic: func(s string) {
				res = s
			},
		},
	)

	if res != "http://example.com" {
		t.Fail()
	}
}

func TestExtractTokenCSS(t *testing.T) {
	t.Parallel()

	const raw = `<html><style>foo {bar:url(test.png);}</style></html>`

	var res string

	ExtractHTML(
		bytes.NewBufferString(raw),
		testBase,
		HTMLParams{
			Filter:  AllowALL,
			ScanCSS: true,
			HandleStatic: func(s string) {
				res = s
			},
		},
	)

	if !strings.HasSuffix(res, "test.png") {
		t.Fail()
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
		lnk     string
		hasLink bool
	}{
		{
			name:    "ok-1",
			raw:     raw1,
			hasLink: true,
			lnk:     testRes1,
		},
		{
			name:    "ok-4",
			raw:     raw4,
			hasLink: true,
			lnk:     testRes1,
		},
		{
			name:    "comment",
			raw:     raw3,
			hasLink: true,
			lnk:     testRes1,
		},
		{
			name:    "bad",
			raw:     raw2,
			hasLink: false,
			lnk:     "",
		},
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
				HandleHTML: func(_ atom.Atom, s string) {
					res = s
				},
			})

			if tc.hasLink {
				if res != tc.lnk {
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

	if slices.Compare(res, want) != 0 {
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
		HandleHTML: func(_ atom.Atom, s string) {
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
