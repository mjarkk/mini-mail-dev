package src

import (
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"mime"
	"strings"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/mjarkk/mini-mail-dev/go/sanitize"
	"github.com/oklog/ulid/v2"
)

var newEmailSubscribersLock sync.Mutex
var newEmailSubscribers = map[chan struct{}]bool{}

func subscribeToNewEmails() chan struct{} {
	ch := make(chan struct{}, 1)
	newEmailSubscribersLock.Lock()
	defer newEmailSubscribersLock.Unlock()
	newEmailSubscribers[ch] = true
	return ch
}

func unsubscribeFromNewEmails(ch chan struct{}) {
	newEmailSubscribersLock.Lock()
	defer newEmailSubscribersLock.Unlock()
	delete(newEmailSubscribers, ch)
}

// notifyNewEmail never blocks, so a stalled browser tab can't hold up an SMTP delivery.
func notifyNewEmail() {
	newEmailSubscribersLock.Lock()
	defer newEmailSubscribersLock.Unlock()
	for ch := range newEmailSubscribers {
		select {
		case ch <- struct{}{}:
		default:
			// One pending notification is enough: the tab refetches the whole list.
		}
	}
}

// searchEmails returns pointers into list, newest first.
func searchEmails(list []Email, search string) []*Email {
	search = strings.ToLower(strings.TrimSpace(search))
	result := make([]*Email, 0, len(list))
	for i := len(list) - 1; i >= 0; i-- {
		if strings.Contains(list[i].searchText, search) {
			result = append(result, &list[i])
		}
	}
	return result
}

// ErrorResponse is the response send by the server when an error occurs
type ErrorResponse struct {
	Error string `json:"error"`
}

// findEmail finds an email by its id
func findEmail(id string) (Email, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return Email{}, err
	}

	emailsLock.Lock()
	defer emailsLock.Unlock()

	for _, email := range emails {
		if email.ID == parsedID {
			return email, nil
		}
	}

	return Email{}, errors.New("email not found")
}

// StartWebserverOptions are the options for the webserver
type StartWebserverOptions struct {
	Addr              string
	BasicAuthUsername string
	BasicAuthPassword string
}

// StartWebserver starts the webserver
func StartWebserver(dist embed.FS, opts StartWebserverOptions) {
	app := fiber.New()

	if opts.BasicAuthUsername != "" || opts.BasicAuthPassword != "" {
		// basicauth only accepts hashed passwords.
		passwordHash := sha256.Sum256([]byte(opts.BasicAuthPassword))
		app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				opts.BasicAuthUsername: "{SHA256}" + base64.StdEncoding.EncodeToString(passwordHash[:]),
			},
		}))
	}

	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(logger.New())

	apiGroup := app.Group("/api")

	apiGroup.Use(func(c fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}

		return c.Status(400).JSON(ErrorResponse{
			Error: err.Error(),
		})
	})

	apiGroup.Get("/emails-events", websocket.New(func(c *websocket.Conn) {
		// Returning leaves a hijacked connection open.
		defer c.Close()

		newEmail := subscribeToNewEmails()
		defer unsubscribeFromNewEmails(newEmail)

		// The tab never sends anything; reading is how a closed tab is noticed.
		closed := make(chan struct{})
		go func() {
			defer close(closed)
			for {
				if _, _, err := c.ReadMessage(); err != nil {
					return
				}
			}
		}()

		for {
			select {
			case <-newEmail:
				if c.WriteMessage(websocket.TextMessage, []byte("new-email")) != nil {
					return
				}
			case <-closed:
				return
			}
		}
	}))

	apiGroup.Get("/emails", func(c fiber.Ctx) error {
		// Held through c.JSON: the response points into emails.
		emailsLock.Lock()
		defer emailsLock.Unlock()

		return c.JSON(searchEmails(emails, c.Query("search")))
	})

	apiGroup.Get("/emails/:id/remainder", func(c fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}
		return c.JSON(email.Remainder)
	})

	apiGroup.Get("/emails/:id/raw", func(c fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}
		c.Response().Header.SetContentType(fiber.MIMETextPlainCharsetUTF8)
		return c.Send(email.Remainder.Raw)
	})

	apiGroup.Get("/email-clients", func(c fiber.Ctx) error {
		return c.JSON(sanitize.Clients)
	})

	apiGroup.Get("/emails/:id/page", func(c fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}

		client := sanitize.Clients[0]
		if id := c.Query("client"); id != "" {
			client = sanitize.ClientByID(id)
			if client == nil {
				return fmt.Errorf("unknown email client %q", id)
			}
		}

		page, err := sanitize.Parse(email.Remainder.HTMLBody, client)
		if err != nil {
			return err
		}

		c.Response().SetBodyRaw([]byte(page))
		// The body was converted to UTF-8 while parsing; this header overrides any <meta charset> left in the HTML.
		c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		c.Response().Header.Set("Content-Security-Policy", client.ContentSecurityPolicy())
		return nil
	})

	apiGroup.Post("/emails/:id/read", func(c fiber.Ctx) error {
		id, err := ulid.Parse(c.Params("id"))
		if err != nil {
			return err
		}

		emailsLock.Lock()
		defer emailsLock.Unlock()

		for i := range emails {
			if emails[i].ID == id {
				emails[i].Read = true
				return c.SendStatus(fiber.StatusNoContent)
			}
		}

		return errors.New("email not found")
	})

	apiGroup.Delete("/emails/:id", func(c fiber.Ctx) error {
		id, err := ulid.Parse(c.Params("id"))
		if err != nil {
			return err
		}

		emailsLock.Lock()
		defer emailsLock.Unlock()

		for i, email := range emails {
			if email.ID == id {
				emails = append(emails[:i], emails[i+1:]...)
				break
			}
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	apiGroup.Get("/emails/:id/attachments/:nr", func(c fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}

		var invalidNr = errors.New("invalid attachment number")

		nr := fiber.Params(c, "nr", -1)
		if nr < 0 || nr >= len(email.Remainder.Attachments) {
			return invalidNr
		}

		attachment := email.Remainder.Attachments[nr]
		c.Response().Header.SetContentType(attachment.ContentType)
		c.Response().Header.Add("Content-Disposition", "attachment; filename=\""+attachment.Filename+"\"")
		return c.Send(attachment.Data)
	})

	apiGroup.Get("/emails/:id/embeddedFiles/:nr", func(c fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}

		var invalidNr = errors.New("invalid embedded file number")

		nr := fiber.Params(c, "nr", -1)
		if nr < 0 || nr >= len(email.Remainder.EmbeddedFiles) {
			return invalidNr
		}

		attachment := email.Remainder.EmbeddedFiles[nr]
		c.Response().Header.SetContentType(attachment.ContentType)
		return c.Send(attachment.Data)
	})

	// Go's built-in table lacks this type, and the scratch image has no /etc/mime.types to fall back on.
	mime.AddExtensionType(".webmanifest", "application/manifest+json")
	app.Get("/*", static.New("dist", static.Config{FS: dist}))

	fmt.Println("Running Web server at", opts.Addr)
	log.Fatal(app.Listen(opts.Addr, fiber.ListenConfig{DisableStartupMessage: true}))
}
