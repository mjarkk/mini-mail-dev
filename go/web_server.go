package src

import (
	"errors"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

// ErrorResposne is the response send by the server when an error occurs
type ErrorResposne struct {
	Error string `json:"error"`
}

// StartWebserver starts the webserver
func StartWebserver() {
	app := fiber.New()

	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(logger.New())

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}

		return c.JSON(ErrorResposne{
			Error: err.Error(),
		})
	})

	apiGroup := app.Group("/api")

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
		hints := make([]EmailHint, len(emails))

		for idx, email := range emails {
			hints[len(hints)-idx-1] = EmailToEmailHint(email)
		}

		return c.JSON(hints)
	})

	apiGroup.Get("/emails/:id", func(c *fiber.Ctx) error {
		id, err := ulid.Parse(c.Params("id"))
		if err != nil {
			return err
		}

		for _, email := range emails {
			if email.ID == id {
				return c.JSON(email)
			}
		}

		return errors.New("email not found")
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
		id, err := ulid.Parse(c.Params("id"))
		if err != nil {
			return err
		}

		var invalidNr = errors.New("invalid attachment number")

		nr, err := c.ParamsInt("nr")
		if err != nil {
			return invalidNr
		}

		if nr < 0 {
			return invalidNr
		}

		var paramEmail *Email
		for _, email := range emails {
			if email.ID == id {
				paramEmail = &email
				break
			}

		}

		if paramEmail == nil {
			return errors.New("email not found")
		}

		if nr >= len(paramEmail.Attachments) {
			return invalidNr
		}

		attachment := paramEmail.Attachments[nr]
		c.Response().Header.SetContentType(attachment.ContentType)
		c.Response().Header.Add("Content-Disposition", "attachment; filename=\""+attachment.Filename+"\"")
		return c.Send(attachment.Data)
	})

	apiGroup.Get("/emails/:id/embeddedFiles/:nr", func(c *fiber.Ctx) error {
		id, err := ulid.Parse(c.Params("id"))
		if err != nil {
			return err
		}

		var invalidNr = errors.New("invalid embedded file number")

		nr, err := c.ParamsInt("nr")
		if err != nil {
			return invalidNr
		}

		if nr < 0 {
			return invalidNr
		}

		var paramEmail *Email
		for _, email := range emails {
			if email.ID == id {
				paramEmail = &email
				break
			}

		}

		if paramEmail == nil {
			return errors.New("email not found")
		}

		if nr >= len(paramEmail.EmbeddedFiles) {
			return invalidNr
		}

		attachment := paramEmail.EmbeddedFiles[nr]
		c.Response().Header.SetContentType(attachment.ContentType)
		return c.Send(attachment.Data)
	})

	log.Fatal(app.Listen(":3000"))
}
