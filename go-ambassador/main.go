package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/sing3demons/ambassador/src/database"
	"github.com/sing3demons/ambassador/src/routes"
)

var (
	production = "production"
)

func init() {
	if os.Getenv("APP_ENV") != production {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	if os.Getenv("APP_ENV") != "seed" {
		database.Connect()
		database.AutoMigrate()
	}
}

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
