package routes

import (
	"github.com/gofiber/fiber/v2"
	"staff/handler"
	"staff/middleware"
)

func RouteInit(app *fiber.App) {
	// AUTH
	app.Post("/login", handler.StaffLoginHandler)
	app.Get("/profile", handler.GetProfile)
	app.Post("/logout", handler.StaffLogoutHandler)

	// CREATE STAFF
	app.Get("/staffs", middleware.StaffAuth, handler.StaffHandlerGetAll)
	app.Get("/staff/:id", middleware.StaffAuth, handler.StaffHandlerGetById)
	app.Post("/staff", handler.CreateStaff)
	app.Put("/staff/:id", middleware.StaffAuth, handler.UpdateStaff)
	app.Delete("/staff/:id", middleware.StaffAuth, handler.DeleteStaff)

}
