package handler

import (
	dorm "dorm/models/entity"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"user/database"
	"user/models/entity"
	"user/utils"
)

var PathImageProduct = "./Public"

func getDormByID(categoryID int) (*dorm.Dorm, error) {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8002/category/%d", categoryID))
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	var category dorm.Dorm
	if err := json.NewDecoder(resp.Body).Decode(&category); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return &category, nil
}

func init() {
	if _, err := os.Stat(PathImageProduct); os.IsNotExist(err) {
		err := os.Mkdir(PathImageProduct, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func UserHandlerGetAll(ctx *fiber.Ctx) error {

	userInfo := ctx.Locals("userInfo")
	log.Println("user info data :: ", userInfo)

	// membuat slice untuk entity users
	var users []entity.User

	// memanggil DB pada package database (cara 1)
	result := database.DB.Find(&users)

	if result.Error != nil {
		log.Println(result.Error)
	}

	return ctx.Status(200).JSON(users)

}

func CreateUser(ctx *fiber.Ctx) error {
	// Parse form data
	user := new(entity.User)

	// Menangani error saat parsing request body
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": "Invalid form data",
		})
	}

	//VALIDATION REQUEST
	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	// pemanggilan hashed password
	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
	// passing password yang sudah di hasing ke entity user (JSON)
	user.Password = hashedPassword

	image, err := ctx.FormFile("profilePicture")

	if err == nil {
		filename := utils.GenerateImageFile(user.Name, image.Filename)
		if err := ctx.SaveFile(image, filepath.Join(PathImageProduct, filename)); err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"status":  "failed",
				"message": "Can't save file image",
			})
		}
		user.ProfilePicture = &filename
	} else {
		user.ProfilePicture = nil
	}

	// Mencoba membuat entitas baru dan menangani errornya
	if err := database.DB.Create(&user).Error; err != nil {
		// Mengembalikan respon error 500 dengan pesan yang sesuai
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   err.Error(), // Menambahkan pesan error ke respon JSON
		})
	}

	// Mengembalikan respon sukses dengan data baru yang telah dibuat
	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})
}

func UserHandlerGetById(ctx *fiber.Ctx) error {

	// mencari user parameter id.
	userId := ctx.Params("id")

	// mendeklarasikan variabel user dengan tipe data userEntity
	var user entity.User

	// Query Statement dengan GORM
	err := database.DB.First(&user, "?", userId).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})
}

func UpdateUser(ctx *fiber.Ctx) error {
	// Parse form data
	userRequest := new(entity.User)
	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "bad request",
		})
	}

	// Get user ID from URL
	var user entity.User
	// logic
	userID := ctx.Params("id")

	// CHECK AVAILABLE USER
	err := database.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	// UPDATE USER DATA
	if userRequest.Name != "" {
		user.Name = userRequest.Name
	}
	user.Address = userRequest.Address

	user.Phone = userRequest.Phone

	// Process image if provided
	image, err := ctx.FormFile("profilePicture")
	if err == nil {
		filename := utils.GenerateImageFile(user.Name, image.Filename)
		if err := ctx.SaveFile(image, filepath.Join(PathImageProduct, filename)); err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"status":  "failed",
				"message": "Can't save file image",
			})
		}
		user.ProfilePicture = &filename
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	// Update the medicine in the database
	if err := database.DB.Save(&user).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "failed to update data",
			"error":   err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user entity.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
