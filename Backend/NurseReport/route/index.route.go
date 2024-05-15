package route

import (
	"github.com/gofiber/fiber/v2"
	"nurseReport/handler"
	"nurseReport/middleware"
)

func RouteInit(r *fiber.App) {
	// NURSE REPORT
	r.Get("/nurse-reports", middleware.StaffAuth, handler.NurseReportGetAll)
	r.Get("/nurse-report/:id", middleware.StaffAuth, handler.NurseReportGetById)
	r.Post("/nurse-report", middleware.StaffAuth, handler.CreateNurseReport)
	r.Put("/nurse-report/:id", middleware.StaffAuth, handler.UpdateNurseReport)
	r.Delete("/nurse-report/:id", middleware.StaffAuth, handler.DeleteNurseReport)

}
