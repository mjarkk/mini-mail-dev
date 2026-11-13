package sanitize

import (
	"bytes"
	"slices"
	"strings"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

func tokensString(tokens []css.Token) string {
	var b strings.Builder
	for _, token := range tokens {
		b.Write(token.Data)
	}
	return b.String()
}

// filterInlineStyle returns the declarations of a style attribute on element that the client keeps.
func (c *Client) filterInlineStyle(element, style string) string {
	p := css.NewParser(parse.NewInputString(style), true)
	var kept []string
	for {
		gt, tt, data := p.Next()
		if gt == css.ErrorGrammar {
			if tt == css.ErrorToken || !p.HasParseError() {
				break
			}
			continue
		}
		if gt != css.DeclarationGrammar && gt != css.CustomPropertyGrammar {
			continue
		}
		if c.poisons != nil && c.poisons(p.Values()) {
			return ""
		}
		if declaration, ok := c.declaration(element, []string{element}, string(data), p.Values()); ok {
			kept = append(kept, declaration)
		}
	}
	return strings.Join(kept, ";")
}

// filterStylesheet returns the part of a <style> block the client keeps,
// or ok=false when the client discards the whole block. usage maps ".class"
// and "#id" to the elements in the message that carry them.
func (c *Client) filterStylesheet(text string, usage map[string]map[string]bool) (string, bool) {
	if c.truncateStyleAtDoubleBrace {
		if i := strings.Index(text, "}}"); i >= 0 {
			text = text[:i+2]
		}
	}

	p := css.NewParser(parse.NewInputString(text), false)
	var out bytes.Buffer
	skipping := 0
	type block struct {
		name         string   // at-rule name, "" for a ruleset
		start, inner int      // offsets in out of the block and of its content
		elements     []string // elements a ruleset can match, nil when unknown
	}
	var blocks []block
	inRuleset := func() bool {
		return slices.ContainsFunc(blocks, func(b block) bool { return b.name == "" })
	}
	open := func(name, prelude string, elements []string) {
		start := out.Len()
		out.WriteString(prelude + "{")
		blocks = append(blocks, block{name, start, out.Len(), elements})
	}
	for {
		gt, tt, data := p.Next()
		if gt == css.ErrorGrammar {
			if tt == css.ErrorToken || !p.HasParseError() {
				break
			}
			if c.dropStyleOnSyntaxError {
				return "", false
			}
			continue
		}

		if skipping > 0 {
			switch gt {
			case css.BeginAtRuleGrammar, css.BeginRulesetGrammar:
				skipping++
			case css.EndAtRuleGrammar, css.EndRulesetGrammar:
				skipping--
			}
			continue
		}

		switch gt {
		case css.AtRuleGrammar:
			name := string(data)
			if !c.atRuleSupported(name, p.Values()) {
				if c.dropStyleOnUnsupportedAtRule {
					return "", false
				}
				continue
			}
			out.WriteString(atRule(name, p.Values()) + ";")
		case css.BeginAtRuleGrammar:
			name := string(data)
			if !c.atRuleSupported(name, p.Values()) {
				if c.dropStyleOnUnsupportedAtRule {
					return "", false
				}
				skipping = 1
				continue
			}
			open(name, atRule(name, p.Values()), nil)
		case css.BeginRulesetGrammar:
			if inRuleset() && c.unsupported("css-nesting") {
				skipping = 1
				continue
			}
			selector, subjects, ok := c.selectorList(p.Values())
			if !ok {
				skipping = 1
				continue
			}
			open("", selector, subjectElements(subjects, usage))
		case css.EndRulesetGrammar, css.EndAtRuleGrammar:
			if len(blocks) == 0 {
				continue
			}
			b := blocks[len(blocks)-1]
			blocks = blocks[:len(blocks)-1]
			if out.Len() == b.inner {
				out.Truncate(b.start)
				continue
			}
			out.WriteString("}")
		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			if c.poisons != nil && c.poisons(p.Values()) {
				return "", false
			}
			if len(blocks) > 0 && !inRuleset() {
				// Descriptors of @font-face, @page and the like, not properties.
				out.WriteString(string(data) + ":" + tokensString(p.Values()) + ";")
				continue
			}
			var elements []string
			if len(blocks) > 0 {
				elements = blocks[len(blocks)-1].elements
			}
			if declaration, ok := c.declaration("", elements, string(data), p.Values()); ok {
				out.WriteString(declaration + ";")
			}
		case css.TokenGrammar:
			if tt != css.CDOToken && tt != css.CDCToken {
				out.Write(data)
			}
		}
	}

	// Re-serializing drops whitespace around some delimiters, which could assemble
	// a "</style" that ends the element early; an escaped "<" is the same character to CSS.
	return strings.ReplaceAll(out.String(), "<", `\3c `), true
}

func atRule(name string, prelude []css.Token) string {
	if p := strings.TrimSpace(tokensString(prelude)); p != "" {
		return name + " " + p
	}
	return name
}

// subjectElements resolves selector subjects to the element names they can match; nil means any.
func subjectElements(subjects []string, usage map[string]map[string]bool) []string {
	if subjects == nil {
		return nil
	}
	elements := []string{}
	for _, subject := range subjects {
		if subject[0] != '.' && subject[0] != '#' {
			elements = append(elements, subject)
			continue
		}
		for element := range usage[subject] {
			elements = append(elements, element)
		}
	}
	return elements
}

// declaration returns "property:value" as the client keeps it, or ok=false when
// the client drops it. element is the element of a style attribute and empty
// in a stylesheet; elements are the elements the declaration can apply to, nil when unknown.
func (c *Client) declaration(element string, elements []string, property string, value []css.Token) (string, bool) {
	if c.rewriteDeclaration != nil {
		var keep bool
		property, value, keep = c.rewriteDeclaration(element, property, value)
		if !keep {
			return "", false
		}
	}
	if !strings.HasPrefix(property, "--") && c.properties != nil && !c.properties[property] {
		return "", false
	}
	if allowed, ok := c.propertyElements[property]; ok && elements != nil {
		if !slices.ContainsFunc(elements, func(e string) bool { return slices.Contains(allowed, e) }) {
			return "", false
		}
	}
	for _, slug := range declarationFeatures(property, value) {
		if c.unsupported(slug) {
			return "", false
		}
	}
	return property + ":" + tokensString(value), true
}

func (c *Client) atRuleSupported(name string, prelude []css.Token) bool {
	switch strings.TrimPrefix(unprefixed(strings.TrimPrefix(name, "@")), "-") {
	case "media":
		return c.mediaQuerySupported(prelude)
	case "font-face":
		return !c.unsupported("css-at-font-face")
	case "import":
		return !c.unsupported("css-at-import")
	case "keyframes":
		return !c.unsupported("css-at-keyframes")
	case "supports":
		return !c.unsupported("css-at-supports")
	case "charset", "namespace":
		return false
	}
	return true
}

func (c *Client) mediaQuerySupported(prelude []css.Token) bool {
	if c.unsupported("css-at-media") {
		return false
	}
	for i, token := range prelude {
		if token.TokenType != css.IdentToken {
			continue
		}
		name := strings.ToLower(string(token.Data))
		if i > 0 && prelude[i-1].TokenType == css.LeftParenthesisToken {
			if slug := mediaFeatureFeature(name); slug != "" && c.unsupported(slug) {
				return false
			}
			if c.mediaFeatures != nil && !c.mediaFeatures[name] {
				return false
			}
			continue
		}
		switch name {
		case "and", "only":
			continue
		case "not", "or":
			if c.mediaTypes == nil {
				continue
			}
		}
		if c.mediaTypes != nil && !c.mediaTypes[name] {
			return false
		}
	}
	return true
}

// selectorList returns the selectors of a rule the client keeps, rewritten the
// way the client rewrites them, or ok=false when none of them would match.
// subjects are the selectors' subjects (see complexSelector), nil when any of
// them could be any element.
func (c *Client) selectorList(tokens []css.Token) (selector string, subjects []string, ok bool) {
	var kept []string
	selectors := splitSelectors(tokens)
	if len(selectors) > 1 && c.unsupported("css-selector-grouping") {
		return "", nil, false
	}
	anyElement := false
	for _, selector := range selectors {
		if s, subject, ok := c.complexSelector(selector); ok {
			kept = append(kept, s)
			subjects = append(subjects, subject)
			anyElement = anyElement || subject == ""
		}
	}
	if anyElement {
		subjects = nil
	}
	return strings.Join(kept, ","), subjects, len(kept) > 0
}

func splitSelectors(tokens []css.Token) [][]css.Token {
	var selectors [][]css.Token
	depth, start := 0, 0
	for i, token := range tokens {
		switch token.TokenType {
		case css.FunctionToken, css.LeftParenthesisToken, css.LeftBracketToken:
			depth++
		case css.RightParenthesisToken, css.RightBracketToken:
			depth--
		case css.CommaToken:
			if depth == 0 {
				selectors = append(selectors, trimWhitespace(tokens[start:i]))
				start = i + 1
			}
		}
	}
	return append(selectors, trimWhitespace(tokens[start:]))
}

func trimWhitespace(tokens []css.Token) []css.Token {
	for len(tokens) > 0 && tokens[0].TokenType == css.WhitespaceToken {
		tokens = tokens[1:]
	}
	for len(tokens) > 0 && tokens[len(tokens)-1].TokenType == css.WhitespaceToken {
		tokens = tokens[:len(tokens)-1]
	}
	return tokens
}

// legacyPseudoElements may be written with a single colon.
var legacyPseudoElements = map[string]bool{"before": true, "after": true, "first-letter": true, "first-line": true}

// complexSelector also returns the subject of the selector: its rightmost
// type, or else its rightmost first class or id as ".name" or "#name"; "" when
// it could be any element.
func (c *Client) complexSelector(tokens []css.Token) (string, string, bool) {
	var slugs []string
	var out strings.Builder
	// The current compound selector: classes and ids (caniemail's "chaining"
	// counts only these), other simple selectors, its type and its first class or id.
	classIDs, others := 0, 0
	typeName, firstClassID := "", ""
	chained := false
	nesting := 0 // inside brackets or pseudo-class arguments, which hold no selectors of their own to check

	endCompound := func() {
		if classIDs > 1 {
			chained = true
		}
		classIDs, others = 0, 0
		typeName, firstClassID = "", ""
	}
	// classOrID reports whether the client ignores this class or id because it isn't the first in its compound.
	classOrID := func(slug, key string) bool {
		classIDs++
		slugs = append(slugs, slug)
		if classIDs == 1 {
			firstClassID = key
		}
		return classIDs > 1 && c.chainedSelectors == chainFirstOnly
	}
	other := func(slug string) {
		others++
		slugs = append(slugs, slug)
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		data := string(token.Data)
		if nesting > 0 {
			switch token.TokenType {
			case css.FunctionToken, css.LeftParenthesisToken, css.LeftBracketToken:
				nesting++
			case css.RightParenthesisToken, css.RightBracketToken:
				nesting--
			}
			out.WriteString(data)
			continue
		}

		switch token.TokenType {
		case css.WhitespaceToken:
			endCompound()
			if next := i + 1; next < len(tokens) && isCombinator(tokens[next]) {
				continue
			}
			slugs = append(slugs, "css-selector-descendant")
		case css.DelimToken:
			switch data {
			case ">":
				endCompound()
				slugs = append(slugs, "css-selector-child")
			case "+":
				endCompound()
				slugs = append(slugs, "css-selector-adjacent-sibling")
			case "~":
				endCompound()
				slugs = append(slugs, "css-selector-general-sibling")
			case "*":
				other("css-selector-universal")
			case ".":
				if i+1 >= len(tokens) {
					break
				}
				name := string(tokens[i+1].Data)
				if classOrID("css-selector-class", "."+name) {
					i++
					continue
				}
				out.WriteString(data)
				i++
				data = name
			}
		case css.HashToken:
			if classOrID("css-selector-id", data) {
				continue
			}
		case css.IdentToken:
			typeName = strings.ToLower(data)
			slugs = append(slugs, "css-selector-type")
			switch typeName {
			case "html":
				if !c.keepRootSelectors {
					return "", "", false
				}
			case "body":
				if !c.keepBodySelectors {
					return "", "", false
				}
				if c.webmail {
					// The client applies body rules to the <div> that replaces <body>.
					data = "div[" + bodyMarker + "]"
				}
			}
		case css.LeftBracketToken:
			if c.attributeOperators != nil && !c.attributeOperators[attributeOperator(tokens[i:])] {
				return "", "", false
			}
			nesting++
			other("css-selector-attribute")
		case css.ColonToken:
			pseudoElement := i+1 < len(tokens) && tokens[i+1].TokenType == css.ColonToken
			if pseudoElement {
				out.WriteString(data)
				i++
			}
			if i+1 >= len(tokens) {
				break
			}
			i++
			name := tokens[i].Data
			out.WriteString(data)
			data = string(name)
			if tokens[i].TokenType == css.FunctionToken {
				nesting++
			}
			lower := strings.ToLower(strings.TrimSuffix(data, "("))
			if lower == "root" && !c.keepRootSelectors {
				return "", "", false
			}
			if c.pseudoOnTypeOnly && (typeName == "" || classIDs+others > 0) {
				return "", "", false
			}
			if pseudoElement || legacyPseudoElements[lower] {
				slugs = append(slugs, "css-pseudo-element-"+lower)
			} else {
				slugs = append(slugs, "css-pseudo-class-"+lower)
			}
		}
		out.WriteString(data)
	}
	subject := typeName
	if subject == "" {
		subject = firstClassID
	}
	endCompound()
	if chained {
		slugs = append(slugs, "css-selector-chaining")
	}

	for _, slug := range slugs {
		if !caniemailFeatures[slug] || !c.unsupported(slug) {
			continue
		}
		if slug == "css-selector-chaining" && c.chainedSelectors == chainFirstOnly {
			continue
		}
		return "", "", false
	}
	return out.String(), subject, true
}

// attributeOperator returns the operator of the attribute selector starting at
// tokens[0], such as "=" or "^=", or "" for a bare [attr].
func attributeOperator(tokens []css.Token) string {
	for _, token := range tokens[1:] {
		switch token.TokenType {
		case css.RightBracketToken:
			return ""
		case css.DelimToken:
			if string(token.Data) == "=" {
				return "="
			}
		case css.IncludeMatchToken, css.DashMatchToken, css.PrefixMatchToken, css.SuffixMatchToken, css.SubstringMatchToken:
			return string(token.Data)
		}
	}
	return ""
}

func isCombinator(token css.Token) bool {
	if token.TokenType != css.DelimToken {
		return false
	}
	switch string(token.Data) {
	case ">", "+", "~":
		return true
	}
	return false
}
