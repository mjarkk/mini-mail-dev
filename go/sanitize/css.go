package sanitize

import (
	"fmt"
	"io"
	"strings"

	"github.com/tdewolff/parse/css"
)

// https://www.caniemail.com/features/

type cssParserT struct {
	lexer  *css.Lexer
	output string
}

func (p *cssParserT) parseSelector() error {
	toAppend := ""

	for {
		tt, text := p.lexer.Next()
		fmt.Printf("%s: \"%s\"\n", tt, string(text))
		switch tt {
		case css.ColonToken:
			toAppend += ":"
		case css.ErrorToken:
			return p.lexer.Err()
		case css.DelimToken:
			switch string(text) {
			case ".":
				toAppend += "."
			case "=":
				toAppend += "="
			case "*":
				toAppend += "*"
			case ">":
				toAppend += ">"
			case "+":
				toAppend += "+"
			case "~":
				toAppend += "~"
			}
		case css.WhitespaceToken:
			toAppend += " "
		case css.HashToken:
			toAppend += string(text)
		case css.IdentToken:
			toAppend += string(text)
		case css.LeftBracketToken:
			toAppend += "["
		case css.RightBracketToken:
			toAppend += "]"
		case css.LeftBraceToken: // {
			p.output += strings.TrimSpace(toAppend) + "{"
			return p.parseSelectorContent()
		case css.StringToken: // "Foo", 'Foo'
			toAppend += string(text)
		case css.SuffixMatchToken: // $=
			toAppend += "$="
		case css.DashMatchToken: // |=
			toAppend += "|="
		case css.PrefixMatchToken: // ^=
			toAppend += "^="
		case css.IncludeMatchToken: // ~=
			toAppend += "~="
		case css.SubstringMatchToken: // *=
			toAppend += "*="
		case css.CommaToken:
			toAppend += ","
		case css.FunctionToken:
		case css.AtKeywordToken:
		case css.BadStringToken:
		case css.URLToken:
		case css.BadURLToken:
		case css.NumberToken:
		case css.PercentageToken:
		case css.DimensionToken:
		case css.UnicodeRangeToken:
		case css.ColumnToken:
		case css.CDOToken:
		case css.CDCToken:
		case css.SemicolonToken:
		case css.LeftParenthesisToken:
		case css.RightParenthesisToken:
		case css.RightBraceToken:
		case css.CommentToken:
		case css.EmptyToken:
		case css.CustomPropertyNameToken:
		case css.CustomPropertyValueToken:
		}
	}
}

func (p *cssParserT) parseSelectorContent() error {
	for {
		tt, text := p.lexer.Next()
		fmt.Printf("%s: \"%s\"\n", tt, string(text))
		switch tt {
		case css.ErrorToken:
			return p.lexer.Err()
		case css.IdentToken:
			textStr := string(text)
			p.output += textStr
		case css.RightBraceToken:
			p.output += "}"
			return nil
		case css.ColonToken:
			p.output += ":"
		case css.SemicolonToken:
			p.output += ";"
		case css.FunctionToken:
			err := p.handleFunction(string(text))
			if err != nil {
				return err
			}
		case css.LeftParenthesisToken:
			p.output += "("
		case css.RightParenthesisToken:
			p.output += ")"
		case css.DimensionToken:
			p.handleDimention(string(text))
		case css.NumberToken:
			p.output += string(text)
		case css.URLToken:
			p.output += string(text)
		case css.DelimToken:
			switch string(text) {
			case "-":
				p.output += "-"
			}
		case css.AtKeywordToken:
		case css.HashToken:
		case css.StringToken:
		case css.BadStringToken:
		case css.BadURLToken:
		case css.PercentageToken:
		case css.UnicodeRangeToken:
		case css.IncludeMatchToken:
		case css.DashMatchToken:
		case css.PrefixMatchToken:
		case css.SuffixMatchToken:
		case css.SubstringMatchToken:
		case css.ColumnToken:
		case css.WhitespaceToken:
		case css.CDOToken:
		case css.CDCToken:
		case css.CommaToken:
		case css.LeftBracketToken:
		case css.RightBracketToken:
		case css.LeftBraceToken:
		case css.CommentToken:
		case css.EmptyToken:
		case css.CustomPropertyNameToken:
		case css.CustomPropertyValueToken:
		}
	}
}

func (p *cssParserT) handleDimention(text string) {
	// https://developer.mozilla.org/en-US/docs/Web/CSS/dimension

	if len(text) == 0 {
		return
	}

	splitIndex := 0
	for index, c := range text {
		if c >= '0' && c <= '9' {
			continue
		}
		if index == 0 && c == '-' {
			continue
		}
		if index == len(text)-1 && c == '.' {
			continue
		}
		splitIndex = index
	}

	unit := text
	if splitIndex > 0 {
		unit = unit[splitIndex-1:]
	}

	switch unit {
	case "em", "ex", "ch", "rem", "lh", "vw", "vh", "vmin", "vmax", "cm", "mm", "Q", "in", "pc", "pt", "px": // Relative length units
		p.output += text
	case "deg", "grad", "rad", "turn": // Angle units
		p.output += text
	case "s", "ms": // Time units
		p.output += text
	case "Hz", "kHz": // Frequency units
		p.output += text
	case "dpi", "dpcm", "dppx": // Resolution units
		p.output += text
	case "fr": // Grid units
		p.output += text
	case "":
		p.output += text
	}
}

func (p *cssParserT) handleFunction(text string) error {
	// https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Functions
	switch strings.TrimRight(text, "(") {
	case "calc": // Basic functions
		p.output += text
	case "min", "max", "clamp": // Comparison functions
		p.output += text
	case "round", "mod", "rem": // Stepped value functions
		p.output += text
	case "sin", "cos", "tan", "asin", "acos", "atan", "atan2": // Trigonometric functions
		p.output += text
	case "pow", "sqrt", "hypot", "log", "exp": // Exponential functions
		p.output += text
	case "abs", "sign": // Sign-related functions
		p.output += text
	case "rgba", "hsl", "hwb", "lch", "oklch", "lab", "oklab", "color", "color-mix", "color-contrast", "device-cmyk": // Color functions
		p.output += text
	case "linear-gradient", "radial-gradient", "conic-gradient", "repeating-linear-gradient", "repeating-radial-gradient", "repeating-conic-gradient": // Gradient functions
		p.output += text
	case "counter", "counters", "symbols": // Counter functions
		p.output += text
	case "attr":
		p.output += text
	case "url":
		// TODO: Look into xss attacks
		p.output += text
	case "fit-content", "minmax", "repeat": // Grid functions
		p.output += text
	case "styleistic", "styleset", "character-variant", "swash", "ornaments", "annotation":
		p.output += text
	default:
		// Not supported by gmail:
		//   - Translate functions
		//   - Rotation functions
		//   - Scaling functions
		//   - Skew functions
		//   - Matrix functions
		//   - Perspective functions
		//   - Filter functions
		//   - Image functions
		//   - Shape functions
		//   - Easing functions
		//   - Animation functions
		// Parse unsupported by gmail:
		//   - Reference functions
		return p.parseUnsupportedFunction()
	}
	return nil
}

func (p *cssParserT) parseUnsupportedFunction() error {
	for {
		tt, _ := p.lexer.Next()

		switch tt {
		case css.ErrorToken:
			return p.lexer.Err()
		case css.FunctionToken:
			return p.parseUnsupportedFunction()
		case css.LeftParenthesisToken:
			return p.parseUnsupportedFunction()
		case css.RightParenthesisToken:
			p.output += ")"
			return nil
		}
	}
}

func (p *cssParserT) parse() error {
	for {
		err := p.parseSelector()
		if err == nil {
			continue
		}

		if p.lexer.Err() == io.EOF {
			return nil
		}
		return err
	}
}

func parseCss(in string) (string, error) {
	parser := &cssParserT{
		lexer:  css.NewLexer(strings.NewReader(in)),
		output: "",
	}

	err := parser.parse()
	return parser.output, err
}
