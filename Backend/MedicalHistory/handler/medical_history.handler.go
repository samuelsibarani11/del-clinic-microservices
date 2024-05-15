package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/khensin166/PA2-Kel9/database"
	"github.com/khensin166/PA2-Kel9/model/entity"
	"github.com/khensin166/PA2-Kel9/utils"
	"log"
)

func GetAllMedicalHistoryByToken(ctx *fiber.Ctx) error {
	//VALIDATION
	// membuat token
	token := ctx.Get("Authorization")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated err",
		})
	}

	i, ok := claims["id"]
	log.Println(i)

	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse user ID",
		})
	}

	// membuat slice untuk entitas medical history
	var medicalHistory []entity.MedicalHistory

	// Memuat medical history yang sesuai dengan ID pengguna
	result := database.DB.Preload("User").Preload("DoctorReport").Preload("DoctorReport.NurseReport").Where("user_id = ?", i).Find(&medicalHistory)

	if result.Error != nil {
		log.Println(result.Error)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch medical history",
		})
	}

	return ctx.Status(200).JSON(medicalHistory)
}

func CreateMedicalHistory(ctx *fiber.Ctx) error {
	medicalHistory := new(entity.MedicalHistoryResponse)

	// Menangani error saat parsing request body
	if err := ctx.BodyParser(medicalHistory); err != nil {
		return err
	}

	//VALIDATION REQUEST
	validate := validator.New()
	err := validate.Struct(medicalHistory)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	// Mencoba membuat entitas baru dan menangani errornya
	if err := database.DB.Create(&medicalHistory).Error; err != nil {
		// Mengembalikan respon error 500 dengan pesan yang sesuai
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(), // Menambahkan pesan error ke respon JSON
		})
	}

	// Mengembalikan respon sukses dengan data baru yang telah dibuat
	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    medicalHistory,
	})
}
