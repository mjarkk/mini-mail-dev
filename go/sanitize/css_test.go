package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCssSelectors(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		// Selectors
		{"element", "a { color: red; }", "a{color:red;}"},
		{"class", ".foo { color: red; }", ".foo{color:red;}"},
		{"id", "#foo { color: red; }", "#foo{color:red;}"},
		{"all", "* { color: red; }", "*{color:red;}"},
		{"attribute", "[foo] { color: red; }", "[foo]{color:red;}"},
		{"attribute value", "[foo=value] {}", "[foo=value]{}"},
		{"attribute value double qoutes", "[foo=\"value\"] {}", "[foo=\"value\"]{}"},
		{"attribute value single qoutes", "[foo='value'] {}", "[foo='value']{}"},
		{"attribute value ends with", "[foo$=value] {}", "[foo$=value]{}"},
		{"attribute value starts with hyphen", "[foo|=value] {}", "[foo|=value]{}"},
		{"attribute value starts with", "[foo^=value] {}", "[foo^=value]{}"},
		{"attribute value contains with rules", "[foo~=value] { }", "[foo~=value]{}"},
		{"attribute value contains", "[foo*=value] { color: red; }", "[foo*=value]{color:red;}"},
		{"combinator", "a b {}", "a b{}"},
		{"combinator", "a, b {}", "a, b{}"},
		{"combinator child", "a > b {}", "a > b{}"},
		{"combinator next", "a + b {}", "a + b{}"},
		{"combinator subsquent sibling", "a ~ b {}", "a ~ b{}"},
		{"pseudo-class", "a:hover { color: red; }", "a:hover{color:red;}"},
		{"double Pseudo-class", "a::before {}", "a::before{}"},
		{"chain", "a.foo#bar:hover::before {}", "a.foo#bar:hover::before{}"},
		{"multiple", "a { color: red; } b { color: blue; }", "a{color:red;}b{color:blue;}"},

		// Bogus selectors
		{"bogus multiple class dots", "..a { }", "..a{}"},
		{"bogus multiple hashtags", "##a { }", "#a{}"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			output, err := parseCss(test.input)
			assert.NoError(t, err, test.name)
			assert.Equal(t, test.output, output, test.name)
		})
	}
}

func TestParseCssProperties(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"color", "color: red;", "color:red;"},
		{"zero width", "width: 0;", "width:0;"},
		{"calc", "width: calc(100vw - 100px);", "width:calc(100vw-100px);"},
		{"background-image", "url(https://example.com)", "url(https://example.com)"},
		{"background-image 2", "url(\"https://example.com\")", "url(\"https://example.com\")"},
		{"background-image 3", "url('https://example.com')", "url('https://example.com')"},
		{"@import", "@import url(\"custom.css\") print;", ""},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			output, err := parseCss("a{" + test.input + "}")
			assert.NoError(t, err, test.name)
			assert.Equal(t, "a{"+test.output+"}", output, test.name)
		})
	}
}
