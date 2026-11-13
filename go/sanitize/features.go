package sanitize

import (
	"net/url"
	"path"
	"strings"

	"github.com/tdewolff/parse/v2/css"
)

// Properties whose caniemail slug isn't "css-" plus the property or one of its
// dash-separated prefixes; that rule would file e.g. border-top-left-radius under css-border.
var propertyFeatures = map[string]string{
	"border-top-left-radius":     "css-border-radius",
	"border-top-right-radius":    "css-border-radius",
	"border-bottom-left-radius":  "css-border-radius",
	"border-bottom-right-radius": "css-border-radius",
	"border-start-start-radius":  "css-border-radius-logical",
	"border-start-end-radius":    "css-border-radius-logical",
	"border-end-start-radius":    "css-border-radius-logical",
	"border-end-end-radius":      "css-border-radius-logical",

	"border-inline": "css-border-inline-block",
	"border-block":  "css-border-inline-block",

	"margin-inline":       "css-margin-inline-block",
	"margin-block":        "css-margin-inline-block",
	"margin-inline-start": "css-margin-inline-start-end",
	"margin-inline-end":   "css-margin-inline-start-end",
	"margin-block-start":  "css-margin-block-start-end",
	"margin-block-end":    "css-margin-block-start-end",

	"padding-inline":       "css-padding-inline-block",
	"padding-block":        "css-padding-inline-block",
	"padding-inline-start": "css-padding-inline-start-end",
	"padding-inline-end":   "css-padding-inline-start-end",
	"padding-block-start":  "css-padding-block-start-end",
	"padding-block-end":    "css-padding-block-start-end",

	"block-size": "css-block-inline-size",

	"left":   "css-left-right-top-bottom",
	"right":  "css-left-right-top-bottom",
	"top":    "css-left-right-top-bottom",
	"bottom": "css-left-right-top-bottom",

	"gap":             "css-gap",
	"row-gap":         "css-gap",
	"column-gap":      "css-gap",
	"grid-gap":        "css-gap",
	"grid-row-gap":    "css-gap",
	"grid-column-gap": "css-gap",

	"columns":           "css-column-layout-properties",
	"column-fill":       "css-column-layout-properties",
	"column-rule":       "css-column-layout-properties",
	"column-rule-color": "css-column-layout-properties",
	"column-rule-style": "css-column-layout-properties",
	"column-rule-width": "css-column-layout-properties",
	"column-span":       "css-column-layout-properties",
	"column-width":      "css-column-layout-properties",

	"flex-flow": "css-flex-wrap",
	"mask":      "css-mask-image",
}

func init() {
	for _, side := range []string{"inline", "block"} {
		for _, part := range []string{"color", "style", "width", "start", "end"} {
			propertyFeatures["border-"+side+"-"+part] = "css-border-inline-block-individual"
		}
		for _, edge := range []string{"start", "end"} {
			for _, part := range []string{"color", "style", "width"} {
				propertyFeatures["border-"+side+"-"+edge+"-"+part] = "css-border-inline-block-longhand"
			}
		}
	}
}

func unprefixed(property string) string {
	for _, prefix := range []string{"-webkit-", "-moz-", "-ms-", "-o-"} {
		if strings.HasPrefix(property, prefix) {
			return property[len(prefix):]
		}
	}
	return property
}

func propertyFeature(property string) string {
	property = unprefixed(property)
	if slug, ok := propertyFeatures[property]; ok {
		return slug
	}
	for name := property; ; {
		if caniemailFeatures["css-"+name] {
			return "css-" + name
		}
		i := strings.LastIndexByte(name, '-')
		if i <= 0 {
			return ""
		}
		name = name[:i]
	}
}

var functionFeatures = map[string]string{
	"calc":                      "css-unit-calc",
	"clamp":                     "css-function-clamp",
	"min":                       "css-function-min",
	"max":                       "css-function-max",
	"var":                       "css-variables",
	"light-dark":                "css-function-light-dark",
	"rgb":                       "css-rgb",
	"rgba":                      "css-rgba",
	"lab":                       "css-modern-color",
	"lch":                       "css-modern-color",
	"oklab":                     "css-modern-color",
	"oklch":                     "css-modern-color",
	"fit-content":               "css-intrinsic-size",
	"linear-gradient":           "css-linear-gradient",
	"repeating-linear-gradient": "css-linear-gradient",
	"radial-gradient":           "css-radial-gradient",
	"repeating-radial-gradient": "css-radial-gradient",
	"conic-gradient":            "css-conic-gradient",
	"repeating-conic-gradient":  "css-conic-gradient",
}

var keywordFeatures = map[string]string{
	"initial":       "css-unit-initial",
	"min-content":   "css-intrinsic-size",
	"max-content":   "css-intrinsic-size",
	"fit-content":   "css-intrinsic-size",
	"system-ui":     "css-sytem-ui",
	"ui-serif":      "css-sytem-ui",
	"ui-sans-serif": "css-sytem-ui",
	"ui-rounded":    "css-sytem-ui",
	"ui-monospace":  "css-sytem-ui",
}

// declarationFeatures lists every caniemail feature a declaration depends on;
// a client lacking any of them drops the declaration.
func declarationFeatures(property string, value []css.Token) []string {
	var slugs []string
	if strings.HasPrefix(property, "--") {
		return []string{"css-variables"}
	}
	if slug := propertyFeature(property); slug != "" {
		slugs = append(slugs, slug)
	}

	hasURL := false
	for i, token := range value {
		data := strings.ToLower(string(token.Data))
		switch token.TokenType {
		case css.FunctionToken:
			name := unprefixed(strings.TrimSuffix(data, "("))
			if name == "url" {
				hasURL = true
			} else if slug, ok := functionFeatures[name]; ok {
				slugs = append(slugs, slug)
			}
		case css.URLToken:
			hasURL = true
		case css.DimensionToken:
			unit := strings.TrimLeft(data, "+-.0123456789eE")
			if caniemailFeatures["css-unit-"+unit] {
				slugs = append(slugs, "css-unit-"+unit)
			}
		case css.PercentageToken:
			slugs = append(slugs, "css-unit-percent")
		case css.IdentToken:
			if slug, ok := keywordFeatures[data]; ok {
				slugs = append(slugs, slug)
			}
		case css.DelimToken:
			if data == "!" && i+1 < len(value) && strings.EqualFold(string(value[i+1].Data), "important") {
				slugs = append(slugs, "css-important")
			}
		}
	}

	switch unprefixed(property) {
	case "display":
		switch displayValue(value) {
		case "flex", "inline-flex":
			slugs = append(slugs, "css-display-flex")
		case "grid", "inline-grid":
			slugs = append(slugs, "css-display-grid")
		case "none":
			slugs = append(slugs, "css-display-none")
		}
	case "background":
		if hasURL {
			slugs = append(slugs, "css-background-image")
		}
	}
	return slugs
}

func displayValue(value []css.Token) string {
	for _, token := range value {
		if token.TokenType == css.IdentToken {
			return strings.ToLower(string(token.Data))
		}
	}
	return ""
}

func isImportant(value []css.Token) bool {
	n := len(value)
	return n >= 2 && string(value[n-2].Data) == "!" && strings.EqualFold(string(value[n-1].Data), "important")
}

var mediaFeatureSlugs = map[string]string{
	"prefers-color-scheme":   "css-at-media-prefers-color-scheme",
	"prefers-reduced-motion": "css-at-media-prefers-reduced-motion",
	"orientation":            "css-at-media-orientation",
	"hover":                  "css-at-media-hover",
	"any-hover":              "css-at-media-hover",
}

// mediaFeatureFeature maps a media feature name such as "max-width" to its caniemail slug, if any.
func mediaFeatureFeature(name string) string {
	name = unprefixed(name)
	if strings.Contains(name, "device-pixel-ratio") {
		return "css-at-media-device-pixel-ratio"
	}
	return mediaFeatureSlugs[name]
}

var elementFeatures = map[string]string{
	"ul": "html-lists", "ol": "html-lists", "li": "html-lists",
	"dl": "html-lists", "dt": "html-lists", "dd": "html-lists",

	"article": "html-semantics", "aside": "html-semantics", "details": "html-semantics",
	"figcaption": "html-semantics", "figure": "html-semantics", "footer": "html-semantics",
	"header": "html-semantics", "main": "html-semantics", "mark": "html-semantics",
	"nav": "html-semantics", "section": "html-semantics", "summary": "html-semantics",
	"time": "html-semantics",

	"h1": "html-h1-h6", "h2": "html-h1-h6", "h3": "html-h1-h6",
	"h4": "html-h1-h6", "h5": "html-h1-h6", "h6": "html-h1-h6",

	"select": "html-select", "option": "html-select", "optgroup": "html-select",
}

var inputFeatures = map[string]string{
	"checkbox": "html-input-checkbox",
	"hidden":   "html-input-hidden",
	"radio":    "html-input-radio",
	"reset":    "html-input-reset",
	"submit":   "html-input-submit",
	"text":     "html-input-text",
}

func elementFeature(tag string, attrs map[string]string) string {
	switch tag {
	case "input":
		if slug, ok := inputFeatures[strings.ToLower(attrs["type"])]; ok {
			return slug
		}
		return "html-input-text"
	case "button":
		if strings.EqualFold(attrs["type"], "reset") {
			return "html-button-reset"
		}
		return "html-button-submit"
	}
	if slug, ok := elementFeatures[tag]; ok {
		return slug
	}
	if caniemailFeatures["html-"+tag] {
		return "html-" + tag
	}
	return ""
}

var attributeFeatures = map[string]string{
	"align":       "html-align",
	"background":  "html-background",
	"cellpadding": "html-cellpadding",
	"cellspacing": "html-cellspacing",
	"dir":         "html-dir",
	"height":      "html-height",
	"width":       "html-width",
	"lang":        "html-lang",
	"role":        "html-role",
	"srcset":      "html-srcset",
	"sizes":       "html-srcset",
	"loading":     "html-loading-attribute",
	"hidden":      "html-hidden",
	"required":    "html-required",
	"inert":       "css-inert-attribute",
	"valign":      "html-valign",

	"aria-describedby": "html-aria-describedby",
	"aria-hidden":      "html-aria-hidden",
	"aria-label":       "html-aria-label",
	"aria-labelledby":  "html-aria-labelledby",
	"aria-live":        "html-aria-live",
}

var imageExtensionFeatures = map[string]string{
	".apng": "image-apng",
	".avif": "image-avif",
	".bmp":  "image-bmp",
	".gif":  "image-gif",
	".heic": "image-heif",
	".heif": "image-heif",
	".ico":  "image-ico",
	".jpg":  "image-jpg",
	".jpeg": "image-jpg",
	".mp4":  "image-mp4",
	".png":  "image-png",
	".svg":  "image-svg",
	".tif":  "image-tiff",
	".tiff": "image-tiff",
	".webp": "image-webp",
}

var imageMIMEFeatures = map[string]string{
	"image/apng":    "image-apng",
	"image/avif":    "image-avif",
	"image/bmp":     "image-bmp",
	"image/gif":     "image-gif",
	"image/heic":    "image-heif",
	"image/heif":    "image-heif",
	"image/x-icon":  "image-ico",
	"image/jpeg":    "image-jpg",
	"image/png":     "image-png",
	"image/svg+xml": "image-svg",
	"image/tiff":    "image-tiff",
	"image/webp":    "image-webp",
}

// imageFeatures guesses the image format of src from its data: MIME type or file extension.
func imageFeatures(src string) []string {
	lower := strings.ToLower(strings.TrimSpace(src))
	if rest, ok := strings.CutPrefix(lower, "data:"); ok {
		mime, _, _ := strings.Cut(rest, ";")
		mime, _, _ = strings.Cut(mime, ",")
		return []string{"image-base64", imageMIMEFeatures[mime]}
	}
	u, err := url.Parse(lower)
	if err != nil {
		return nil
	}
	return []string{imageExtensionFeatures[path.Ext(u.Path)]}
}
