package route

import (
	"github.com/gofiber/fiber/v2"
	"medicine/handler"
)

func RouteInit(r *fiber.App) {

	// MEDICINE
	r.Get("/medicines", handler.MedicineGetAll)
	r.Post("/medicine", handler.CreateMedicine)
	r.Put("/medicine/:id", handler.UpdateMedicine)
	r.Delete("/medicine/:id", handler.DeleteMedicine)

}
