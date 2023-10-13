package src

import (
	"io"
	"net/mail"
	"strings"
	"sync"
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
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"-"`
}

// ConvertAttachment converts parsemail.Attachment to a Attachment
func ConvertAttachment(at parsemail.Attachment) (Attachment, error) {
	data, err := io.ReadAll(at.Data)
	return Attachment{
		Filename:    at.Filename,
		ContentType: at.ContentType,
		Data:        data,
	}, err
}

// EmbeddedFile with json attributes
type EmbeddedFile struct {
	M *sync.Mutex `json:"-"`

	CID         string `json:"cid"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"-"`
}

// ConvertEmbeddedFile converts parsemail.EmbeddedFile to a EmbeddedFile
func ConvertEmbeddedFile(eb parsemail.EmbeddedFile) (EmbeddedFile, error) {
	data, err := io.ReadAll(eb.Data)
	return EmbeddedFile{
		M:           &sync.Mutex{},
		CID:         eb.CID,
		ContentType: eb.ContentType,
		Data:        data,
	}, err
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

	ContentType string `json:"contentType"`

	HTMLBody string `json:"htmlBody"`
	TextBody string `json:"textBody"`

	Attachments   []Attachment   `json:"attachments"`
	EmbeddedFiles []EmbeddedFile `json:"embeddedFiles"`
}

// ConvertEmail converts parsemail.Email to a Email
func ConvertEmail(bluemondayPolicy *bluemonday.Policy, ulid ulid.ULID, em parsemail.Email, realDate time.Time, realFrom, realTo string) (Email, error) {
	attachments, err := ConvertAttachmentList(em.Attachments)
	if err != nil {
		return Email{}, err
	}

	embeddedFiles, err := ConvertEmbeddedFileList(em.EmbeddedFiles)
	if err != nil {
		return Email{}, err
	}

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
		HTMLBody:        bluemondayPolicy.Sanitize(em.HTMLBody),
		TextBody:        em.TextBody,
		Attachments:     attachments,
		EmbeddedFiles:   embeddedFiles,
	}

	if resp.Sender == nil && realFrom != "" {
		resp.Sender = &Address{Address: realFrom}
	}

	if len(resp.To) == 0 && realTo != "" {
		resp.To = []*Address{{Address: realTo}}
	}

	return resp, nil
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
func ConvertAttachmentList(atList []parsemail.Attachment) ([]Attachment, error) {
	if atList == nil {
		return nil, nil
	}

	var err error
	results := make([]Attachment, len(atList))
	for idx, at := range atList {
		results[idx], err = ConvertAttachment(at)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

// ConvertEmbeddedFileList converts a list of parsemail.EmbeddedFile to a list of EmbeddedFile
func ConvertEmbeddedFileList(ebList []parsemail.EmbeddedFile) ([]EmbeddedFile, error) {
	if ebList == nil {
		return nil, nil
	}

	var err error
	results := make([]EmbeddedFile, len(ebList))
	for idx, eb := range ebList {
		results[idx], err = ConvertEmbeddedFile(eb)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
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
