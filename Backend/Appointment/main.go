package main

import (
	"appointment/database"
	"appointment/database/migration"
	"appointment/route"
	"github.com/gofiber/fiber/v2"
)

func main() {

	//INITIAL DATABASE
	database.DatabaseInit()

	// RUN MIGRATION
	migration.Migration()

	// menginisialisasikan go fiber (di passing ke route)
	app := fiber.New()

	// INITIAL ROUTE
	route.RouteInit(app)

	app.Listen(":8080")
}
