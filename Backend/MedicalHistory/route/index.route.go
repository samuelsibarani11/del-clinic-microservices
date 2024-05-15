package route

import (
	"github.com/gofiber/fiber/v2"
	"medicalHistory/handler"
	"medicalHistory/middleware"
)

func RouteInit(r *fiber.App) {

	// MEDICAL HISTORY
	r.Get("/medical-histories", middleware.Auth, handler.GetAllMedicalHistoryByToken)
	r.Post("/medical-history", middleware.StaffAuth, handler.CreateMedicalHistory)

}
