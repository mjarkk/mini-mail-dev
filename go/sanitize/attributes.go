package sanitize

import (
	"strconv"
	"strings"
)

// Presentational attributes email HTML relies on; none of them can run script or load a document.
var globalAttributes = set(
	"class", "id", "style", "title", "dir", "lang", "role", "hidden", "inert",
	"align", "valign", "width", "height", "bgcolor", "background", "border", "nowrap",
)

var elementAttributes = map[string]map[string]bool{
	"a":        set("href", "name"),
	"area":     set("href", "shape", "coords", "alt"),
	"map":      set("name"),
	"img":      set("src", "alt", "srcset", "sizes", "usemap", "hspace", "vspace", "loading"),
	"table":    set("cellpadding", "cellspacing", "frame", "rules", "summary"),
	"td":       set("colspan", "rowspan", "headers", "scope", "abbr"),
	"th":       set("colspan", "rowspan", "headers", "scope", "abbr"),
	"col":      set("span"),
	"colgroup": set("span"),
	"font":     set("face", "size", "color"),
	"ol":       set("type", "start", "reversed"),
	"ul":       set("type"),
	"li":       set("type", "value"),
	"body":     set("text", "link", "vlink", "alink"),
	"hr":       set("size", "noshade", "color"),
	"input":    set("type", "name", "value", "checked", "disabled", "placeholder", "readonly", "maxlength", "size", "min", "max", "step", "required"),
	"button":   set("type", "name", "value", "disabled"),
	"select":   set("name", "multiple", "disabled", "size", "required"),
	"option":   set("value", "selected", "disabled", "label"),
	"optgroup": set("label", "disabled"),
	"textarea": set("name", "rows", "cols", "placeholder", "disabled", "readonly", "required", "maxlength"),
	"label":    set("for"),
	"form":     set("name"),
	"fieldset": set("name", "disabled"),
	"video":    set("src", "poster", "controls", "autoplay", "muted", "loop", "playsinline", "preload"),
	"audio":    set("src", "controls", "autoplay", "muted", "loop", "preload"),
	"source":   set("src", "srcset", "sizes", "type", "media"),
	"details":  set("open"),
	"dialog":   set("open"),
	"meter":    set("value", "min", "max", "low", "high", "optimum"),
	"progress": set("value", "max"),
	"marquee":  set("behavior", "direction", "scrollamount", "scrolldelay", "loop"),
	"link":     set("rel", "href", "media", "type"),
	"time":     set("datetime"),
	"del":      set("datetime"),
	"ins":      set("datetime"),
}

var urlAttributes = set("href", "src", "poster", "background")

// baselineAttribute decides whether an attribute may be rendered at all, whatever the client.
func baselineAttribute(tag, key, value string) (string, bool) {
	if !globalAttributes[key] && !elementAttributes[tag][key] && !dataOrARIA(key) {
		return "", false
	}
	switch {
	case urlAttributes[key]:
		return value, safeURL(value, key != "href")
	case key == "srcset":
		for _, candidate := range strings.Split(value, ",") {
			url, _, _ := strings.Cut(strings.TrimSpace(candidate), " ")
			if !safeURL(url, true) {
				return "", false
			}
		}
	}
	return value, true
}

func dataOrARIA(key string) bool {
	name, ok := strings.CutPrefix(key, "data-")
	if !ok {
		name, ok = strings.CutPrefix(key, "aria-")
	}
	if !ok || name == "" {
		return false
	}
	for _, r := range name {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.') {
			return false
		}
	}
	return true
}

// safeURL reports whether a URL can only navigate to or load ordinary content.
// dataImages permits data:image/ URLs, which are inert outside of navigation.
func safeURL(raw string, dataImages bool) bool {
	// Browsers drop tabs and newlines anywhere in a URL and trim leading control
	// characters and spaces, so "java\tscript:" still runs script.
	cleaned := strings.Map(func(r rune) rune {
		if r == '\t' || r == '\n' || r == '\r' {
			return -1
		}
		return r
	}, raw)
	cleaned = strings.TrimLeft(cleaned, "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x0b\x0c\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f ")

	i := strings.IndexAny(cleaned, ":/?#")
	if i <= 0 || cleaned[i] != ':' {
		return true
	}
	switch strings.ToLower(cleaned[:i]) {
	case "http", "https", "mailto", "tel", "sms", "cid":
		return true
	case "data":
		return dataImages && strings.HasPrefix(strings.ToLower(cleaned[i+1:]), "image/")
	}
	return false
}

func isHTTP(raw string) bool {
	lower := strings.ToLower(strings.TrimSpace(raw))
	return strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "http://")
}

// evalCondition evaluates a conditional comment expression such as "gte mso 9"
// or "(mso)|(IE)" the way Outlook running Office version mso does.
func evalCondition(expr string, mso int) bool {
	e := conditionEval{tokens: conditionTokens(expr), mso: float64(mso)}
	return e.or()
}

func conditionTokens(expr string) []string {
	var tokens []string
	word := strings.Builder{}
	flush := func() {
		if word.Len() > 0 {
			tokens = append(tokens, strings.ToLower(word.String()))
			word.Reset()
		}
	}
	for _, r := range expr {
		switch r {
		case '(', ')', '!', '&', '|':
			flush()
			tokens = append(tokens, string(r))
		case ' ', '\t', '\n', '\r':
			flush()
		default:
			word.WriteRune(r)
		}
	}
	flush()
	return tokens
}

type conditionEval struct {
	tokens []string
	mso    float64
}

func (e *conditionEval) peek() string {
	if len(e.tokens) == 0 {
		return ""
	}
	return e.tokens[0]
}

func (e *conditionEval) take() string {
	t := e.peek()
	if len(e.tokens) > 0 {
		e.tokens = e.tokens[1:]
	}
	return t
}

func (e *conditionEval) or() bool {
	result := e.and()
	for e.peek() == "|" {
		e.take()
		result = e.and() || result
	}
	return result
}

func (e *conditionEval) and() bool {
	result := e.unary()
	for e.peek() == "&" {
		e.take()
		result = e.unary() && result
	}
	return result
}

func (e *conditionEval) unary() bool {
	switch e.peek() {
	case "!":
		e.take()
		return !e.unary()
	case "(":
		e.take()
		result := e.or()
		if e.peek() == ")" {
			e.take()
		}
		return result
	}
	return e.comparison()
}

func (e *conditionEval) comparison() bool {
	operator := ""
	switch e.peek() {
	case "lt", "lte", "gt", "gte":
		operator = e.take()
	}

	var have float64
	switch e.take() {
	case "mso":
		have = e.mso
	case "vml":
		return true
	default:
		// IE and anything else Outlook isn't.
		if _, err := strconv.ParseFloat(e.peek(), 64); err == nil {
			e.take()
		}
		return false
	}

	want, err := strconv.ParseFloat(e.peek(), 64)
	if err != nil {
		return true
	}
	e.take()
	switch operator {
	case "lt":
		return have < want
	case "lte":
		return have <= want
	case "gt":
		return have > want
	case "gte":
		return have >= want
	}
	return have == want
}
