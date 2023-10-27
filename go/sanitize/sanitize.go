package sanitize

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Parser struct {
	tokenizer *html.Tokenizer
	resp      []byte

	// pre tag counter is used to keep track of the n-th depth of pre tags
	// if the number is more than 0 we should not trim white space in text tokens
	preTagCounter uint
}

// Parse parses a html document and returns a sanitized version of it
func Parse(document string) (string, error) {
	document = strings.TrimSpace(document)
	if document == "" {
		return "", nil
	}

	parser := &Parser{
		tokenizer:     html.NewTokenizer(strings.NewReader(document)),
		resp:          []byte{},
		preTagCounter: 0,
	}
	err := parser.parse()
	return string(parser.resp), err
}

func (p *Parser) append(content []byte) {
	p.resp = append(p.resp, content...)
}

func (p *Parser) appendStr(content string) {
	p.resp = append(p.resp, content...)
}

func (p *Parser) next() html.TokenType {
	return p.tokenizer.Next()
}

func (p *Parser) err() error {
	return p.tokenizer.Err()
}

func (p *Parser) parseStyle() error {
	token := p.next()
	if token == html.ErrorToken {
		return p.err()
	}
	if token == html.EndTagToken {
		tokenValue := p.tokenizer.Token().Data
		if tokenValue != "style" {
			return fmt.Errorf("expected style end tag token but got %v", tokenValue)
		}
		return nil
	}
	if token != html.TextToken {
		return fmt.Errorf("expected style text token, got %v", token)
	}

	p.appendStr("<style>")
	p.append(p.tokenizer.Raw())

	defer func() {
		p.appendStr("</style>")
	}()

	token = p.next()
	if token == html.ErrorToken {
		return p.err()
	}
	if token != html.EndTagToken {
		return fmt.Errorf("expected style end tag token but got %v", token)
	}
	tokenValue := p.tokenizer.Token().Data
	if tokenValue != "style" {
		return fmt.Errorf("expected style end tag token but got %v", tokenValue)
	}

	return nil
}

func (p *Parser) parseIgnoreScript() error {
	token := p.next()
	if token == html.ErrorToken {
		return p.err()
	}
	if token == html.EndTagToken {
		tokenValue := p.tokenizer.Token().Data
		if tokenValue != "script" {
			return fmt.Errorf("expected script end tag token but got %v", tokenValue)
		}
		return nil
	}
	if token != html.TextToken {
		return fmt.Errorf("expected script text token, got %v", token)
	}

	token = p.next()
	if token == html.ErrorToken {
		return p.err()
	}
	if token != html.EndTagToken {
		return fmt.Errorf("expected script end tag token but got %v", token)
	}
	tokenValue := p.tokenizer.Token().Data
	if tokenValue != "script" {
		return fmt.Errorf("expected script end tag token but got %v", tokenValue)
	}

	return nil
}

func formatErr(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}

func elementTokenAllowed(el string) bool {
	lastIdx := len(el) - 1
	for idx, l := range el {
		if l >= 'a' && l <= 'z' {
			continue
		}
		if l >= 'A' && l <= 'Z' {
			continue
		}
		if idx != 0 {
			if l >= '0' && l <= '9' {
				continue
			}
			if idx != lastIdx && (l == '-' || l == '_') {
				continue
			}
		}
		return false
	}

	switch el {
	case "link":
		return false
	case "area", "audio", "track", "video":
		return false
	case "embed", "iframe", "object", "picture", "portal", "source":
		return false
	case "canvas", "script":
		return false
	case "button", "datalist", "form", "input", "meter", "optgroup", "option", "output", "select", "textarea":
		return false
	case "dialog", "menu":
		return false
	case "marquee":
		return false
	case "applet":
		return false
	case "xml", "svg":
		// TODO Also ignore everything within the xml/svg tag, should write a full ignore tag method
		return false
	case "style":
		// Style should be parsed separately
		return false
	}

	return true
}

func attributeKeyValueAllowed(tagName, key, value string) (string, bool) {
	switch key {
	case "lang":
		if tagName != "html" {
			return "", false
		}

		return value, true
	case "style", "class", "id":
		return value, true
	case "href":
		if tagName != "a" {
			return "", false
		}

		// TODO: This is probably not enough xss protection
		if strings.Contains(value, "javascript:") {
			return "", false
		}

		return value, true
	case "src", "alt", "sizes", "srcset":
		if tagName != "img" {
			return "", false
		}

		// TODO: This is probably not enough xss protection
		if strings.Contains(value, "javascript:") {
			return "", false
		}

		return value, true
	case "colspan", "rowspan", "headers", "align", "border", "cellpadding", "cellspacing": // tables
		switch tagName {
		case "td", "th", "table", "tr", "col", "colgroup", "tbody", "thead", "tfoot":
			return value, true
		}
		return "", false
	case "valign", "width", "height":
		return value, true
	}

	if strings.HasPrefix(key, "data-") {
		lastIdx := len(key) - 1
		for idx, c := range strings.TrimPrefix(key, "data-") {
			if c >= 'a' && c <= 'z' {
				continue
			}
			if c >= 'A' && c <= 'Z' {
				continue
			}
			if c >= '0' && c <= '9' {
				continue
			}
			if idx != lastIdx {
				if c == '-' || c == '_' {
					continue
				}
			}
			return "", false
		}

		return value, true
	}

	return "", false
}

func (p *Parser) appendTag(token html.Token, selfClosing bool) {
	if !elementTokenAllowed(token.Data) {
		return
	}

	props := []string{}
	for _, v := range token.Attr {
		value, ok := attributeKeyValueAllowed(token.Data, v.Key, strings.TrimSpace(v.Val))
		if !ok {
			// fmt.Println("DISALLOWED ATTR", token.Data, v.Key, strings.TrimSpace(v.Val))
			continue
		}
		if v.Val == "" {
			props = append(props, v.Key)
		}

		props = append(props, v.Key+"=\""+html.EscapeString(value)+"\"")
	}

	propsAsString := strings.Join(props, " ")
	if len(props) > 0 {
		propsAsString = " " + propsAsString
	}

	p.appendStr("<" + token.Data + propsAsString)
	if selfClosing {
		p.appendStr("/>")
	} else {
		p.appendStr(">")
	}
}

func (p *Parser) parse() error {
	for {
		token := p.next()

		switch token {
		case html.ErrorToken:
			return formatErr(p.err())
		case html.TextToken:
			textContents := html.EscapeString(string(p.tokenizer.Raw()))

			if p.preTagCounter == 0 {
				textContents = strings.TrimSpace(textContents)
			}

			p.appendStr(textContents)
		case html.StartTagToken:
			token := p.tokenizer.Token()
			switch token.Data {
			case "style":
				err := p.parseStyle()
				if err != nil {
					return formatErr(err)
				}
			case "script":
				err := p.parseIgnoreScript()
				if err != nil {
					return formatErr(err)
				}
			case "pre":
				p.preTagCounter++
				fallthrough
			default:
				p.appendTag(token, false)
			}
		case html.EndTagToken:
			token := p.tokenizer.Token().Data
			if token == "pre" && p.preTagCounter > 0 {
				p.preTagCounter--
			}
			if elementTokenAllowed(token) {
				p.appendStr("</" + token + ">")
			}
		case html.SelfClosingTagToken:
			p.appendTag(p.tokenizer.Token(), true)
		case html.CommentToken:
			// Ignore
		case html.DoctypeToken:
			if len(p.resp) == 0 {
				p.append(p.tokenizer.Raw())
			}
		}
	}
}
