package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"path/filepath"
	"staff/database"
	"staff/models/entity"
	"staff/utils"
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
func StaffHandlerGetAll(ctx *fiber.Ctx) error {

	staffInfo := ctx.Locals("staffInfo")
	log.Println("staff info data :: ", staffInfo)

	var staff []entity.Staff

	result := database.DB.Find(&staff)

	if result.Error != nil {
		log.Println(result.Error)
	}

	return ctx.Status(200).JSON(staff)

}

func CreateStaff(ctx *fiber.Ctx) error {
	staff := new(entity.Staff)

	// Menangani error saat parsing request body
	if err := ctx.BodyParser(staff); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": "Invalid form data",
		})
	}

	//VALIDATION REQUEST
	validate := validator.New()
	err := validate.Struct(staff)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"mesage": "failed",
			"error":  err.Error(),
		})
	}

	// pemanggilan hashed password
	hashedPassword, err := utils.HashingPassword(staff.Password)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
	// passing password yang sudah di hasing ke entity user (JSON)
	staff.Password = hashedPassword

	image, err := ctx.FormFile("profilePicture")

	if err == nil {
		filename := utils.GenerateImageFile(staff.Name, image.Filename)
		if err := ctx.SaveFile(image, filepath.Join(PathImageProduct, filename)); err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"status":  "failed",
				"message": "Can't save file image",
			})
		}
		staff.ProfilePicture = &filename
	} else {
		staff.ProfilePicture = nil
	}

	// Mencoba membuat entitas baru dan menangani errornya
	if err := database.DB.Create(&staff).Error; err != nil {
		// Mengembalikan respon error 500 dengan pesan yang sesuai
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(), // Menambahkan pesan error ke respon JSON
		})
	}

	// Mengembalikan respon sukses dengan data baru yang telah dibuat
	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    staff,
	})
}

func StaffHandlerGetById(ctx *fiber.Ctx) error {
	// mencari user parameter id.
	staffId := ctx.Params("id")

	// mendeklarasikan variabel user dengan tipe data userEntity
	var staff entity.Staff

	// Query Statement dengan GORM
	err := database.DB.First(&staff, "?", staffId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "staff not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    staff,
	})
}

func UpdateStaff(ctx *fiber.Ctx) error {

	staffRequest := new(entity.Staff)
	if err := ctx.BodyParser(staffRequest); err != nil {
		return ctx.Status(404).JSON(fiber.Map{"message": "bad request"})
	}

	var staff entity.Staff

	// logic
	userID := ctx.Params("id")
	// CHECK AVAILABLE USER
	err := database.DB.First(&staff, "id = ?", userID).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	// UPDATE USER DATA
	if staffRequest.Name != "" {
		staff.Name = staffRequest.Name
	}
	staff.Address = staffRequest.Address
	staff.Phone = staffRequest.Phone

	// Process image if provided
	image, err := ctx.FormFile("profilePicture")
	if err == nil {
		filename := utils.GenerateImageFile(staff.Name, image.Filename)
		if err := ctx.SaveFile(image, filepath.Join(PathImageProduct, filename)); err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"status":  "failed",
				"message": "Can't save file image",
			})
		}
		staff.ProfilePicture = &filename
	}

	validate := validator.New()
	if err := validate.Struct(staff); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	// Update the medicine in the database
	if err := database.DB.Save(&staff).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "failed to update data",
			"error":   err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   staff,
	})
}

func DeleteStaff(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var staff entity.Staff

	if err := database.DB.First(&staff, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff not found",
		})
	}

	if err := database.DB.Delete(&staff).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete staff",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Staff deleted successfully",
	})
}
