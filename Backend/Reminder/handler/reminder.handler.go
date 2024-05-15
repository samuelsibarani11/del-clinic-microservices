package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"reminder/database"
	"reminder/model/entity"
)

func ReminderGetAll(ctx *fiber.Ctx) error {
	var reminders []entity.Reminder

	database.DB.Find(&reminders)

	// Create a slice to store the response data
	response := make([]fiber.Map, len(reminders))

	// Iterate through medicines and populate response with required fields
	for i, reminder := range reminders {
		response[i] = fiber.Map{
			"id":          reminder.ID,
			"date":        reminder.Date,
			"time":        reminder.Time,
			"frequency":   reminder.Frequency,
			"duration":    reminder.Frequency,
			"medicine_id": reminder.MedicineID,
		}
	}

	return ctx.Status(200).JSON(fiber.Map{
		"reminder": response,
	})
}

func CreateReminder(ctx *fiber.Ctx) error {
	reminder := &entity.Reminder{}

	if err := ctx.BodyParser(reminder); err != nil {
		return err
	}

	// VALIDATION Request
	validate := validator.New()
	if err := validate.Struct(reminder); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	if err := database.DB.Create(reminder).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(),
		})
	}
	response := entity.Reminder{
		ID:         reminder.ID,
		Date:       reminder.Time,
		Frequency:  reminder.Frequency,
		Duration:   reminder.Duration,
		MedicineID: reminder.MedicineID,
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}

func UpdateReminder(ctx *fiber.Ctx) error {

	reminderRequest := new(entity.Reminder)

	if err := ctx.BodyParser(reminderRequest); err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "bad request",
		})
	}

	var reminder entity.Reminder

	medicineID := ctx.Params("id")

	err := database.DB.First(&reminder, "id = ?", medicineID).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "reminder not found",
		})
	}

	// UPDATE USER DATA
	if reminder.Date != "" {
		reminder.Date = reminderRequest.Date
	}

	if reminder.Time != "" {
		reminder.Time = reminderRequest.Time
	}

	if reminder.Frequency != "" {
		reminder.Frequency = reminderRequest.Frequency
	}

	if reminder.Duration != "" {
		reminder.Duration = reminderRequest.Duration
	}

	if reminder.MedicineID != 0 {
		reminder.MedicineID = reminderRequest.MedicineID
	}

	errUpdate := database.DB.Save(&reminder).Error

	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "internal server error",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    reminder,
	})

}

func DeleteReminder(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var reminder entity.Reminder

	if err := database.DB.First(&reminder, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "reminder not found",
		})
	}

	if err := database.DB.Delete(&reminder).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete reminder",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "reminder deleted successfully",
	})
}
