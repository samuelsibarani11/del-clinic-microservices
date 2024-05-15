package handler

import (
	"appointment/database"
	"appointment/model/entity"
	"appointment/utils"
	"github.com/gofiber/fiber/v2"

	"log"
)

func AppointmentGetByAuth(ctx *fiber.Ctx) error {
	// Mendapatkan token dari header Authorization
	token := ctx.Get("Authorization")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Mendekode token untuk mendapatkan informasi pengguna
	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated err",
		})
	}

	// Mendapatkan ID pengguna dari token
	userID := int(claims["id"].(float64))

	// Melakukan query ke basis data untuk mendapatkan appointment yang sesuai dengan ID pengguna
	// Anda perlu menyesuaikan sesuai dengan struktur tabel dan kueri yang benar
	var appointments []entity.Appointment
	err = database.DB.Where("requested_id = ?", userID).Find(&appointments).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "appointment not found",
		})
	}

	// Membuat slice baru untuk menyimpan data appointment yang akan dikirim ke client
	var responseAppointments []interface{}
	for _, appointment := range appointments {
		responseAppointment := map[string]interface{}{
			"id":        appointment.ID,
			"date":      appointment.Date,
			"time":      appointment.Time,
			"complaint": appointment.Complaint,
			"approved":  appointment.Approved,
		}
		responseAppointments = append(responseAppointments, responseAppointment)
	}

	// Mengembalikan data appointment dalam format JSON
	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    responseAppointments,
	})
}

func AppointmentGetAll(ctx *fiber.Ctx) error {
	var appointments []entity.Appointment

	// Memuat entitas terkait menggunakan Preload
	result := database.DB.Preload("Approved").Preload("Requested").Find(&appointments)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch appointments",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success":      "get data success",
		"appointments": appointments,
	})
}

func CreateAppointment(ctx *fiber.Ctx) error {
	appointment := new(entity.AppointmentResponse)

	//PARSE TO OBJECT STRUCT
	if err := ctx.BodyParser(appointment); err != nil {
		return ctx.Status(503).JSON(fiber.Map{
			"err": err,
		})
	}

	// VALIDATION
	//log.Println(appointment.RequestedID)
	if appointment.RequestedID == 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"err": "requested_id is required",
		})
	}

	if err := database.DB.Create(&appointment).Error; err != nil {
		// Mengembalikan respon error 500 dengan pesan yang sesuai
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(), // Menambahkan pesan error ke respon JSON
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message":     "create data successfully",
		"appointment": appointment,
	})
}

func UpdateApprovedID(ctx *fiber.Ctx) error {
	// Get the ID from the URL
	id := ctx.Params("id")

	// Find the appointment in the database
	appointment := new(entity.AppointmentResponse)
	if err := database.DB.First(&appointment, id).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "appointment not found",
			"error":   err.Error(),
		})
	}

	// Parse the updated data
	updatedAppointment := new(entity.AppointmentResponse)
	if err := ctx.BodyParser(updatedAppointment); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "error parsing data",
			"error":   err.Error(),
		})
	}

	// Validate the parsed data
	if updatedAppointment.ApprovedID == nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "approved_id is required",
		})
	}

	// Update the appointment in the database
	appointment.ApprovedID = updatedAppointment.ApprovedID

	if err := database.DB.Save(&appointment).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to update data",
			"error":   err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message":     "update data successfully",
		"appointment": appointment,
	})
}

func UpdateAppointment(ctx *fiber.Ctx) error {
	appointmentRequest := new(entity.AppointmentUpdate)
	if err := ctx.BodyParser(appointmentRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "Bad request"})
	}

	var appointment entity.Appointment

	appointmentID := ctx.Params("id")
	// CHECK AVAILABLE APPOINTMENT
	err := database.DB.First(&appointment, "id = ?", appointmentID).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Appointment not found",
		})
	}

	// UPDATE APPOINTMENT DATA
	if appointmentRequest.Date != "" {
		appointment.Date = appointmentRequest.Date
	}
	if appointmentRequest.Time != "" {
		appointment.Time = appointmentRequest.Time
	}
	if appointmentRequest.Complaint != "" {
		appointment.Complaint = appointmentRequest.Complaint
	}

	errUpdate := database.DB.Model(&appointment).Updates(entity.Appointment{
		Date:      appointment.Date,
		Time:      appointment.Time,
		Complaint: appointment.Complaint,
	}).Error

	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    appointment,
	})
}

func DeleteAppointment(c *fiber.Ctx) error {
	id := c.Params("id")

	var appointment entity.Appointment
	if err := database.DB.First(&appointment, id).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Appointment not found",
		})
	}

	if err := database.DB.Delete(&appointment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Appointment deleted successfully",
	})
}
