package sanitize

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
	"golang.org/x/net/html"
)

// Classic Outlook for Windows (2016 and later, Microsoft 365) renders with Microsoft Word.
// https://learn.microsoft.com/en-us/previous-versions/office/developer/office-2007/aa338201(v=office.12)
var outlookWindows = &Client{
	ID:         "outlook-windows",
	Name:       "Outlook (desktop)",
	msoVersion: 16,
	vml:        true,
	// Word ignores web fonts; caniemail's @import support is about the stylesheet, not the fonts it loads.
	blockWebFonts:     true,
	keepBodySelectors: true,
	keepRootSelectors: true,
	chainedSelectors:  chainFirstOnly,
	firstClassOnly:    true,
	overrides: map[string]support{
		"css-at-font-face": unsupported,
		"css-text-justify": unsupported,
	},
	properties: set(
		"color", "font", "font-family", "font-style", "font-variant", "font-size", "font-weight",
		"text-decoration", "text-transform", "text-align", "text-indent", "vertical-align",
		"letter-spacing", "line-height", "background", "background-color",
		"border", "border-color", "border-style", "border-width",
		"border-top", "border-right", "border-bottom", "border-left",
		"border-top-color", "border-top-style", "border-top-width",
		"border-right-color", "border-right-style", "border-right-width",
		"border-bottom-color", "border-bottom-style", "border-bottom-width",
		"border-left-color", "border-left-style", "border-left-width", "border-collapse",
		"margin", "margin-top", "margin-right", "margin-bottom", "margin-left",
		"padding", "padding-top", "padding-right", "padding-bottom", "padding-left",
		"width", "height", "max-width", "display", "list-style", "list-style-type",
		"direction", "word-break",
	),
	propertyElements: outlookPropertyElements(),
	stripAttributes: set(
		"title", "disabled", "checked", "selected", "readonly", "multiple", "srcset", "sizes",
		"loading", "inert", "required", "reversed", "hidden",
	),
	rewriteDeclaration: outlookDeclaration,
	// Word's defaults for HTML without styles; Calibri/Aptos is only the compose font.
	baseCSS: `body{margin:0;padding:12px;font-family:'Times New Roman',Times,serif;font-size:12pt;color:#000000}a{color:#0563C1;text-decoration:underline}img{border:0;vertical-align:baseline}ul,ol{padding-left:0;margin-left:36pt}`,
}

func outlookPropertyElements() map[string][]string {
	cells := []string{"td", "th"}
	blocks := []string{"p", "div", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "li", "dl", "dd", "dt", "blockquote", "pre", "center", "address", "hr", "table"}
	elements := map[string][]string{
		"width":     {"table", "td", "th"},
		"height":    {"table", "tr", "td", "th"},
		"max-width": {"table"},
	}
	for _, side := range []string{"", "-top", "-right", "-bottom", "-left"} {
		elements["padding"+side] = cells
		elements["margin"+side] = blocks
	}
	return elements
}

var hiddenDisplay = []css.Token{{TokenType: css.IdentToken, Data: []byte("none")}}

func outlookDeclaration(element, property string, value []css.Token) (string, []css.Token, bool) {
	if element != "" {
		// Word ignores !important in style attributes but keeps the declaration.
		value = withoutImportant(value)
	}

	switch property {
	case "mso-hide":
		if strings.EqualFold(displayValue(value), "all") {
			return "display", hiddenDisplay, true
		}
		return "", nil, false
	case "display":
		return property, value, displayValue(value) == "none" && element != "img"
	case "background":
		if color := colorToken(value); color != nil {
			return "background-color", color, true
		}
		return "", nil, false
	case "list-style":
		for _, token := range value {
			if token.TokenType == css.IdentToken && !listStylePositions[strings.ToLower(string(token.Data))] {
				return "list-style-type", []css.Token{token}, true
			}
		}
		return "", nil, false
	case "font-weight":
		if len(value) == 1 && value[0].TokenType == css.NumberToken {
			weight, _ := strconv.Atoi(string(value[0].Data))
			keyword := "normal"
			if weight >= 600 {
				keyword = "bold"
			}
			return property, []css.Token{{TokenType: css.IdentToken, Data: []byte(keyword)}}, true
		}
	case "text-transform":
		return property, value, displayValue(value) != "lowercase"
	case "text-align":
		switch displayValue(value) {
		case "start", "end", "match-parent":
			return "", nil, false
		}
	case "text-decoration":
		return property, value, len(value) == 1 && displayValue(value) != "overline"
	}

	if strings.HasPrefix(property, "margin") {
		for _, token := range value {
			if strings.EqualFold(string(token.Data), "auto") || strings.HasPrefix(string(token.Data), "-") {
				return "", nil, false
			}
		}
	}
	return property, value, true
}

var listStylePositions = set("inside", "outside", "none")

func withoutImportant(value []css.Token) []css.Token {
	if !isImportant(value) {
		return value
	}
	return trimWhitespace(value[:len(value)-2])
}

// colorToken picks the color out of a background shorthand.
func colorToken(value []css.Token) []css.Token {
	for i, token := range value {
		switch token.TokenType {
		case css.HashToken:
			return value[i : i+1]
		case css.FunctionToken:
			if strings.EqualFold(string(token.Data), "rgb(") {
				for j := i; j < len(value); j++ {
					if value[j].TokenType == css.RightParenthesisToken {
						return value[i : j+1]
					}
				}
			}
		case css.IdentToken:
			if !backgroundKeywords[strings.ToLower(string(token.Data))] {
				return value[i : i+1]
			}
		}
	}
	return nil
}

var backgroundKeywords = set(
	"none", "repeat", "repeat-x", "repeat-y", "no-repeat", "space", "round", "top", "bottom", "left", "right",
	"center", "fixed", "scroll", "local", "border-box", "padding-box", "content-box", "text", "cover", "contain", "auto",
	"inherit", "initial", "unset",
)

// vmlShape is a VML shape whose HTML is written once its <v:fill> child, if any, has been seen.
type vmlShape struct {
	name  string
	tag   string
	attrs map[string]string
	fill  map[string]string
}

var vmlShapes = set("v:rect", "v:roundrect", "v:oval")

// vmlStart handles a VML start tag and reports whether it was one.
func (p *parser) vmlStart(name string, attrs map[string]string, selfClosing bool) bool {
	switch {
	case vmlShapes[name]:
		p.flushShape()
		tag := "div"
		if href := attrs["href"]; href != "" && safeURL(href, false) {
			tag = "a"
		}
		p.shape = &vmlShape{name: name, tag: tag, attrs: attrs}
		if selfClosing {
			p.flushShape()
			p.out.WriteString("</" + tag + ">")
		}
		return true
	case name == "v:fill":
		if p.shape != nil {
			p.shape.fill = attrs
		}
		return true
	case name == "v:textbox":
		p.flushShape()
		inset := vmlInset(attrs["inset"])
		p.out.WriteString(`<div style="` + html.EscapeString("padding:"+inset) + `">`)
		if !selfClosing {
			p.open = append(p.open, openElement{name: name, out: "div"})
		}
		return true
	}
	return false
}

// flushShape writes the start tag of a pending VML shape.
func (p *parser) flushShape() {
	s := p.shape
	if s == nil {
		return
	}
	p.shape = nil

	style := vmlStyle(s.attrs["style"])
	width, height := style["width"], style["height"]
	decls := []string{"box-sizing:border-box", "text-decoration:none", "color:inherit", "vertical-align:top"}
	if width != "" {
		decls = append(decls, "width:"+width)
	}
	if height != "" {
		decls = append(decls, "height:"+height)
	}

	// VML anchors text at the top unless v-text-anchor says otherwise.
	justify := "flex-start"
	switch style["v-text-anchor"] {
	case "middle", "middle-center", "middle-baseline":
		justify = "center"
	case "bottom", "bottom-center", "bottom-baseline":
		justify = "flex-end"
	}
	decls = append(decls, "display:inline-flex", "flex-direction:column", "justify-content:"+justify)

	fill := "#ffffff"
	if c := vmlColor(s.attrs["fillcolor"]); c != "" {
		fill = c
	}
	if vmlFalse(s.attrs["filled"]) || vmlFalse(s.attrs["fill"]) {
		fill = "transparent"
	}
	if s.fill != nil {
		if c := vmlColor(s.fill["color"]); c != "" {
			fill = c
		}
		if src := s.fill["src"]; isHTTP(src) && safeURL(src, false) {
			decls = append(decls, `background-image:url("`+cssString(src)+`")`)
			if strings.EqualFold(s.fill["type"], "frame") {
				decls = append(decls, "background-size:100% 100%")
			}
		}
	}
	decls = append(decls, "background-color:"+fill)

	if vmlFalse(s.attrs["stroke"]) || vmlFalse(s.attrs["stroked"]) {
		decls = append(decls, "border:0")
	} else {
		color := "#000000"
		if c := vmlColor(s.attrs["strokecolor"]); c != "" {
			color = c
		}
		weight := "0.75pt"
		if w := vmlLength(s.attrs["strokeweight"]); w != "" {
			weight = w
		}
		decls = append(decls, "border:"+weight+" solid "+color)
	}

	switch s.name {
	case "v:oval":
		decls = append(decls, "border-radius:50%")
	case "v:roundrect":
		if radius := vmlRadius(s.attrs["arcsize"], width, height); radius != "" {
			decls = append(decls, "border-radius:"+radius)
		}
	}

	attrs := ` style="` + html.EscapeString(strings.Join(decls, ";")) + `"`
	if s.tag == "a" {
		attrs = ` href="` + html.EscapeString(s.attrs["href"]) + `" target="_blank" rel="noopener noreferrer"` + attrs
	}
	p.out.WriteString("<" + s.tag + attrs + ">")
	p.open = append(p.open, openElement{name: s.name, out: s.tag})
}

func vmlStyle(style string) map[string]string {
	result := map[string]string{}
	p := css.NewParser(parse.NewInputString(style), true)
	for {
		gt, tt, data := p.Next()
		if gt == css.ErrorGrammar {
			if tt == css.ErrorToken || !p.HasParseError() {
				return result
			}
			continue
		}
		if gt != css.DeclarationGrammar {
			continue
		}
		name := string(data)
		values := p.Values()
		switch name {
		case "width", "height":
			if len(values) == 1 {
				result[name] = vmlLength(string(values[0].Data))
			}
		case "v-text-anchor":
			result[name] = strings.ToLower(tokensString(values))
		}
	}
}

var (
	vmlLengthPattern = regexp.MustCompile(`^(\d+(?:\.\d+)?)(px|pt|in|cm|mm|%)?$`)
	vmlColorPattern  = regexp.MustCompile(`^(#[0-9a-fA-F]{3}|#[0-9a-fA-F]{6}|[a-zA-Z]+)$`)
)

// vmlLength returns a CSS length, or "" when value isn't a plain length. Unitless VML lengths are pixels.
func vmlLength(value string) string {
	m := vmlLengthPattern.FindStringSubmatch(strings.TrimSpace(strings.ToLower(value)))
	if m == nil {
		return ""
	}
	if m[2] == "" {
		return m[1] + "px"
	}
	return m[1] + m[2]
}

func vmlColor(value string) string {
	value = strings.TrimSpace(value)
	if vmlColorPattern.MatchString(value) {
		return value
	}
	return ""
}

func vmlFalse(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "f", "false":
		return true
	}
	return false
}

// vmlInset converts a v:textbox inset ("left,top,right,bottom") to CSS padding.
func vmlInset(inset string) string {
	sides := []string{"0.1in", "0.05in", "0.1in", "0.05in"}
	for i, part := range strings.Split(inset, ",") {
		if i < len(sides) {
			if length := vmlLength(part); length != "" {
				sides[i] = length
			}
		}
	}
	return sides[1] + " " + sides[2] + " " + sides[3] + " " + sides[0]
}

// vmlRadius converts a roundrect arcsize, a fraction of half the shorter side, to a border-radius.
func vmlRadius(arcsize, width, height string) string {
	fraction := 0.2
	if arcsize = strings.TrimSpace(arcsize); arcsize != "" {
		if f, err := strconv.ParseFloat(strings.TrimSuffix(arcsize, "%"), 64); err == nil {
			fraction = f
			if strings.HasSuffix(arcsize, "%") {
				fraction = f / 100
			} else if f > 1 {
				fraction = f / 65536 // fixed-point fraction
			}
		}
	}
	shortest := 0.0
	for _, length := range []string{width, height} {
		if v, err := strconv.ParseFloat(strings.TrimSuffix(length, "px"), 64); err == nil && strings.HasSuffix(length, "px") {
			if shortest == 0 || v < shortest {
				shortest = v
			}
		}
	}
	if shortest == 0 {
		return ""
	}
	return fmt.Sprintf("%.1fpx", fraction*shortest/2)
}

func cssString(s string) string {
	return strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", "", "\r", "").Replace(s)
}
