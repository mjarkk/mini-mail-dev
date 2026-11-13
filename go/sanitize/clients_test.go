package sanitize

import (
	"strings"
	"testing"
)

type expectation struct {
	in            string
	want, wantNot []string
}

func check(t *testing.T, client *Client, cases []expectation) {
	t.Helper()
	for _, c := range cases {
		got := mustParse(t, c.in, client)
		// The client's default styles would otherwise satisfy or break most expectations.
		got = strings.Replace(got, "<style data-client-defaults>"+client.baseCSS+"</style>", "", 1)
		for _, want := range c.want {
			if !strings.Contains(got, want) {
				t.Errorf("%s: Parse(%q)\n got: %q\nwant it to contain %q", client.ID, c.in, got, want)
			}
		}
		for _, wantNot := range c.wantNot {
			if strings.Contains(got, wantNot) {
				t.Errorf("%s: Parse(%q)\n got: %q\nwant it not to contain %q", client.ID, c.in, got, wantNot)
			}
		}
	}
}

func TestGmail(t *testing.T) {
	check(t, gmail, []expectation{
		{
			in:      `<!DOCTYPE html><html><head><style>.a{color:red}</style></head><body bgcolor="#eee" style="color:blue"><style>.b{color:green}</style><p class="a">x</p></body></html>`,
			want:    []string{`<!DOCTYPE html>`, `.a{color:red;}`, `<u></u><div style="color:blue">`},
			wantNot: []string{`.b{`, `bgcolor`, `<body`, `<html`},
		},
		{
			in:      `<p style="color:red;position:absolute;-webkit-text-size-adjust:none;box-shadow:0 0 1px #000">x</p>`,
			want:    []string{`style="color:red"`},
			wantNot: []string{`position`, `text-size-adjust`, `box-shadow`},
		},
		{
			in:      `<p style="color:red !IMPORTANT;margin:0">x</p><p style="color:rgb(0 0 0)">y</p>`,
			wantNot: []string{`style=`},
		},
		{
			in:      `<head><style>@import url(a.css);@font-face{font-family:x;src:url(x.woff)}@media (prefers-color-scheme:dark){.a{color:#fff}}@media screen and (max-width:600px){.b{color:red}}.c:hover{color:red}.d:first-child{color:red}[data-x]{color:red}</style></head>`,
			want:    []string{`@media screen and (max-width:600px){.b{color:red;}}`, `.c:hover{color:red;}`},
			wantNot: []string{`@import`, `@font-face`, `prefers-color-scheme`, `first-child`, `data-x`},
		},
		{
			in:      `<head><style>.a{color:red !Important}</style></head>`,
			wantNot: []string{`<style>`},
		},
		{
			in:      `<head><style>.a{margin:-4px 0;--x:red;color:var(--x)}</style></head><p data-a="1" hidden>x</p>`,
			wantNot: []string{`margin`, `--x`, `data-a`, `hidden`},
		},
		{
			in:   `<img src="data:image/png;base64,AAAA" alt="a"><img src="https://example.com/a.png" srcset="https://example.com/b.png 2x">`,
			want: []string{`<img alt="a">`, `<img src="https://example.com/a.png">`},
		},
		{
			in:      `<details><summary>s</summary>c</details>`,
			want:    []string{`sc`},
			wantNot: []string{`<details`, `<summary`},
		},
	})
}

func TestGmailStyleLimit(t *testing.T) {
	big := strings.Repeat(".a{color:red}", 2000)
	in := `<head><style>.small{color:red}</style><style>` + big + `</style><style>.after{color:red}</style></head>`
	check(t, gmail, []expectation{{in: in, want: []string{`.small{`}, wantNot: []string{`.after{`, big[:100]}}})
}

func TestGmailClipping(t *testing.T) {
	in := "<p>start</p>" + strings.Repeat("<p>"+strings.Repeat("x", 1000)+"</p>", 110) + "<p>the end</p>"
	check(t, gmail, []expectation{{in: in, want: []string{"start", "[Message clipped]"}, wantNot: []string{"the end"}}})
	check(t, appleMail, []expectation{{in: in, want: []string{"the end"}, wantNot: []string{"[Message clipped]"}}})
}

func TestOutlookWindows(t *testing.T) {
	check(t, outlookWindows, []expectation{
		{
			in:      `<div style="width:300px;padding:10px;margin:0 auto;display:block;border-radius:4px;color:red">x</div>`,
			want:    []string{`<div style="color:red">`},
			wantNot: []string{`width`, `padding`, `margin`, `display`, `radius`},
		},
		{
			in:   `<table style="width:600px"><tr><td style="width:300px;padding:10px;background:#fff url(a.png) no-repeat">x</td></tr></table>`,
			want: []string{`<table style="width:600px">`, `style="width:300px;padding:10px;background-color:#fff"`},
		},
		{
			in:      `<p style="mso-hide:all;color:red !important;font-weight:700">x</p><img src="https://example.com/a.webp">`,
			want:    []string{`style="display:none;color:red;font-weight:bold"`, `<img>`},
			wantNot: []string{`mso-hide`, `important`},
		},
		{
			in:      `<head><style>.a.b{color:red}.c > .d{color:red}a:hover{color:red}a:link{color:blue}@media (max-width:600px){.e{color:red}}</style></head><p class="a b">x</p>`,
			want:    []string{`.a{color:red;}`, `a:link{color:blue;}`, `<p class="a">`},
			wantNot: []string{`.d{`, `:hover`, `@media`},
		},
		{
			in:      `<head><style>div.foo{color:red}a.btn{color:red}.wrapper{max-width:600px;padding:8px;color:red}.cell{padding:10px}body{padding:24px;color:#333}</style></head><body><div class="wrapper"><table><tr><td class="cell">x</td></tr></table></div></body>`,
			want:    []string{`div.foo{color:red;}`, `a.btn{color:red;}`, `.wrapper{color:red;}`, `.cell{padding:10px;}`, `body{color:#333;}`},
			wantNot: []string{`max-width`, `padding:8px`, `padding:24px`},
		},
		{
			in: `<!--[if mso]><v:roundrect href="https://example.com" style="height:40px;v-text-anchor:middle;width:200px;" arcsize="10%" stroke="f" fillcolor="#556270"><w:anchorlock/><center>Go</center></v:roundrect><![endif]-->`,
			want: []string{
				`<a href="https://example.com" target="_blank" rel="noopener noreferrer" style="`,
				`width:200px;height:40px;display:inline-flex;flex-direction:column;justify-content:center;background-color:#556270;border:0;border-radius:2.0px`,
				`<center>Go</center></a>`,
			},
		},
		{
			in:   `<!--[if gte mso 9]><v:rect style="width:600px;height:200px" stroke="false"><v:fill type="tile" src="https://example.com/bg.png" color="#7bceeb"/><v:textbox inset="0,0,0,0"><![endif]--><p>x</p><!--[if gte mso 9]></v:textbox></v:rect><![endif]-->`,
			want: []string{`background-image:url(&#34;https://example.com/bg.png&#34;);background-color:#7bceeb;border:0">`, `<div style="padding:0px 0px 0px 0px"><p>x</p></div></div>`},
		},
	})
}

func TestOutlookCom(t *testing.T) {
	check(t, outlookCom, []expectation{
		{
			in:      `<head><style>a:hover{color:red}.btn:hover{color:red}.a.b{color:red}[x]{color:red}[x^=y]{color:red}body{margin:0}@media (max-width:600px){.c{color:red}}.after{color:red}</style></head>`,
			want:    []string{`a:hover{color:red;}`, `[x]{color:red;}`, `div[data-body]{margin:0;}`, `@media (max-width:600px){.c{color:red;}}`},
			wantNot: []string{`.btn`, `.a.b`, `^=`, `.after`},
		},
		{
			in:      `<head><style>td[data-x]{color:red}div.foo{color:red}#a.b{color:red}</style></head><body style="padding:0"><p>x</p></body>`,
			want:    []string{`td[data-x]{color:red;}`, `div.foo{color:red;}`, `<div style="padding:0" data-body>`},
			wantNot: []string{`#a.b`},
		},
		{
			in:   `<p style="position:relative;margin:-5px;color:red">x</p>`,
			want: []string{`style="color:red"`},
		},
	})
}

func TestYahoo(t *testing.T) {
	check(t, yahoo, []expectation{
		{
			in:      `<head><style>@media only screen and (max-width:600px){.a{color:red}}@media (prefers-color-scheme:dark){.b{color:red}}[x=y]{color:red}[x]{color:red}body{margin:0}</style></head>`,
			want:    []string{`@media only screen and (max-width:600px)`, `[x=y]{color:red;}`},
			wantNot: []string{`prefers-color-scheme`, `[x]{`, `body{`},
		},
		{
			in:   `<td style="height:40px">x</td><img class="a" src="https://example.com/a.png" /><img class="b" src="https://example.com/b.png">`,
			want: []string{`style="min-height:40px"`, `<img src="https://example.com/a.png">`, `<img class="b"`},
		},
	})
}

func TestThunderbird(t *testing.T) {
	check(t, thunderbird, []expectation{{
		in:      `<html><head><style>@media (max-width:600px){.a{color:red}}.b{aspect-ratio:1}</style></head><body>x</body></html>`,
		want:    []string{`<!DOCTYPE html>`, `.b{aspect-ratio:1;}`, `<body>`},
		wantNot: []string{`@media`},
	}})
}

func TestSamsungEmail(t *testing.T) {
	check(t, samsungEmail, []expectation{{
		in:   `<div style="position:fixed;top:10%;font-size:80%;color:red">x</div>`,
		want: []string{`style="position:absolute;color:red"`},
	}})
}

func TestClientsAreListed(t *testing.T) {
	if Clients[0] != gmail {
		t.Errorf("default client is %q, want gmail", Clients[0].ID)
	}
	for _, client := range Clients {
		if ClientByID(client.ID) != client {
			t.Errorf("ClientByID(%q) didn't return it", client.ID)
		}
		if _, ok := caniemail[client.ID]; !ok {
			t.Errorf("no caniemail data for %q", client.ID)
		}
	}
}
