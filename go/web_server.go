package src

import (
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
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

var websocketConnectionsLock sync.Mutex
var websocketConnections []*websocket.Conn

func registerWebsocketConnection(c *websocket.Conn) {
	websocketConnectionsLock.Lock()
	defer websocketConnectionsLock.Unlock()
	websocketConnections = append(websocketConnections, c)
}
func unregisterWebsocketConnection(c *websocket.Conn) bool {
	websocketConnectionsLock.Lock()
	defer websocketConnectionsLock.Unlock()
	for _, websocketConnection := range websocketConnections {
		if websocketConnection == c {
			websocketConnections = websocketConnections[:len(websocketConnections)-1]
			return true
		}
	}
	return false
}

func filterSanitize(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func comparitorFn(input string) func(s string) bool {
	input = filterSanitize(input)
	if input == "" {
		return func(s string) bool {
			return true
		}
	}

	return func(s string) bool {
		return strings.Contains(filterSanitize(s), input)
	}
}

func filterEmail(searchValue string, email Email) bool {
	comparitor := comparitorFn(searchValue)

	if comparitor(email.Subject) {
		return true
	}

	if comparitor(email.Sender.Address) {
		return true
	}

	if comparitor(email.Sender.Name) {
		return true
	}

	for _, ad := range email.To {
		if comparitor(ad.Address) {
			return true
		}

		if comparitor(ad.Name) {
			return true
		}
	}

	if email.BodyType == "text" {
		return comparitor(email.Remainder.TextBody)
	}

	return comparitor(email.Remainder.HTMLBody)
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
		registerWebsocketConnection(c)

		var err error
		for {
			_, _, err = c.ReadMessage()
			if err != nil {
				break
			}
		}

		if !unregisterWebsocketConnection(c) {
			log.Println("Was unable to close websocket connection, this should not happen and is a bug")
		}
	}))

	apiGroup.Get("/emails", func(c fiber.Ctx) error {
		searchValue := strings.TrimSpace(c.Query("search"))

		// Held through c.JSON: the response points into emails.
		emailsLock.Lock()
		defer emailsLock.Unlock()

		if searchValue == "" {
			// Take the quick path

			response := make([]*Email, len(emails))

			for idx := len(emails) - 1; idx >= 0; idx-- {
				response[len(emails)-idx-1] = &emails[idx]
			}

			return c.JSON(response)
		}

		response := []*Email{}
		for idx := len(emails) - 1; idx >= 0; idx-- {
			if len(searchValue) > 0 && !filterEmail(searchValue, emails[idx]) {
				continue
			}

			response = append(response, &emails[idx])
		}

		return c.JSON(response)
	})

	apiGroup.Get("/emails/:id/remainder", func(c fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}
		return c.JSON(email.Remainder)
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

	app.Get("/*", static.New("dist", static.Config{FS: dist}))

	fmt.Println("Running Web server at", opts.Addr)
	log.Fatal(app.Listen(opts.Addr, fiber.ListenConfig{DisableStartupMessage: true}))
}
