package src

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	if opts.BasicAuthUsername != "" || opts.BasicAuthPassword != "" {
		app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				opts.BasicAuthUsername: opts.BasicAuthPassword,
			},
		}))
	}

	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(logger.New())

	apiGroup := app.Group("/api")

	apiGroup.Use(func(c *fiber.Ctx) error {
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

	apiGroup.Get("/emails", func(c *fiber.Ctx) error {
		searchValue := c.Query("search")

		emailsCopy := []*Email{}

		for idx := 0; idx < len(emails); idx++ {
			if len(searchValue) > 0 && !filterEmail(searchValue, emails[idx]) {
				continue
			}

			emailsCopy = append(emailsCopy, &emails[idx])
		}

		return c.JSON(emailsCopy)
	})

	apiGroup.Get("/emails/:id/remainder", func(c *fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}
		return c.JSON(email.Remainder)
	})

	apiGroup.Get("/emails/:id/page", func(c *fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}

		c.Response().SetBodyRaw([]byte(email.Remainder.HTMLBody))
		c.Response().Header.SetContentType(fiber.MIMETextHTML)
		return nil
	})

	apiGroup.Delete("/emails/:id", func(c *fiber.Ctx) error {
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

	apiGroup.Get("/emails/:id/attachments/:nr", func(c *fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}

		var invalidNr = errors.New("invalid attachment number")

		nr, err := c.ParamsInt("nr")
		if err != nil || nr < 0 || nr >= len(email.Remainder.Attachments) {
			return invalidNr
		}

		attachment := email.Remainder.Attachments[nr]
		c.Response().Header.SetContentType(attachment.ContentType)
		c.Response().Header.Add("Content-Disposition", "attachment; filename=\""+attachment.Filename+"\"")
		return c.Send(attachment.Data)
	})

	apiGroup.Get("/emails/:id/embeddedFiles/:nr", func(c *fiber.Ctx) error {
		email, err := findEmail(c.Params("id"))
		if err != nil {
			return err
		}

		var invalidNr = errors.New("invalid embedded file number")

		nr, err := c.ParamsInt("nr")
		if err != nil || nr < 0 || nr >= len(email.Remainder.EmbeddedFiles) {
			return invalidNr
		}

		attachment := email.Remainder.EmbeddedFiles[nr]
		c.Response().Header.SetContentType(attachment.ContentType)
		return c.Send(attachment.Data)
	})

	app.Use(filesystem.New(filesystem.Config{
		PathPrefix: "dist",
		Index:      "index.html",
		Root:       http.FS(dist),
	}))

	fmt.Println("Running Web server at", opts.Addr)
	log.Fatal(app.Listen(opts.Addr))
}
