package src

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"github.com/fasthttp/websocket"
	"github.com/microcosm-cc/bluemonday"
	"github.com/oklog/ulid/v2"
)

var emailsLock sync.Mutex
var emails = []Email{}

type EmailCredentials struct {
	Username string
	Password string
}

// The Backend implements SMTP server methods.
type Backend struct {
	Credentials      *EmailCredentials
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

// unauthenticated returns true if the session is not authenticated
func (s *Session) unauthenticated() bool {
	return s.Backend.Credentials != nil && !s.Auth
}

// AuthPlain implements smtp.Session
func (s *Session) AuthPlain(username, password string) error {
	creds := s.Backend.Credentials
	if creds == nil {
		return nil
	}

	if username != creds.Username || password != creds.Password {
		return smtp.ErrAuthFailed
	}

	s.Auth = true
	return nil
}

// Mail implements smtp.Session
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if s.unauthenticated() {
		return smtp.ErrAuthRequired
	}

	s.From = from
	return nil
}

// Rcpt implements smtp.Session
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if s.unauthenticated() {
		return smtp.ErrAuthRequired
	}

	s.To = to
	return nil
}

// Data implements smtp.Session
func (s *Session) Data(r io.Reader) error {
	if s.unauthenticated() {
		return smtp.ErrAuthRequired
	}

	emailContents, err := parsemail.Parse(r)
	if err != nil {
		return err
	}

	now := time.Now()
	ms := ulid.Timestamp(now)

	email, err := ConvertEmail(
		s.Backend.BluemondayPolicy,
		ulid.MustNew(ms, s.Backend.Entropy),
		emailContents,
		now,
		s.From,
		s.To,
	)
	if err != nil {
		return err
	}

	emailsLock.Lock()
	emails = append(emails, email)
	if len(emails) > 100 {
		emails = emails[len(emails)-100:]
	}
	emailsLock.Unlock()

	for _, c := range websocketConnections {
		c.WriteMessage(websocket.TextMessage, []byte("new-email"))
	}

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

type StartEmailServerOptions struct {
	Addr     string
	Domain   string
	Username string
	Password string
}

// StartEmailServer starts the email server
func StartEmailServer(opts StartEmailServerOptions) {
	backend := &Backend{
		Entropy:          rand.New(rand.NewSource(time.Now().UnixNano())),
		BluemondayPolicy: bluemonday.UGCPolicy(),
	}

	if opts.Username != "" && opts.Password != "" {
		backend.Credentials = &EmailCredentials{
			Username: opts.Username,
			Password: opts.Password,
		}
	}

	server := smtp.NewServer(backend)

	server.Addr = opts.Addr
	server.Domain = opts.Domain
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024 * 10 // 10 MB
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true
	server.AuthDisabled = backend.Credentials == nil

	fmt.Println("Running SMTP server at", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
