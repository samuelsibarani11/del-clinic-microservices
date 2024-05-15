package handler

import (
	"doctorReport/database"
	"doctorReport/model/entity"
	"doctorReport/utils"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	md "medicine/model/entity"
	"net/http"
	nr "nurseReport/model/entity"
)

func GetNurseReportById(nurseReportId int) (*nr.NurseReport, error) {
	resp, err := http.Get(fmt.Sprintf("http://172.20.10.4:8004/nurse-report/%d", nurseReportId))
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var nurseReport nr.NurseReport
	if err := json.NewDecoder(resp.Body).Decode(&nurseReport); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return &nurseReport, nil
}

func GetMedicineById(medicineId int) (*md.Medicine, error) {
	resp, err := http.Get(fmt.Sprintf("http://172.20.10.4:8004/medicine/%d", medicineId))
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var medicine md.Medicine
	if err := json.NewDecoder(resp.Body).Decode(&medicine); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return &medicine, nil
}

func DoctorReportGetAll(ctx *fiber.Ctx) error {
	var doctorReports []entity.DoctorReport

	if err := database.DB.Preload("NurseReport").
		Preload("NurseReport.Patient").
		Preload("NurseReport.StaffNurse").
		Preload("StaffDoctor").
		Preload("Medicines"). // Preload medicines association
		Find(&doctorReports).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to fetch doctor reports",
			"error":   err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"doctor_reports": doctorReports,
	})
}

func CreateDoctorReport(ctx *fiber.Ctx) error {
	var doctorReport entity.DoctorReport

	if err := ctx.BodyParser(&doctorReport); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(doctorReport); err != nil {
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

	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated err",
		})
	}

	role := claims["role"].(float64)
	if role == 1 {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Doctor can't create the data",
		})
	}

	reportId := doctorReport.NurseReportID
	medicineId := doctorReport.MedicineID

	nurseReports, err2 := GetNurseReportById(int(reportId))

	if err2 != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	medicines, err3 := GetMedicineById(int(medicineId))
	if err3 != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}
	medicines.ID = uint(medicineId)
	nurseReports.ID = reportId

	userData := entity.DoctorReport{
		NurseReportID: nurseReports.ID,
		Disease:       doctorReport.Disease,
		Amount:        doctorReport.Amount,
		StaffDoctorID: doctorReport.StaffDoctorID,
		MedicineID:    int(medicines.ID),
	}

	// Create the DoctorReport
	if err := database.DB.Create(&userData).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(),
		})
	}

	// Load relations
	if err := database.DB.Preload("NurseReport").Preload("StaffDoctor").Preload("Medicines").First(&doctorReport, doctorReport.ID).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to load relations",
			"error":   err.Error(),
		})
	}

	return ctx.Status(200).JSON(doctorReport)
}

func UpdateDoctorReport(ctx *fiber.Ctx) error {
	reportRequest := new(entity.DoctorReportResponse)
	if err := ctx.BodyParser(reportRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "Bad request"})
	}

	var report entity.DoctorReport

	reportID := ctx.Params("id")
	// CHECK AVAILABLE REPORT
	err := database.DB.First(&report, "id = ?", reportID).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "DoctorReport not found",
		})
	}

	// UPDATE REPORT DATA
	if reportRequest.Disease != "" {
		report.Disease = reportRequest.Disease
	}
	if reportRequest.NurseReportID != 0 {
		report.NurseReportID = reportRequest.NurseReportID
	}

	if reportRequest.StaffDoctorID != 0 {
		report.StaffDoctorID = reportRequest.StaffDoctorID
	}

	if reportRequest.Disease != "" {
		report.Disease = reportRequest.Disease
	}

	errUpdate := database.DB.Save(&report).Error

	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    report,
	})
}

func DeleteDoctorReport(ctx *fiber.Ctx) error {
	reportID := ctx.Params("id")

	var report entity.DoctorReport
	if err := database.DB.First(&report, reportID).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"error": "DoctorReport not found",
		})
	}

	if err := database.DB.Delete(&report).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to delete DoctorReport",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "DoctorReport deleted successfully",
	})
}
