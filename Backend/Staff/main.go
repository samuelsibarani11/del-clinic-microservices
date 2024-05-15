package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"staff/database"
	"staff/database/migration"
	"staff/routes"
)

func main() {
	database.Connect()
	migration.Migration()
	// seeders.SeederData()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "https://gofiber.io",
	}))

	routes.RouteInit(app)

	err := app.Listen("192.168.82.133:8004")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

}
