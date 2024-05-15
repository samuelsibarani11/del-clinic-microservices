package handler

import (
	"dorm/Models/entity"
	"dorm/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func DormGetById(ctx *fiber.Ctx) error {
	// mencari user parameter id.
	dormId := ctx.Params("id")

	// mendeklarasikan variabel user dengan tipe data userEntity
	var dorm entity.Dorm

	// Query Statement dengan GORM
	err := database.DB.First(&dorm, "?", dormId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    dorm,
	})
}

func DormGetAll(ctx *fiber.Ctx) error {
	var dorms []entity.Dorm

	database.DB.Find(&dorms)

	// Create a slice to store the response data
	response := make([]fiber.Map, len(dorms))

	// Iterate through medicines and populate response with required fields
	for i, medicine := range dorms {
		response[i] = fiber.Map{
			"id":     medicine.ID,
			"name":   medicine.Name,
			"status": medicine.Status,
		}
	}

	return ctx.Status(200).JSON(fiber.Map{
		"medicine": response,
	})
}

func CreateDorm(ctx *fiber.Ctx) error {
	dorm := &entity.Dorm{}

	if err := ctx.BodyParser(dorm); err != nil {
		return err
	}

	// VALIDATION Request
	validate := validator.New()
	if err := validate.Struct(dorm); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	if err := database.DB.Create(dorm).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(),
		})
	}
	response := entity.Dorm{
		ID:     dorm.ID,
		Name:   dorm.Name,
		Status: dorm.Status,
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}

func UpdateDorm(ctx *fiber.Ctx) error {

	dormRequest := new(entity.Dorm)

	if err := ctx.BodyParser(dormRequest); err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "bad request",
		})
	}

	var dorm entity.Dorm

	dormID := ctx.Params("id")

	err := database.DB.First(&dorm, "id = ?", dormID).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "medicine not found",
		})
	}

	// UPDATE USER DATA
	if dorm.Name != "" {
		dorm.Name = dormRequest.Name
	}

	if dorm.Name != "" {
		dorm.Status = dormRequest.Status
	}

	errUpdate := database.DB.Save(&dorm).Error

	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "internal server error",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    dorm,
	})

}

func DeleteDorm(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var dorm entity.Dorm

	if err := database.DB.First(&dorm, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "medicine not found",
		})
	}

	if err := database.DB.Delete(&dorm).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete medicine",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Medicine deleted successfully",
	})
}
