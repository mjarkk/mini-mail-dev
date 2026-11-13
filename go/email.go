package src

import (
	"net/mail"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/jhillyerd/enmime/v2"
	"github.com/oklog/ulid/v2"
)

// Address with json attributes
type Address struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Attachment with json attributes
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"-"`
}

// ConvertAttachment converts a decoded enmime.Part to a Attachment
func ConvertAttachment(part *enmime.Part) Attachment {
	return Attachment{
		Filename:    part.FileName,
		ContentType: part.ContentType,
		Data:        part.Content,
	}
}

// EmbeddedFile with json attributes
type EmbeddedFile struct {
	M *sync.Mutex `json:"-"`

	CID         string `json:"cid"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"-"`
}

// ConvertEmbeddedFile converts a decoded enmime.Part to a EmbeddedFile
func ConvertEmbeddedFile(part *enmime.Part) EmbeddedFile {
	return EmbeddedFile{
		M:           &sync.Mutex{},
		CID:         part.ContentID,
		ContentType: part.ContentType,
		Data:        part.Content,
	}
}

type EmailRemainder struct {
	Raw string `json:"raw"`

	Header mail.Header `json:"header"`

	ReplyTo    []*Address `json:"replyTo"`
	Cc         []*Address `json:"cc"`
	Bcc        []*Address `json:"bcc"`
	Date       time.Time  `json:"date"`
	MessageID  string     `json:"messageId"`
	InReplyTo  []string   `json:"inReplyTo"`
	References []string   `json:"references"`

	ResentFrom      []*Address `json:"resentFrom"`
	ResentSender    *Address   `json:"resentSender"`
	ResentTo        []*Address `json:"resentTo"`
	ResentDate      time.Time  `json:"resentDate"`
	ResentCc        []*Address `json:"resentCc"`
	ResentBcc       []*Address `json:"resentBcc"`
	ResentMessageID string     `json:"resentMessageId"`

	ContentType string `json:"contentType"`

	// HTMLBody is the HTML as received; render it through sanitize.Parse.
	HTMLBody string `json:"htmlBody"`
	TextBody string `json:"textBody"`

	Attachments   []Attachment   `json:"attachments"`
	EmbeddedFiles []EmbeddedFile `json:"embeddedFiles"`
}

// Email contains a single email
type Email struct {
	// Added by our server
	ID       ulid.ULID `json:"id"`
	RealDate time.Time `json:"realDate"`
	Subject  string    `json:"subject"`
	BodyType string    `json:"bodyType"` // "html" or "text"
	BodyHint string    `json:"bodyHint"`
	Read     bool      `json:"read"`

	// Recived form the client
	Sender *Address   `json:"sender"`
	From   []*Address `json:"from"`
	To     []*Address `json:"to"`

	Remainder EmailRemainder `json:"-"`
}

// ConvertEmail converts a parsed enmime.Envelope to a Email
func ConvertEmail(ulid ulid.ULID, emlRawData []byte, env *enmime.Envelope, realDate time.Time, realFrom, realTo string) Email {
	attachments, embeddedFiles := splitParts(env)

	bodyHint := strings.Join(strings.Fields(env.Text), " ")
	if hint := []rune(bodyHint); len(hint) > 100 {
		bodyHint = string(hint[:100]) + "..."
	}

	date, _ := env.Date()
	resentDate, _ := mail.ParseDate(env.GetHeader("Resent-Date"))

	resp := Email{
		ID:       ulid,
		RealDate: realDate,
		Subject:  env.GetHeader("Subject"),
		Sender:   firstAddress(env, "Sender"),
		From:     addressList(env, "From"),
		To:       addressList(env, "To"),
		BodyHint: bodyHint,
		BodyType: "text",
		Remainder: EmailRemainder{
			Raw:             string(emlRawData),
			Header:          mail.Header(env.Root.Header),
			ReplyTo:         addressList(env, "Reply-To"),
			Cc:              addressList(env, "Cc"),
			Bcc:             addressList(env, "Bcc"),
			Date:            date,
			MessageID:       messageID(env.GetHeader("Message-ID")),
			InReplyTo:       messageIDs(env.GetHeader("In-Reply-To")),
			References:      messageIDs(env.GetHeader("References")),
			ResentFrom:      addressList(env, "Resent-From"),
			ResentSender:    firstAddress(env, "Resent-Sender"),
			ResentTo:        addressList(env, "Resent-To"),
			ResentDate:      resentDate,
			ResentCc:        addressList(env, "Resent-Cc"),
			ResentBcc:       addressList(env, "Resent-Bcc"),
			ResentMessageID: messageID(env.GetHeader("Resent-Message-ID")),
			ContentType:     env.GetHeader("Content-Type"),
			HTMLBody:        env.HTML,
			TextBody:        env.Text,
			Attachments:     attachments,
			EmbeddedFiles:   embeddedFiles,
		},
	}

	if env.HTML != "" {
		resp.BodyType = "html"
	}

	if resp.Sender == nil && realFrom != "" {
		resp.Sender = &Address{Address: realFrom}
	}

	if len(resp.To) == 0 && realTo != "" {
		resp.To = []*Address{{Address: realTo}}
	}

	return resp
}

// ConvertAddressList converts a list of mail.Address to a list of Address
func ConvertAddressList(adList []*mail.Address) []*Address {
	if adList == nil {
		return nil
	}

	var result []*Address
	for _, ad := range adList {
		result = append(result, &Address{
			Name:    ad.Name,
			Address: ad.Address,
		})
	}
	return result
}

func splitParts(env *enmime.Envelope) (attachments []Attachment, embeddedFiles []EmbeddedFile) {
	for _, part := range env.Attachments {
		attachments = append(attachments, ConvertAttachment(part))
	}

	for _, part := range slices.Concat(env.Inlines, env.OtherParts) {
		if part.ContentID != "" {
			embeddedFiles = append(embeddedFiles, ConvertEmbeddedFile(part))
		} else if part.FileName != "" {
			// Without a Content-ID the HTML can't reference it, so it would otherwise be invisible.
			attachments = append(attachments, ConvertAttachment(part))
		}
	}

	return attachments, embeddedFiles
}

// A missing or malformed address header leaves the field empty instead of rejecting the email.
func addressList(env *enmime.Envelope, key string) []*Address {
	list, _ := env.AddressList(key)
	return ConvertAddressList(list)
}

func firstAddress(env *enmime.Envelope, key string) *Address {
	list := addressList(env, key)
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

func messageID(value string) string {
	return strings.Trim(value, "<> ")
}

func messageIDs(value string) []string {
	var ids []string
	for _, id := range strings.Fields(value) {
		ids = append(ids, messageID(id))
	}
	return ids
}
