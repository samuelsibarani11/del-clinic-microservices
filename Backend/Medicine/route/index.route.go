package route

import (
	"github.com/gofiber/fiber/v2"
	"medicine/handler"
	"medicine/middleware"
)

func RouteInit(r *fiber.App) {

	// MEDICINE
	r.Get("/medicines", middleware.Auth, handler.MedicineGetAll)
	r.Post("/medicine", middleware.Auth, handler.CreateMedicine)
	r.Put("/medicine/:id", middleware.Auth, handler.UpdateMedicine)
	r.Delete("/medicine/:id", middleware.Auth, handler.DeleteMedicine)

}
