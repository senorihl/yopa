package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const transPixel = "\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B"

func Setup(pixelCallback func(query string, remoteAddr string)) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/pixel.gif", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "image/gif")
		go pixelCallback(string(c.Request().URI().QueryString()), c.IP())
		return c.SendString(transPixel)
	})

	app.Post("/pixel.gif", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "image/gif")
		go pixelCallback(string(c.Request().URI().QueryString())+"&p="+string(c.Body()), c.IP())
		return c.SendString(transPixel)
	})

	app.Static("/assets", "/app/web/dist")

	return app
}
