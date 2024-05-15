package main

import (
	"dorm/database"
	"dorm/database/migration"
	"dorm/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
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

	err := app.Listen("172.20.10.4:8001")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

}
