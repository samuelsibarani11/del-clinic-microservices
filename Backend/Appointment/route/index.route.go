package route

import (
	"appointment/handler"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {

	// APPOINTMENT
	r.Get("/appointments", handler.AppointmentGetAll)
	r.Get("/appointments-auth", handler.AppointmentGetByAuth)
	r.Post("/appointment", handler.CreateAppointment)
	r.Put("/appointment/:id/approve", handler.UpdateApprovedID)
	r.Put("/appointment/:id", handler.UpdateAppointment)
	r.Delete("/appointment/:id", handler.DeleteAppointment)

}
