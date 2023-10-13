package src

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/oklog/ulid/v2"
)

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
	emailsGroup := apiGroup.Group("/emails")

	emailsGroup.Get("", func(c *fiber.Ctx) error {
		hints := []EmailHint{}

		emailsLock.Lock()
		for _, email := range emails {
			hints = append(hints, EmailToEmailHint(email))
		}
		emailsLock.Unlock()

		return c.JSON(hints)
	})

	emailsGroup.Get("/:id", func(c *fiber.Ctx) error {
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

	emailsGroup.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := ulid.Parse(c.Params("id"))
		if err != nil {
			return err
		}

		emailsLock.Lock()
		for i, email := range emails {
			if email.ID == id {
				emails = append(emails[:i], emails[i+1:]...)
				break
			}
		}
		emailsLock.Unlock()

		return c.SendStatus(fiber.StatusNoContent)
	})

	log.Fatal(app.Listen(":3000"))
}
