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

type EmailRemainder struct {
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
	BodyHint string    `json:"bodyHint"`

	// Recived form the client
	Sender *Address   `json:"sender"`
	From   []*Address `json:"from"`
	To     []*Address `json:"to"`

	Remainder EmailRemainder `json:"-"`
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

	bodyHint := strings.NewReplacer(
		"\n", " ",
		"\r", " ",
		"\t", " ",
		"   ", " ",
		"  ", " ",
	).Replace(em.TextBody)
	if len(bodyHint) > 100 {
		bodyHint = bodyHint[:100] + "..."
	}

	resp := Email{
		ID:       ulid,
		RealDate: realDate,
		Subject:  em.Subject,
		Sender:   MightConvertAddress(em.Sender),
		From:     ConvertAddressList(em.From),
		To:       ConvertAddressList(em.To),
		BodyHint: bodyHint,
		Remainder: EmailRemainder{
			Header:          em.Header,
			ReplyTo:         ConvertAddressList(em.ReplyTo),
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
		},
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
