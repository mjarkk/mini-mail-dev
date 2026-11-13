package sanitize

import "testing"

func TestParseText(t *testing.T) {
	cases := []struct{ in, want string }{
		{`<p>foo <span>bar</span></p>`, `<p>foo <span>bar</span></p>`},
		{`<p><span>foo</span> <span>bar</span></p>`, `<p><span>foo</span> <span>bar</span></p>`},
		{"<p><b>foo</b>\n<i>bar</i></p>", "<p><b>foo</b>\n<i>bar</i></p>"},
		{"<pre>  a\n  b</pre>", "<pre>  a\n  b</pre>"},
		{`<p>foo&nbsp;bar</p>`, "<p>foo bar</p>"},
		{`<p>a &amp; b &lt;script&gt;</p>`, `<p>a &amp; b &lt;script&gt;</p>`},
	}

	for _, c := range cases {
		got, err := Parse(c.in, &Client{})
		if err != nil {
			t.Fatalf("Parse(%q) error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("Parse(%q)\n got: %q\nwant: %q", c.in, got, c.want)
		}
	}
}

func mustParse(t *testing.T, in string, client *Client) string {
	t.Helper()
	got, err := Parse(in, client)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", in, err)
	}
	return got
}

func TestParseRemovesScript(t *testing.T) {
	cases := []struct{ in, want string }{
		{`<a href="javascript:alert(1)">x</a>`, `<a target="_blank" rel="noopener noreferrer">x</a>`},
		{`<a href="JaVaScRiPt:alert(1)">x</a>`, `<a target="_blank" rel="noopener noreferrer">x</a>`},
		{"<a href=\"java\tscript:alert(1)\">x</a>", `<a target="_blank" rel="noopener noreferrer">x</a>`},
		{`<a href=" &#106;avascript:alert(1)">x</a>`, `<a target="_blank" rel="noopener noreferrer">x</a>`},
		{`<a href="data:text/html,<script>alert(1)</script>">x</a>`, `<a target="_blank" rel="noopener noreferrer">x</a>`},
		{`<a href="https://example.com" onclick="alert(1)">x</a>`, `<a href="https://example.com" target="_blank" rel="noopener noreferrer">x</a>`},
		{`<img src="data:image/png;base64,AAAA" onerror="alert(1)">`, `<img src="data:image/png;base64,AAAA">`},
		{`<p>a<script>alert(1)</script>b</p>`, `<p>ab</p>`},
		{`<p>a<svg><script>alert(1)</script><text>t</text></svg>b</p>`, `<p>ab</p>`},
		{`<p>a<svg/>b</p>`, `<p>ab</p>`},
		{`<p>a<noscript><img src=x onerror=alert(1)></noscript>b</p>`, `<p>ab</p>`},
		{`<p>a<iframe src="https://example.com"></iframe>b</p>`, `<p>ab</p>`},
		{`<form action="https://example.com"><input type="text" name="q" formaction="https://example.com"></form>`, `<form><input type="text" name="q"></form>`},
		{`<style>p{content:"<"}</style>`, `<style>p{content:"\3c ";}</style>`},
		{`<style>p{a:< /style><script>alert(1)</script>}</style>`, `<style>p{a:\3c /style>\3c script>alert(1)\3c /script>;}</style>`},
	}

	for _, c := range cases {
		if got := mustParse(t, c.in, &Client{}); got != c.want {
			t.Errorf("Parse(%q)\n got: %q\nwant: %q", c.in, got, c.want)
		}
	}
}

func TestParseConditionalComments(t *testing.T) {
	in := `<p>a</p><!--[if mso]><p>outlook</p><![endif]--><!--[if !mso]><!--><p>others</p><!--<![endif]--><!--[if gte mso 9]><p>9+</p><![endif]--><!--[if lt mso 12]><p>old</p><![endif]--><!--[if (mso)|(IE)]><p>either</p><![endif]-->`

	if got, want := mustParse(t, in, &Client{}), `<p>a</p><p>others</p>`; got != want {
		t.Errorf("browser\n got: %q\nwant: %q", got, want)
	}
	if got, want := mustParse(t, in, &Client{msoVersion: 16}), `<p>a</p><p>outlook</p><p>9+</p><p>either</p>`; got != want {
		t.Errorf("outlook\n got: %q\nwant: %q", got, want)
	}
}

func TestParseOfficeMarkup(t *testing.T) {
	in := `<!--[if gte mso 9]><xml><o:OfficeDocumentSettings><o:PixelsPerInch>96</o:PixelsPerInch></o:OfficeDocumentSettings></xml><![endif]--><p>a<o:p>&nbsp;</o:p></p>`
	if got, want := mustParse(t, in, &Client{msoVersion: 16}), "<p>a\u00a0</p>"; got != want {
		t.Errorf("got: %q\nwant: %q", got, want)
	}
}

func TestParseStrayEndTags(t *testing.T) {
	in := `<div>a</span></div></p></br>`
	if got, want := mustParse(t, in, &Client{}), `<div>a</div></p><br>`; got != want {
		t.Errorf("got: %q\nwant: %q", got, want)
	}
}

func FuzzParse(f *testing.F) {
	f.Add(`<!DOCTYPE html><html><head><style>@media (max-width:600px){.a{color:red}}.b:hover{x:y}</style></head><body style="margin:0"><p class="a b">x</p></body></html>`)
	f.Add(`<!--[if mso]><v:roundrect style="width:200px;height:40px" arcsize="10%" fillcolor="#556270"><v:fill src="https://e.com/a.png"/><v:textbox inset="1,2,3,4"><center>Go</center></v:textbox></v:roundrect><![endif]--><!--[if !mso]><!--><a href="https://e.com">x</a><!--<![endif]-->`)
	f.Add(`<style>@font-face{src:url(a)}a{b:c!IMPORTANT;color:rgb(0 0 0);--x:y}[x^=y]{}</style><svg><style>p{}</style></svg><img src="data:image/gif;base64,AA" />`)
	f.Fuzz(func(t *testing.T, document string) {
		for _, client := range Clients {
			Parse(document, client)
		}
	})
}
