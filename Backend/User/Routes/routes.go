package routes

import (
	"github.com/gofiber/fiber/v2"
	"user/handler"
	"user/middleware"
)

func RouteInit(app *fiber.App) {
	// USER
	app.Get("/users", middleware.Auth, handler.UserHandlerGetAll)
	app.Get("/user/:id", middleware.Auth, handler.UserHandlerGetById)
	app.Post("/user", handler.CreateUser)
	app.Put("/user/:id", middleware.Auth, handler.UpdateUser)
	app.Delete("/user/:id", handler.DeleteUser)
}
