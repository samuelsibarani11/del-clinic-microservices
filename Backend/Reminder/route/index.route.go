package route

import (
	"github.com/gofiber/fiber/v2"
	"reminder/handler"
	"reminder/middleware"
)

func RouteInit(r *fiber.App) {

	// REMINDER
	r.Get("/reminders", middleware.Auth, handler.ReminderGetAll)
	r.Post("/reminder", middleware.Auth, handler.CreateReminder)
	r.Put("/reminder/:id", middleware.Auth, handler.UpdateReminder)
	r.Delete("/reminder/:id", middleware.Auth, handler.DeleteReminder)

}
