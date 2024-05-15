package route

import (
	"appointment/handler"
	"appointment/middleware"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	// APPOINTMENT
	r.Get("/appointments", handler.AppointmentGetAll)
	r.Get("/appointments-auth", middleware.Auth, handler.AppointmentGetByAuth)
	r.Post("/appointment", middleware.Auth, handler.CreateAppointment)
	r.Put("/appointment/:id/approve", middleware.StaffAuth, handler.UpdateApprovedID)
	r.Put("/appointment/:id", middleware.Auth, handler.UpdateAppointment)
	r.Delete("/appointment/:id", middleware.Auth, handler.DeleteAppointment)

}
