package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"service/user/database"
	"service/user/database/migration"
	"service/user/routes"
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

	err := app.Listen("172.20.10.4:8008")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

}
