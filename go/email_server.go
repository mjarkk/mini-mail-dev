package src

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/oklog/ulid/v2"
)

var emailsLock sync.Mutex
var emails = []Email{}

// The Backend implements SMTP server methods.
type Backend struct {
	Entropy          *rand.Rand
	BluemondayPolicy *bluemonday.Policy
}

// NewSession implements smtp.Backend
func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{
		Backend: b,
	}, nil
}

// A Session is returned after EHLO.
type Session struct {
	Backend *Backend
	From    string
	To      string
	Auth    bool
}

// AuthPlain implements smtp.Session
func (s *Session) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return smtp.ErrAuthFailed
	}

	s.Auth = true
	return nil
}

// Mail implements smtp.Session
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.Auth {
		return smtp.ErrAuthRequired
	}

	s.From = from
	return nil
}

// Rcpt implements smtp.Session
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.Auth {
		return smtp.ErrAuthRequired
	}

	s.To = to
	return nil
}

// Data implements smtp.Session
func (s *Session) Data(r io.Reader) error {
	if !s.Auth {
		return smtp.ErrAuthRequired
	}

	emailContents, err := parsemail.Parse(r)
	if err != nil {
		return err
	}

	now := time.Now()
	ms := ulid.Timestamp(now)

	email := ConvertEmail(
		s.Backend.BluemondayPolicy,
		ulid.MustNew(ms, s.Backend.Entropy),
		emailContents,
		now,
		s.From,
		s.To,
	)

	emailsLock.Lock()
	emails = append(emails, email)
	emailsLock.Unlock()

	return nil
}

// Reset implements smtp.Session
func (s *Session) Reset() {
	s.From = ""
	s.To = ""
}

// Logout implements smtp.Session
func (s *Session) Logout() error {
	return nil
}

// StartEmailServer starts the email server
func StartEmailServer() {
	backend := &Backend{
		Entropy:          rand.New(rand.NewSource(time.Now().UnixNano())),
		BluemondayPolicy: bluemonday.UGCPolicy(),
	}

	server := smtp.NewServer(backend)

	server.Addr = ":1025"
	server.Domain = "localhost"
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024 * 10 // 10 MB
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	log.Println("Starting server at", server.Domain+server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
