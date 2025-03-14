package src

import (
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"github.com/fasthttp/websocket"
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
	Credentials *EmailCredentials
	Entropy     *rand.Rand
	MaxEmails   int
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

	emailRawData, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	emailContents, err := parsemail.Parse(bytes.NewBuffer(emailRawData))
	if err != nil {
		return err
	}

	now := time.Now()
	ms := ulid.Timestamp(now)

	email, err := ConvertEmail(
		ulid.MustNew(ms, s.Backend.Entropy),
		emailRawData,
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
	if len(emails) > s.Backend.MaxEmails {
		emails = emails[len(emails)-s.Backend.MaxEmails:]
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

// StartEmailServerOptions are the options for starting the email server
type StartEmailServerOptions struct {
	Addr      string
	Domain    string
	Username  string
	Password  string
	MaxEmails uint16
	Tls       bool
}

// StartEmailServer starts the email server
func StartEmailServer(opts StartEmailServerOptions) {
	if opts.MaxEmails == 0 {
		opts.MaxEmails = 200
	}

	backend := &Backend{
		Entropy:   rand.New(rand.NewSource(time.Now().UnixNano())),
		MaxEmails: int(opts.MaxEmails),
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

	if opts.Tls {
		tlsConfig, err := generateTLSConfig()
		if err != nil {
			fmt.Println("WARN: failed to generate self signed tls certificate for email server, error:", err)
			os.Exit(1)
		} else {
			server.TLSConfig = tlsConfig
		}
	}

	fmt.Println("Running SMTP server at", server.Addr, "with a emails dequeue length of", opts.MaxEmails)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func generateTLSConfig() (*tls.Config, error) {
	// Generate a private key
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Set up certificate template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := cryptoRand.Int(cryptoRand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Self-Signed Cert"},
			CommonName:   "*",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create the certificate
	derBytes, err := x509.CreateCertificate(
		cryptoRand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		return nil, err
	}

	// Create the TLS certificate
	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
	}

	// Return the TLS config
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		// This is the key configuration - the server won't verify client hostnames
		InsecureSkipVerify: true,
		ClientAuth:         tls.NoClientCert,
		// Skip hostname verification completely - accept ANY hostname
		VerifyConnection: func(state tls.ConnectionState) error {
			return nil // Always approve the connection
		},
	}, nil
}
