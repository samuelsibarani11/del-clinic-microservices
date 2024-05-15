package route

import (
	"doctorReport/handler"
	"doctorReport/middleware"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {

	// DOCTOR REPORT
	r.Get("/doctor-reports", middleware.StaffAuth, handler.DoctorReportGetAll)
	r.Post("/doctor-report", middleware.StaffAuth, handler.CreateDoctorReport)
	r.Put("/doctor-report/:id", middleware.StaffAuth, handler.UpdateDoctorReport)
	r.Delete("/doctor-report/:id", middleware.StaffAuth, handler.DeleteDoctorReport)

}
