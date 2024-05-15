package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"nurseReport/database"
	"nurseReport/model/entity"
	"nurseReport/utils"
)

func NurseReportGetAll(ctx *fiber.Ctx) error {
	var nurseReport []entity.NurseReportResponse
	database.DB.Find(&nurseReport)
	return ctx.Status(200).JSON(fiber.Map{
		"nurse_reports": nurseReport,
	})
}

func CreateNurseReport(ctx *fiber.Ctx) error {
	nurseReport := &entity.NurseReport{}

	if err := ctx.BodyParser(nurseReport); err != nil {
		return err
	}

	// VALIDATION Request
	validate := validator.New()
	if err := validate.Struct(nurseReport); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	//VALIDATION
	// membuat token
	token := ctx.Get("Authorization")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	//_, err := utils.VerifyToken(token)
	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated err",
		})
	}

	//ctx.Locals("userInfo", claims)
	role := claims["role"].(float64)

	log.Println(role)

	if role == 2 {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Doctor can't create the data",
		})
	}

	if err := database.DB.Create(nurseReport).Error; err != nil {
		// Mengembalikan respon error 500 dengan pesan yang sesuai
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(), // Menambahkan pesan error ke respon JSON
		})
	}

	response := entity.NurseReportResponse{
		ID:                     nurseReport.ID,
		Temperature:            nurseReport.Temperature,
		Systole:                nurseReport.Systole,
		Diastole:               nurseReport.Diastole,
		Pulse:                  nurseReport.Pulse,
		OxygenSaturation:       nurseReport.OxygenSaturation,
		Respiration:            nurseReport.Respiration,
		Height:                 nurseReport.Height,
		Weight:                 nurseReport.Weight,
		AbdominalCircumference: nurseReport.AbdominalCircumference,
		Allergy:                nurseReport.Allergy,
		StaffNurseID:           nurseReport.StaffNurseID,
		PatientID:              nurseReport.PatientID,
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message":      "create data successfully",
		"nurse_report": response,
	})

}

func UpdateNurseReport(ctx *fiber.Ctx) error {
	reportRequest := new(entity.NurseReportResponseUpdate)

	if err := ctx.BodyParser(reportRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "Bad Request"})
	}

	var report entity.NurseReport

	reportID := ctx.Params("id")

	// CHECK AVAILABLE DOCTOR REPORT
	err := database.DB.First(&report, "id = ?", reportID).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "NurseReport not found",
		})
	}

	// UPDATE REPORT DATA
	if reportRequest.Temperature != "" {
		report.Temperature = reportRequest.Temperature
	}

	if reportRequest.Systole != "" {
		report.Systole = reportRequest.Systole
	}

	if reportRequest.Diastole != "" {
		report.Diastole = reportRequest.Diastole
	}

	if reportRequest.Pulse != "" {
		report.Pulse = reportRequest.Pulse
	}

	if reportRequest.OxygenSaturation != "" {
		report.OxygenSaturation = reportRequest.OxygenSaturation
	}

	if reportRequest.Respiration != "" {
		report.Respiration = reportRequest.Respiration
	}

	if reportRequest.Height != 0 {
		report.Height = reportRequest.Height
	}

	if reportRequest.Weight != 0 {
		report.Weight = reportRequest.Weight
	}

	if reportRequest.AbdominalCircumference != 0 {
		report.AbdominalCircumference = reportRequest.AbdominalCircumference
	}

	if reportRequest.Allergy != "" {
		report.Allergy = reportRequest.Allergy
	}

	errUpdate := database.DB.Save(&report).Error
	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success adding data",
		"data":    report,
	})
}

func DeleteNurseReport(ctx *fiber.Ctx) error {
	reportID := ctx.Params("id")

	var report entity.NurseReport
	if err := database.DB.First(&report, reportID).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error": "Nurse Report not found",
		})
	}

	if err := database.DB.Delete(&report).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to delete Nurse Report",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Nurse Report deleted successfully",
	})
}
