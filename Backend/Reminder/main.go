package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"reminder/database"
	"reminder/database/migration"
	"reminder/route"
)

func main() {
	database.DatabaseInit()
	migration.Migration()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "https://gofiber.io",
	}))

	route.RouteInit(app)

	err := app.Listen("172.20.10.4:8007")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
