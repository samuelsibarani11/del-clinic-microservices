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
	"service/user/database"
	"service/user/models/entity"
	"service/user/utils"
	"strconv"
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

func getDormByID(categoryID int) (*dorm.Dorm, error) {
	resp, err := http.Get(fmt.Sprintf("http://172.20.10.4:8001/dorm/%d", categoryID))
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

	// Handle error when parsing request body
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": "Invalid form data",
		})
	}

	// Validation request
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	// Parse dorm ID from the user data
	dormId := user.DormID

	// Hash the password
	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": "Internal server error",
		})
	}
	// Assign hashed password to user
	user.Password = hashedPassword

	// Handle profile picture upload
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

	// Retrieve dorm by ID
	dormm, err := getDormByID(int(dormId))
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}
	dormm.ID = uint(dormId)

	userData := entity.User{
		Name:           user.Name,
		Address:        user.Address,
		Weight:         user.Weight,
		Height:         user.Height,
		Role:           user.Role,
		Gender:         user.Gender,
		Phone:          user.Phone,
		ProfilePicture: user.ProfilePicture,
		Birthday:       user.Birthday,
		NIK:            user.NIK,
		Age:            user.Age,
		Username:       user.Username,
		Password:       user.Password,
		DormID:         dormm.ID,
	}

	// Create new entity and handle error
	if err := database.DB.Create(&userData).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "Failed to store data",
			"error":   err.Error(),
		})
	}

	// Return success response with new data
	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   userData,
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

	dormId, err := strconv.Atoi(ctx.FormValue("dormID"))
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	if dormId != 0 {
		dorm, err := getDormByID(dormId)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{
				"status":  "failed",
				"message": err.Error(),
			})
		}
		dorm.ID = uint(dormId)
		user.DormID = uint(dormId)
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
