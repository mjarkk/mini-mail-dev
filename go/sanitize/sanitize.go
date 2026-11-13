package sanitize

import (
	"io"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

type action uint8

const (
	keep action = iota
	unwrap
	remove // drop the element together with its content
)

// Never rendered regardless of client: these carry script, plugins, or document
// metadata, or their content would show up as literal text.
var baselineActions = map[string]action{
	"script": remove, "template": remove, "svg": remove, "math": remove, "xml": remove,
	"iframe": remove, "noscript": remove, "noembed": remove, "noframes": remove,
	"xmp": remove, "plaintext": remove, "title": remove, "frameset": remove,
	"embed": remove, "frame": remove, "meta": remove, "base": remove, "param": remove,
	"portal": remove, "keygen": remove,
	"object": unwrap, "applet": unwrap, "canvas": unwrap,
}

// Elements a client without support removes outright instead of showing their content.
var removeWhenUnsupported = map[string]bool{
	"svg": true, "select": true, "textarea": true, "dialog": true, "input": true, "link": true,
}

var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true, "hr": true, "img": true,
	"input": true, "keygen": true, "link": true, "meta": true, "param": true, "source": true,
	"track": true, "wbr": true,
}

// headElements don't move the parser into the body.
var headElements = map[string]bool{
	"html": true, "head": true, "title": true, "meta": true, "link": true, "style": true,
	"base": true, "script": true, "noscript": true, "template": true, "xml": true,
}

// bodyMarker marks the <div> a webmail client puts in place of <body>.
const bodyMarker = "data-body"

type openElement struct {
	name string
	out  string // "" when the start tag wasn't written
}

type parser struct {
	client *Client
	root   *html.Tokenizer
	out    strings.Builder

	read       int // bytes of the original document consumed
	begun      bool
	inBody     bool
	clipped    bool
	doctype    bool
	styleBytes int // size of the <style> blocks seen so far

	open []openElement

	removing      string
	removingDepth int

	shape *vmlShape

	// One entry per open downlevel-revealed conditional comment: whether the client hides its content.
	conditions []bool

	usage map[string]map[string]bool // see Client.filterStylesheet
}

// Parse returns document the way client would display it, reduced to markup
// that is safe to render on the mini-mail-dev origin. The result is a complete
// HTML document including the client's default styles; a blank document gives "".
func Parse(document string, client *Client) (string, error) {
	document = strings.TrimSpace(document)
	if document == "" {
		return "", nil
	}

	p := &parser{client: client}
	if client.propertyElements != nil {
		p.usage = client.elementUsage(document)
	}
	p.root = html.NewTokenizer(strings.NewReader(document))
	err := p.run(p.root)
	p.flushShape()
	if p.clipped {
		p.out.WriteString(client.clipNotice)
	}
	return p.out.String(), err
}

func (p *parser) next(z *html.Tokenizer) html.TokenType {
	tt := z.Next()
	if z == p.root {
		p.read += len(z.Raw())
	}
	return tt
}

func (p *parser) run(z *html.Tokenizer) error {
	for !p.clipped {
		tt := p.next(z)
		if tt == html.ErrorToken {
			if err := z.Err(); err != io.EOF {
				return err
			}
			return nil
		}
		p.token(z, tt)
		if z == p.root && p.client.clipBytes > 0 && p.read > p.client.clipBytes {
			p.clipped = true
		}
	}
	return nil
}

func (p *parser) token(z *html.Tokenizer, tt html.TokenType) {
	if p.removing != "" {
		p.skipRemoved(z, tt)
		return
	}

	switch tt {
	case html.CommentToken:
		p.comment(string(z.Text()))
		return
	case html.DoctypeToken:
		if !p.begun && !p.client.webmail && !p.client.standardsMode {
			p.out.WriteString("<!DOCTYPE " + html.EscapeString(string(z.Text())) + ">")
		}
		p.doctype = true
		return
	}

	if p.hidden() {
		return
	}
	p.begin()

	switch tt {
	case html.TextToken:
		text := string(z.Text())
		if strings.TrimSpace(text) == "" {
			if p.inBody {
				p.out.WriteString(text)
			}
			return
		}
		p.inBody = true
		p.flushShape()
		p.out.WriteString(html.EscapeString(text))
	case html.StartTagToken, html.SelfClosingTagToken:
		p.startTag(z, tt)
	case html.EndTagToken:
		name, _ := z.TagName()
		if string(name) != "v:fill" {
			p.flushShape()
			p.endTag(string(name))
		}
	}
}

func (p *parser) begin() {
	if p.begun {
		return
	}
	p.begun = true
	if p.client.webmail || p.client.standardsMode {
		p.out.WriteString("<!DOCTYPE html>")
	}
	if p.client.baseCSS != "" {
		p.out.WriteString("<style data-client-defaults>" + p.client.baseCSS + "</style>")
	}
}

func (p *parser) skipRemoved(z *html.Tokenizer, tt html.TokenType) {
	if tt != html.StartTagToken && tt != html.EndTagToken {
		return
	}
	name, _ := z.TagName()
	if string(name) != p.removing {
		return
	}
	if tt == html.StartTagToken {
		p.removingDepth++
		return
	}
	p.removingDepth--
	if p.removingDepth == 0 {
		p.removing = ""
	}
}

func (p *parser) startTag(z *html.Tokenizer, tt html.TokenType) {
	token := z.Token()
	name := token.Data

	if name == "style" {
		p.style(z)
		return
	}
	if !headElements[name] {
		p.inBody = true
	}

	attrs := make(map[string]string, len(token.Attr))
	for _, attr := range token.Attr {
		if _, seen := attrs[attr.Key]; !seen {
			attrs[attr.Key] = strings.TrimSpace(attr.Val)
		}
	}

	// Office markup is XML, where a self-closing tag has no content or end tag.
	void := voidElements[name] || tt == html.SelfClosingTagToken && strings.Contains(name, ":")
	if p.client.vml && p.vmlStart(name, attrs, tt == html.SelfClosingTagToken) {
		return
	}
	p.flushShape()

	switch p.elementAction(name, attrs) {
	case remove:
		// The tokenizer reports "<svg/>" as self-closing, and unlike HTML elements a browser honours that for SVG and MathML.
		selfClosingForeign := tt == html.SelfClosingTagToken && (name == "svg" || name == "math")
		if !void && !selfClosingForeign {
			p.removing, p.removingDepth = name, 1
		}
		return
	case unwrap:
		if !void {
			p.open = append(p.open, openElement{name: name})
		}
		return
	}

	out := name
	marker := ""
	if name == "body" && p.client.webmail {
		out = "div"
		if p.doctype {
			p.out.WriteString(p.client.doctypeMarker)
		}
		if p.client.keepBodySelectors {
			marker = " " + bodyMarker
		}
	}
	p.out.WriteString("<" + out + p.attributes(name, token.Attr, attrs, tt == html.SelfClosingTagToken) + marker + ">")
	if !void {
		p.open = append(p.open, openElement{name: name, out: out})
	}
}

func (p *parser) endTag(name string) {
	for i := len(p.open) - 1; i >= 0; i-- {
		if p.open[i].name != name {
			continue
		}
		if out := p.open[i].out; out != "" {
			p.out.WriteString("</" + out + ">")
		}
		p.open = p.open[:i]
		return
	}

	// Browsers ignore stray end tags except these two, which create elements.
	switch name {
	case "p":
		p.out.WriteString("</p>")
	case "br":
		p.out.WriteString("<br>")
	}
}

func (p *parser) style(z *html.Tokenizer) {
	text := ""
	if p.next(z) == html.TextToken {
		text = string(z.Text())
		p.next(z)
	}

	c := p.client
	if c.unsupported("html-style") || c.styleOnlyInHead && p.inBody && !p.inside("head") {
		return
	}
	p.styleBytes += len(text)
	if c.maxStyleBytes > 0 && p.styleBytes > c.maxStyleBytes {
		return
	}
	css, ok := c.filterStylesheet(text, p.usage)
	if !ok || css == "" {
		return
	}
	p.out.WriteString("<style>" + css + "</style>")
}

func (p *parser) inside(name string) bool {
	return slices.ContainsFunc(p.open, func(e openElement) bool { return e.name == name })
}

func (p *parser) elementAction(name string, attrs map[string]string) action {
	if a, ok := baselineActions[name]; ok {
		return a
	}
	if p.client.unwrapElements[name] {
		return unwrap
	}
	if !validElementName(name) {
		// Namespaced Office markup such as <o:p> and <v:rect>: browsers render their content as if the tags weren't there.
		return unwrap
	}

	c := p.client
	switch name {
	case "html", "head":
		if c.webmail {
			return unwrap
		}
	case "link":
		if !strings.EqualFold(attrs["rel"], "stylesheet") || !safeURL(attrs["href"], false) || !isHTTP(attrs["href"]) {
			return remove
		}
		if c.webmail {
			return remove
		}
	}

	if slug := elementFeature(name, attrs); slug != "" && c.unsupported(slug) {
		if removeWhenUnsupported[name] {
			return remove
		}
		return unwrap
	}
	return keep
}

func validElementName(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
		case i > 0 && (r >= '0' && r <= '9' || r == '-' || r == '_'):
		default:
			return false
		}
	}
	return true
}

func (p *parser) attributes(tag string, list []html.Attribute, attrs map[string]string, selfClosing bool) string {
	c := p.client
	var b strings.Builder
	written := map[string]bool{}
	imageBlocked := tag == "img" && !c.imageSupported(attrs["src"])

	for _, attr := range list {
		key := attr.Key
		if written[key] {
			continue
		}
		value, ok := baselineAttribute(tag, key, strings.TrimSpace(attr.Val))
		if !ok || c.stripAttributes[key] || c.stripAttributes["data-*"] && strings.HasPrefix(key, "data-") {
			continue
		}
		if key == "class" && selfClosing && c.dropClassOnSelfClosing[tag] {
			continue
		}
		if c.webmail && tag == "body" && !c.bodyAttributes[key] {
			continue
		}
		if slug, ok := attributeFeatures[key]; ok && c.unsupported(slug) {
			continue
		}
		switch key {
		case "style":
			value = c.filterInlineStyle(tag, value)
			if value == "" {
				continue
			}
		case "src", "srcset":
			if imageBlocked {
				continue
			}
		case "class":
			if c.firstClassOnly {
				value, _, _ = strings.Cut(strings.TrimSpace(value), " ")
			}
		}
		written[key] = true
		b.WriteString(" " + key + `="` + html.EscapeString(value) + `"`)
	}

	if tag == "a" || tag == "area" {
		// The message is shown in an iframe; following a link inside it would replace the message.
		b.WriteString(` target="_blank" rel="noopener noreferrer"`)
	}
	return b.String()
}

func (c *Client) imageSupported(src string) bool {
	for _, slug := range imageFeatures(src) {
		if slug != "" && c.unsupported(slug) {
			return false
		}
	}
	return true
}

var (
	conditionalBlock = regexp.MustCompile(`(?is)^\[if ([^\]]*)\]>(.*)<!\[endif\]$`)
	conditionalStart = regexp.MustCompile(`(?is)^\[if ([^\]]*)\](><!.*)?$`)
	conditionalEnd   = regexp.MustCompile(`(?is)^(<!)?\[endif\]$`)
)

// comment handles Outlook's conditional comments; any other comment is dropped.
func (p *parser) comment(text string) {
	text = strings.TrimSpace(text)
	mso := p.client.msoVersion

	if m := conditionalBlock.FindStringSubmatch(text); m != nil {
		if mso > 0 && !p.hidden() && evalCondition(m[1], mso) {
			p.run(html.NewTokenizer(strings.NewReader(m[2])))
		}
		return
	}
	if m := conditionalStart.FindStringSubmatch(text); m != nil {
		p.conditions = append(p.conditions, mso > 0 && !evalCondition(m[1], mso))
		return
	}
	if conditionalEnd.MatchString(text) && len(p.conditions) > 0 {
		p.conditions = p.conditions[:len(p.conditions)-1]
	}
}

func (p *parser) hidden() bool {
	for _, hidden := range p.conditions {
		if hidden {
			return true
		}
	}
	return false
}

// elementUsage maps ".class" and "#id" to the elements carrying them in document,
// counting markup in conditional comments only when the client renders it.
func (c *Client) elementUsage(document string) map[string]map[string]bool {
	usage := map[string]map[string]bool{}
	add := func(key, element string) {
		if usage[key] == nil {
			usage[key] = map[string]bool{}
		}
		usage[key][element] = true
	}

	var scan func(string)
	scan = func(document string) {
		z := html.NewTokenizer(strings.NewReader(document))
		for {
			switch z.Next() {
			case html.ErrorToken:
				return
			case html.CommentToken:
				m := conditionalBlock.FindStringSubmatch(strings.TrimSpace(string(z.Text())))
				if m != nil && c.msoVersion > 0 && evalCondition(m[1], c.msoVersion) {
					scan(m[2])
				}
			case html.StartTagToken, html.SelfClosingTagToken:
				token := z.Token()
				for _, attr := range token.Attr {
					switch attr.Key {
					case "class":
						names := strings.Fields(attr.Val)
						if c.firstClassOnly && len(names) > 1 {
							names = names[:1]
						}
						for _, name := range names {
							add("."+name, token.Data)
						}
					case "id":
						add("#"+strings.TrimSpace(attr.Val), token.Data)
					}
				}
			}
		}
	}
	scan(document)
	return usage
}
