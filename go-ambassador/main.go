package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sing3demons/ambassador/src/database"
	"github.com/sing3demons/ambassador/src/routes"
)

func main() {
	database.Connect()
	database.AutoMigrate()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	routes.Setup(app)

	log.Fatal(app.Listen(":8000"))
}
