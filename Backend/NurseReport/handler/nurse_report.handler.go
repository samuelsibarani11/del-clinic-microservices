package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"nurseReport/database"
	"nurseReport/model/entity"
	"nurseReport/utils"
	usr "service/user/Models/entity"
	staff "staff/models/entity"
)

func NurseReportGetById(ctx *fiber.Ctx) error {
	// mencari user parameter id.
	nurseReportId := ctx.Params("id")

	// mendeklarasikan variabel user dengan tipe data userEntity
	var nurseReport entity.NurseReport

	// Query Statement dengan GORM
	err := database.DB.First(&nurseReport, "?", nurseReportId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "staff not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    nurseReport,
	})
}

func getStaffByID(staffID int) (*staff.Staff, error) {
	resp, err := http.Get(fmt.Sprintf("http://172.20.10.4:8004/staff/%d", staffID))
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var staffs staff.Staff
	if err := json.NewDecoder(resp.Body).Decode(&staffs); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return &staffs, nil
}

func getUserByID(staffID int) (*usr.User, error) {
	resp, err := http.Get(fmt.Sprintf("http://172.20.10.4:8004/staff/%d", staffID))
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var users usr.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return &users, nil
}

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
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed to parse request body",
			"error":   err.Error(),
		})
	}

	// VALIDATION Request
	validate := validator.New()
	if err := validate.Struct(nurseReport); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	// Validate Authorization Token
	token := ctx.Get("Authorization")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
			"error":   err.Error(),
		})
	}

	role := claims["role"].(float64)
	log.Println("Role:", role)

	if role == 2 {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Doctor can't create the data",
		})
	}

	if err := database.DB.Create(nurseReport).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(),
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
		"message":      "created data successfully",
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
