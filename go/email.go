package src

import (
	"io"
	"net/mail"
	"strings"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/microcosm-cc/bluemonday"
	"github.com/oklog/ulid/v2"
)

// Address with json attributes
type Address struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// MightConvertAddress converts *mail.Address to a *Address
func MightConvertAddress(ad *mail.Address) *Address {
	if ad == nil {
		return nil
	}
	return &Address{
		Name:    ad.Name,
		Address: ad.Address,
	}
}

// ConvertAddress converts mail.Address to a Address
func ConvertAddress(ad mail.Address) Address {
	return Address{
		Name:    ad.Name,
		Address: ad.Address,
	}
}

// Attachment with json attributes
type Attachment struct {
	Filename    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	Data        io.Reader `json:"-"`
}

// ConvertAttachment converts parsemail.Attachment to a Attachment
func ConvertAttachment(at parsemail.Attachment) Attachment {
	return Attachment{
		Filename:    at.Filename,
		ContentType: at.ContentType,
		Data:        at.Data,
	}
}

// EmbeddedFile with json attributes
type EmbeddedFile struct {
	CID         string    `json:"cid"`
	ContentType string    `json:"contentType"`
	Data        io.Reader `json:"-"`
}

// ConvertEmbeddedFile converts parsemail.EmbeddedFile to a EmbeddedFile
func ConvertEmbeddedFile(eb parsemail.EmbeddedFile) EmbeddedFile {
	return EmbeddedFile{
		CID:         eb.CID,
		ContentType: eb.ContentType,
		Data:        eb.Data,
	}
}

// Email contains a single email
type Email struct {
	// Added by our server
	ID       ulid.ULID `json:"id"`
	RealDate time.Time `json:"realDate"`

	// Recived form the client
	Header mail.Header `json:"header"`

	Subject    string     `json:"subject"`
	Sender     *Address   `json:"sender"`
	From       []*Address `json:"from"`
	ReplyTo    []*Address `json:"replyTo"`
	To         []*Address `json:"to"`
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

	ContentType string    `json:"contentType"`
	Content     io.Reader `json:"-"`

	HTMLBody string `json:"htmlBody"`
	TextBody string `json:"textBody"`

	Attachments   []Attachment   `json:"attachments"`
	EmbeddedFiles []EmbeddedFile `json:"embeddedFiles"`
}

// ConvertEmail converts parsemail.Email to a Email
func ConvertEmail(bluemondayPolicy *bluemonday.Policy, ulid ulid.ULID, em parsemail.Email, realDate time.Time, realFrom, realTo string) Email {
	resp := Email{
		ID:              ulid,
		RealDate:        realDate,
		Header:          em.Header,
		Subject:         em.Subject,
		Sender:          MightConvertAddress(em.Sender),
		From:            ConvertAddressList(em.From),
		ReplyTo:         ConvertAddressList(em.ReplyTo),
		To:              ConvertAddressList(em.To),
		Cc:              ConvertAddressList(em.Cc),
		Bcc:             ConvertAddressList(em.Bcc),
		Date:            em.Date,
		MessageID:       em.MessageID,
		InReplyTo:       em.InReplyTo,
		References:      em.References,
		ResentFrom:      ConvertAddressList(em.ResentFrom),
		ResentSender:    MightConvertAddress(em.ResentSender),
		ResentTo:        ConvertAddressList(em.ResentTo),
		ResentDate:      em.ResentDate,
		ResentCc:        ConvertAddressList(em.ResentCc),
		ResentBcc:       ConvertAddressList(em.ResentBcc),
		ResentMessageID: em.ResentMessageID,
		ContentType:     em.ContentType,
		Content:         em.Content,
		HTMLBody:        bluemondayPolicy.Sanitize(em.HTMLBody),
		TextBody:        em.TextBody,
		Attachments:     ConvertAttachmentList(em.Attachments),
		EmbeddedFiles:   ConvertEmbeddedFileList(em.EmbeddedFiles),
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

// ConvertAttachmentList converts a list of parsemail.Attachment to a list of Attachment
func ConvertAttachmentList(atList []parsemail.Attachment) []Attachment {
	if atList == nil {
		return nil
	}

	var result []Attachment
	for _, at := range atList {
		result = append(result, Attachment{
			Filename:    at.Filename,
			ContentType: at.ContentType,
			Data:        at.Data,
		})
	}
	return result
}

// ConvertEmbeddedFileList converts a list of parsemail.EmbeddedFile to a list of EmbeddedFile
func ConvertEmbeddedFileList(ebList []parsemail.EmbeddedFile) []EmbeddedFile {
	if ebList == nil {
		return nil
	}

	var result []EmbeddedFile
	for _, eb := range ebList {
		result = append(result, EmbeddedFile{
			CID:         eb.CID,
			ContentType: eb.ContentType,
			Data:        eb.Data,
		})
	}
	return result
}

// EmailHint contains the hints of a single email
type EmailHint struct {
	ID           ulid.ULID    `json:"id"`
	RealDate     time.Time    `json:"realDate"`
	Subject      string       `json:"subject"`
	Sender       *Address     `json:"sender"`
	From         []*Address   `json:"from"`
	To           []*Address   `json:"to"`
	TextBodyHint string       `json:"textBodyHint"`
	Attachments  []Attachment `json:"attachments"`
}

// EmailToEmailHint converts a Email to a EmailHint
func EmailToEmailHint(email Email) EmailHint {
	bodyHint := strings.NewReplacer(
		"\n", " ",
		"\r", " ",
		"\t", " ",
		"   ", " ",
		"  ", " ",
	).Replace(email.TextBody)
	if len(bodyHint) > 100 {
		bodyHint = bodyHint[:100] + "..."
	}

	hint := EmailHint{
		ID:           email.ID,
		RealDate:     email.RealDate,
		Subject:      email.Subject,
		Sender:       email.Sender,
		From:         email.From,
		To:           email.To,
		TextBodyHint: bodyHint,
		Attachments:  email.Attachments,
	}

	return hint
}
