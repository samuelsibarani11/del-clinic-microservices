package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"medicine/database"
	"medicine/model/entity"
	"medicine/utils"
	"os"
	"path/filepath"
)

var PathImageProduct = "./Public"

func init() {
	if _, err := os.Stat(PathImageProduct); os.IsNotExist(err) {
		err := os.Mkdir(PathImageProduct, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func MedicineGetAll(ctx *fiber.Ctx) error {
	var medicines []entity.Medicine

	database.DB.Find(&medicines)

	// Create a slice to store the response data
	response := make([]fiber.Map, len(medicines))

	// Iterate through medicines and populate response with required fields
	for i, medicine := range medicines {
		response[i] = fiber.Map{
			"id":      medicine.ID,
			"name":    medicine.Name,
			"amount":  medicine.Amount,
			"expired": medicine.Expired,
		}
	}

	return ctx.Status(200).JSON(fiber.Map{
		"medicine": response,
	})
}

func CreateMedicine(ctx *fiber.Ctx) error {
	// Parse form data
	input := new(entity.Medicine)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": "Invalid form data",
		})
	}

	// Validate the input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	// Process image
	image, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": "Image is required",
		})
	}

	filename := utils.GenerateImageFile(input.Name, image.Filename)
	if err := ctx.SaveFile(image, filepath.Join(PathImageProduct, filename)); err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "Can't save file image",
		})
	}
	input.Image = filename

	// Create the Medicine
	if err := database.DB.Create(input).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "failed to store data",
			"error":   err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   input,
	})
}

func UpdateMedicine(ctx *fiber.Ctx) error {
	// Parse form data
	input := new(entity.Medicine)
	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": "Invalid form data",
		})
	}

	// Get medicine ID from URL
	medicineID := ctx.Params("id")
	var medicine entity.Medicine

	// Find the existing medicine record
	if err := database.DB.First(&medicine, "id = ?", medicineID).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"status":  "failed",
			"message": "medicine not found",
		})
	}

	// Update fields
	if input.Name != "" {
		medicine.Name = input.Name
	}

	if input.Amount != 0 {
		medicine.Amount = input.Amount
	}

	if input.Expired != "" {
		medicine.Expired = input.Expired
	}

	// Process image if provided
	image, err := ctx.FormFile("image")
	if err == nil {
		filename := utils.GenerateImageFile(medicine.Name, image.Filename)
		if err := ctx.SaveFile(image, filepath.Join(PathImageProduct, filename)); err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"status":  "failed",
				"message": "Can't save file image",
			})
		}
		medicine.Image = filename
	}

	// Validate the struct
	validate := validator.New()
	if err := validate.Struct(medicine); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	// Update the medicine in the database
	if err := database.DB.Save(&medicine).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "failed to update data",
			"error":   err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   medicine,
	})
}

func DeleteMedicine(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var medicine entity.Medicine

	if err := database.DB.First(&medicine, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "medicine not found",
		})
	}

	if err := database.DB.Delete(&medicine).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete medicine",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Medicine deleted successfully",
	})
}
