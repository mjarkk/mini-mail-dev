package sanitize

import (
	"strings"

	"github.com/tdewolff/parse/v2/css"
)

// Gmail in a desktop browser with a Google account.
// https://developers.google.com/workspace/gmail/design/reference/supported_css
var gmail = &Client{
	ID:              "gmail",
	Name:            "Gmail",
	webmail:         true,
	bodyAttributes:  set("class", "id", "style", "title", "dir", "lang"),
	doctypeMarker:   "<u></u>",
	maxStyleBytes:   16 * 1024,
	styleOnlyInHead: true,
	poisons:         gmailPoisons,
	clipBytes:       102400,
	clipNotice:      `<div style="margin:16px 0;font:13px Arial,Helvetica,sans-serif;color:#222">[Message clipped]&nbsp;&nbsp;<span style="color:#15c;text-decoration:underline">View entire message</span></div>`,
	overrides: map[string]support{
		// Only [class~=x] survives, which hardly anyone writes.
		"css-selector-attribute": unsupported,
	},
	properties: set(
		"azimuth", "background", "background-blend-mode", "background-clip", "background-color",
		"background-image", "background-origin", "background-position", "background-repeat",
		"background-size", "border", "border-bottom", "border-bottom-color", "border-bottom-left-radius",
		"border-bottom-right-radius", "border-bottom-style", "border-bottom-width", "border-collapse",
		"border-color", "border-left", "border-left-color", "border-left-style", "border-left-width",
		"border-radius", "border-right", "border-right-color", "border-right-style", "border-right-width",
		"border-spacing", "border-style", "border-top", "border-top-color", "border-top-left-radius",
		"border-top-right-radius", "border-top-style", "border-top-width", "border-width", "box-sizing",
		"break-after", "break-before", "break-inside", "caption-side", "clear", "color", "column-count",
		"column-fill", "column-gap", "column-rule", "column-rule-color", "column-rule-style",
		"column-rule-width", "column-span", "column-width", "columns", "direction", "display", "elevation",
		"empty-cells", "float", "font", "font-family", "font-feature-settings", "font-kerning", "font-size",
		"font-size-adjust", "font-stretch", "font-style", "font-synthesis", "font-variant",
		"font-variant-alternates", "font-variant-caps", "font-variant-east-asian", "font-variant-ligatures",
		"font-variant-numeric", "font-weight", "height", "image-orientation", "image-resolution", "ime-mode",
		"isolation", "layout-flow", "layout-grid", "layout-grid-char", "layout-grid-char-spacing",
		"layout-grid-line", "layout-grid-mode", "layout-grid-type", "letter-spacing", "line-break",
		"line-height", "list-style", "list-style-position", "list-style-type", "margin", "margin-bottom",
		"margin-left", "margin-right", "margin-top", "marker-offset", "max-height", "max-width",
		"min-height", "min-width", "mix-blend-mode", "object-fit", "object-position", "opacity", "outline",
		"outline-color", "outline-style", "outline-width", "overflow", "overflow-x", "overflow-y", "padding",
		"padding-bottom", "padding-left", "padding-right", "padding-top", "page-break-after",
		"page-break-before", "page-break-inside", "pause", "pause-after", "pause-before", "pitch",
		"pitch-range", "quotes", "richness", "speak", "speak-header", "speak-numeral", "speak-punctuation",
		"speech-rate", "stress", "table-layout", "text-align", "text-align-last", "text-autospace",
		"text-combine-upright", "text-decoration", "text-decoration-color", "text-decoration-line",
		"text-decoration-skip", "text-decoration-style", "text-emphasis", "text-emphasis-color",
		"text-emphasis-style", "text-indent", "text-justify", "text-kashida-space", "text-orientation",
		"text-overflow", "text-transform", "text-underline-position", "unicode-bidi", "vertical-align",
		"voice-family", "white-space", "width", "word-break", "word-spacing", "word-wrap", "writing-mode",
		"zoom",
	),
	mediaTypes:         set("all", "screen"),
	mediaFeatures:      set("min-width", "max-width", "min-device-width", "max-device-width", "orientation", "min-resolution", "max-resolution"),
	unwrapElements:     set("details", "summary", "main", "nav", "picture", "source", "video", "track"),
	stripAttributes:    set("data-*", "srcset", "sizes", "loading", "hidden", "reversed", "required", "inert"),
	rewriteDeclaration: gmailDeclaration,
	// a[href] matches the specificity Gmail's ".ii a[href]" has over the email's own prefixed rules.
	baseCSS: `body{margin:0;padding:0;background:#fff;color:#222;font:small/1.5 Arial,Helvetica,sans-serif;word-wrap:break-word;overflow-wrap:break-word}a[href]{color:#15c}`,
}

func gmailDeclaration(_, property string, value []css.Token) (string, []css.Token, bool) {
	switch {
	case strings.HasPrefix(property, "margin"), property == "text-indent":
		return property, value, !hasNegativeLength(value)
	case property == "list-style":
		return property, value, !hasImage(value)
	}
	return property, value, true
}

// gmailPoisons reports values that make Gmail discard the whole style attribute or <style> block.
func gmailPoisons(value []css.Token) bool {
	depth := 0
	for i, token := range value {
		data := string(token.Data)
		switch token.TokenType {
		case css.FunctionToken:
			switch strings.ToLower(strings.TrimSuffix(data, "(")) {
			case "lab", "lch", "oklab", "oklch", "color":
				return true
			case "rgb", "rgba":
				if !commaSeparated(value[i+1:]) {
					return true
				}
			}
			depth++
		case css.LeftParenthesisToken:
			depth++
		case css.RightParenthesisToken:
			depth--
		case css.DelimToken:
			if data == "!" {
				if i+1 >= len(value) || string(value[i+1].Data) != "important" {
					return true
				}
				continue
			}
			if depth == 0 && strings.ContainsAny(data, "*$&@^?|<>\\`={}[]") {
				return true
			}
		case css.IdentToken, css.DimensionToken, css.NumberToken:
			// IE hacks such as "red\9".
			if strings.Contains(data, `\`) {
				return true
			}
		}
	}
	return false
}

// commaSeparated reports whether the function arguments up to the closing parenthesis use commas.
func commaSeparated(args []css.Token) bool {
	for _, token := range args {
		switch token.TokenType {
		case css.CommaToken:
			return true
		case css.RightParenthesisToken:
			return false
		}
	}
	return false
}

func hasNegativeLength(value []css.Token) bool {
	for _, token := range value {
		switch token.TokenType {
		case css.DimensionToken, css.NumberToken, css.PercentageToken:
			if strings.HasPrefix(string(token.Data), "-") {
				return true
			}
		}
	}
	return false
}

func hasImage(value []css.Token) bool {
	for _, token := range value {
		switch token.TokenType {
		case css.URLToken:
			return true
		case css.FunctionToken:
			name := strings.ToLower(string(token.Data))
			if name == "url(" || strings.Contains(name, "gradient(") {
				return true
			}
		}
	}
	return false
}

// Outlook.com and the new Outlook for Windows, which is the same web app.
var outlookCom = &Client{
	ID:                         "outlook-com",
	Name:                       "Outlook",
	webmail:                    true,
	bodyAttributes:             set("style", "class", "id", "dir", "lang"),
	keepBodySelectors:          true,
	pseudoOnTypeOnly:           true,
	attributeOperators:         set("", "=", "~=", "|="),
	truncateStyleAtDoubleBrace: true,
	unwrapElements:             set("article", "aside", "details", "summary", "figure", "figcaption", "footer", "header", "main", "mark", "nav", "section", "hgroup"),
	stripAttributes:            set("hidden", "loading", "srcset", "sizes", "reversed"),
	overrides: map[string]support{
		"css-pseudo-class-target": unsupported,
	},
	rewriteDeclaration: outlookComDeclaration,
	baseCSS:            `body{margin:0;padding:0;background:#fff;color:#242424;font-family:"Segoe UI","Segoe UI Web (West European)",-apple-system,BlinkMacSystemFont,Roboto,"Helvetica Neue",sans-serif;font-size:15px;word-wrap:break-word}`,
}

func outlookComDeclaration(_, property string, value []css.Token) (string, []css.Token, bool) {
	switch {
	case strings.HasPrefix(property, "margin"):
		return property, value, !hasNegativeLength(value)
	case property == "position":
		switch displayValue(value) {
		case "static", "sticky":
			return property, value, true
		}
		return "", nil, false
	case property == "gap", property == "row-gap", property == "overflow-block", property == "overflow-inline":
		return "", nil, false
	}
	return property, value, true
}

// Yahoo Mail in a desktop browser; AOL webmail shares its engine.
var yahoo = &Client{
	ID:                     "yahoo",
	Name:                   "Yahoo",
	webmail:                true,
	bodyAttributes:         set("style", "class", "id", "dir", "lang"),
	mediaTypes:             set("screen"),
	mediaFeatures:          set("min-width", "max-width", "min-height", "max-height"),
	attributeOperators:     set("="),
	stripAttributes:        set("aria-label", "aria-hidden", "aria-describedby", "aria-labelledby", "aria-live", "hidden", "loading", "srcset", "sizes", "reversed"),
	dropClassOnSelfClosing: set("img", "hr"),
	rewriteDeclaration:     yahooDeclaration,
	baseCSS:                `body{margin:0;padding:0;background:#fff;color:#1d2228;font-family:"Helvetica Neue",Helvetica,Arial,sans-serif;font-size:13px}`,
}

func yahooDeclaration(element, property string, value []css.Token) (string, []css.Token, bool) {
	switch {
	case property == "height":
		if element != "" && displayValue(value) != "auto" {
			return "min-height", value, true
		}
	case strings.HasPrefix(property, "margin"), property == "text-indent":
		return property, value, !hasNegativeLength(value)
	case property == "position":
		switch displayValue(value) {
		case "static", "relative":
			return property, value, true
		}
		return "", nil, false
	case property == "display":
		switch displayValue(value) {
		case "inline-flex", "inline-grid", "flow-root", "contents", "revert":
			return "", nil, false
		}
	case strings.HasPrefix(property, "transition-"):
		return "", nil, false
	case property == "font-weight":
		if len(value) == 1 && value[0].TokenType == css.NumberToken && !strings.HasSuffix(string(value[0].Data), "00") {
			return "", nil, false
		}
	}
	return property, value, true
}

// Apple Mail on macOS renders the message as its own document with WebKit.
var appleMail = &Client{
	ID:                 "apple-mail",
	Name:               "Apple Mail",
	keepBodySelectors:  true,
	keepRootSelectors:  true,
	rewriteDeclaration: webViewDeclaration,
	// From Mail's MUIWebDocument.css; :where keeps the link colour below the email's own rules.
	baseCSS: `html{font:12px Helvetica}body{margin:18px 20px 22px;padding:0;word-wrap:break-word;word-break:break-word}table{word-break:normal}:where(:link){color:#0068DA}`,
}

// Samsung Email on Android, with a non-Microsoft account.
var samsungEmail = &Client{
	ID:                 "samsung-email",
	Name:               "Samsung",
	Width:              360,
	keepBodySelectors:  true,
	keepRootSelectors:  true,
	rewriteDeclaration: samsungDeclaration,
	baseCSS:            `html{font-family:Roboto,SamsungOne,sans-serif;font-size:16px;color:#000}body{margin:8px;word-wrap:break-word}:where(:link){color:#0000EE}`,
}

// Thunderbird with remote content allowed for the sender.
var thunderbird = &Client{
	ID:                "thunderbird",
	Name:              "Thunderbird",
	standardsMode:     true,
	keepBodySelectors: true,
	keepRootSelectors: true,
	// caniemail's results predate Gecko support that Thunderbird's sanitizer leaves alone.
	overrides: map[string]support{
		"css-aspect-ratio":     supported,
		"css-modern-color":     supported,
		"css-pseudo-class-has": supported,
	},
	baseCSS: `html{font-family:"Lucida Grande",sans-serif;font-size:15px;color:#000}body{margin:8px;line-height:initial!important}code,kbd,pre,samp,tt{font-family:Menlo,Monaco,monospace;font-size:12px}:where(:link){color:#0B6CDA}img{max-width:calc(100vw - 16px)!important}`,
}

// webViewDeclaration covers what Apple's and Samsung's WebKit/Chromium views change.
func webViewDeclaration(_, property string, value []css.Token) (string, []css.Token, bool) {
	switch property {
	case "position":
		switch displayValue(value) {
		case "fixed":
			return property, []css.Token{{TokenType: css.IdentToken, Data: []byte("absolute")}}, true
		case "sticky":
			return property, []css.Token{{TokenType: css.IdentToken, Data: []byte("static")}}, true
		}
	case "overflow-block", "overflow-inline":
		return "", nil, false
	}
	return property, value, true
}

func samsungDeclaration(element, property string, value []css.Token) (string, []css.Token, bool) {
	switch property {
	case "font-size":
		switch displayValue(value) {
		case "larger", "smaller":
			return "", nil, false
		}
		if len(value) > 0 && value[0].TokenType == css.PercentageToken {
			return "", nil, false
		}
	case "top", "right", "bottom", "left":
		for _, token := range value {
			if token.TokenType == css.PercentageToken {
				return "", nil, false
			}
		}
	case "text-align":
		if displayValue(value) == "match-parent" {
			return "", nil, false
		}
	case "border-inline-color", "border-inline-width", "border-block-color", "border-block-width":
		return "", nil, false
	}
	return webViewDeclaration(element, property, value)
}
