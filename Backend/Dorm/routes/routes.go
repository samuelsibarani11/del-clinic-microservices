package routes

import (
	"dorm/handler"
	"dorm/middleware"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(app *fiber.App) {
	// DORM
	app.Get("/dorm/:id", handler.DormGetById)
	app.Get("/dorms", middleware.StaffAuth, handler.DormGetAll)
	app.Post("/dorm", middleware.StaffAuth, handler.CreateDorm)
	app.Put("/dorm/:id", middleware.StaffAuth, handler.UpdateDorm)
	app.Delete("/dorm/:id", middleware.StaffAuth, handler.DeleteDorm)

}
