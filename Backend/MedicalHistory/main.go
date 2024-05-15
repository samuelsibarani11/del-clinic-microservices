package main

import (
	"github.com/gofiber/fiber/v2"
	"medicalHistory/database"
	"medicalHistory/database/migration"
	"medicalHistory/route"
)

func main() {

	// INITIAL DATABASE
	database.DatabaseInit()

	// RUN MIGRATION
	migration.Migration()

	// menginisialisasikan go fiber (di passing ke route)
	app := fiber.New()

	// INITIAL ROUTE
	route.RouteInit(app)

	app.Listen(":8080")
}
