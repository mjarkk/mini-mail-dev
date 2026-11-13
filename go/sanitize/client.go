package sanitize

import "github.com/tdewolff/parse/v2/css"

//go:generate go run ./internal/caniemailgen

type support uint8

const (
	supported support = iota
	partial
	unsupported
)

type chainMode uint8

const (
	chainAsCaniemail chainMode = iota
	// chainFirstOnly applies ".a.b" as ".a", like Outlook's Word engine.
	chainFirstOnly
)

// Client is an email client whose rendering Parse simulates.
type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Width is the viewport width in CSS pixels for a mobile client, 0 for a desktop one.
	Width int `json:"width,omitempty"`

	// overrides take precedence over caniemail's data; research resolves its
	// partial ("a") entries into these.
	overrides map[string]support

	// webmail clients embed the message in their own page: the document always
	// renders in standards mode and <body> becomes a <div>.
	webmail        bool
	bodyAttributes map[string]bool
	// standardsMode renders every message in standards mode, whatever its doctype.
	standardsMode bool

	// msoVersion > 0 evaluates Outlook conditional comments as that Office version.
	msoVersion int

	// baseCSS stands in for the client's own styles around the message.
	baseCSS string

	// properties, when set, is the complete list of CSS properties the client keeps.
	properties map[string]bool
	// propertyElements limits inline CSS properties to the listed elements.
	propertyElements map[string][]string

	mediaTypes    map[string]bool
	mediaFeatures map[string]bool

	// maxStyleBytes caps the combined size of all <style> blocks; the block
	// crossing it and every block after it are dropped.
	maxStyleBytes                int
	styleOnlyInHead              bool
	dropStyleOnSyntaxError       bool
	dropStyleOnUnsupportedAtRule bool
	// truncateStyleAtDoubleBrace drops everything after the first "}}" in a <style>.
	truncateStyleAtDoubleBrace bool
	// poisons reports a declaration value that makes the client discard the
	// whole style attribute or <style> block it appears in.
	poisons func(value []css.Token) bool

	keepBodySelectors bool // body {} still styles the message
	keepRootSelectors bool // html {} and :root {} still match
	chainedSelectors  chainMode
	// pseudoOnTypeOnly keeps pseudo-classes and pseudo-elements only on a bare type selector, as in a:hover.
	pseudoOnTypeOnly bool
	// attributeOperators, when set, lists the attribute selector operators the
	// client keeps; "" stands for a bare [attr].
	attributeOperators map[string]bool

	unwrapElements map[string]bool
	// doctypeMarker replaces the doctype in front of the message, like Gmail's <u></u>.
	doctypeMarker string

	// stripAttributes may contain "data-*" to strip every data attribute.
	stripAttributes map[string]bool
	// dropClassOnSelfClosing lists elements whose class is dropped when written as <img ... />.
	dropClassOnSelfClosing map[string]bool
	// firstClassOnly keeps only the first name of a class attribute.
	firstClassOnly bool
	// vml draws Outlook's VML shapes (bulletproof buttons and backgrounds) as HTML.
	vml bool
	// blockWebFonts stops fonts from loading, including those an @import pulls in.
	blockWebFonts bool

	// rewriteDeclaration runs before any other check; keep=false drops the declaration.
	// element is empty for declarations in a stylesheet.
	rewriteDeclaration func(element, property string, value []css.Token) (newProperty string, newValue []css.Token, keep bool)

	clipBytes  int
	clipNotice string
}

func (c *Client) unsupported(slug string) bool {
	level, ok := c.overrides[slug]
	if !ok {
		level = caniemail[c.ID][slug]
	}
	return level == unsupported
}

// ContentSecurityPolicy is the policy to serve a page from Parse with.
func (c *Client) ContentSecurityPolicy() string {
	policy := "script-src 'none'; object-src 'none'; frame-src 'none'; form-action 'none'"
	if c.blockWebFonts {
		policy += "; font-src 'none'"
	}
	return policy
}

// Clients lists the simulated clients, the default first.
var Clients = []*Client{gmail, appleMail, outlookCom, outlookWindows, yahoo, samsungEmail, thunderbird}

// ClientByID returns nil for an unknown id.
func ClientByID(id string) *Client {
	for _, client := range Clients {
		if client.ID == id {
			return client
		}
	}
	return nil
}

func set(items ...string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, item := range items {
		m[item] = true
	}
	return m
}
