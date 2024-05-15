package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"staff/database"
	"staff/middleware"
	"staff/models/entity"
	"staff/models/request"
	"staff/utils"
	"time"
)

func isActiveToken(token string) bool {
	// Periksa apakah token ada dalam daftar ActiveTokens
	_, isActive := middleware.ActiveTokens[token]
	return isActive
}

func StaffLoginHandler(ctx *fiber.Ctx) error {
	// Ambil token dari header Authorization
	token := ctx.Get("Authorization")

	// Periksa apakah token masih aktif
	if isActiveToken(token) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "token masih aktif",
		})
	}

	loginRequest := new(request.LoginRequest)

	// Menangani error saat parsing request body
	if err := ctx.BodyParser(loginRequest); err != nil {
		return err
	}

	//VALIDATION REQUEST
	validate := validator.New()
	errValidate := validate.Struct(loginRequest)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"mesage": "failed",
			"error":  errValidate.Error(),
		})
	}

	// CHECK AVAILABLE USER
	var staff entity.Staff

	err := database.DB.First(&staff, "username = ?", loginRequest.Username).Error
	log.Println(err)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "username not found",
		})
	}

	// CHECK VALIDATION PASSWORD
	isValid := utils.CheckPasswordHash(loginRequest.Password, staff.Password)
	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "wrong password",
		})
	}

	//GENERATE JWT (don't forget to install package : github.com/dgrijalva/jwt-go)
	claims := jwt.MapClaims{}
	claims["id"] = staff.ID
	claims["name"] = staff.Name
	claims["age"] = staff.Age
	claims["weight"] = staff.Weight
	claims["height"] = staff.Height
	claims["nip"] = staff.NIP
	claims["birthday"] = staff.Birthday
	claims["gender"] = staff.Gender
	claims["address"] = staff.Address
	claims["phone"] = staff.Phone
	claims["username"] = staff.Username
	claims["role"] = staff.Role
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	token, errGenerateToken := utils.GenerateToken(&claims)

	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "wrong credentials",
		})

	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success login",
		"data":    claims,
		"token":   token,
	})

}

func StaffLogoutHandler(ctx *fiber.Ctx) error {
	// Ambil token dari header Authorization
	authHeader := ctx.Get("Authorization")

	// Verifikasi token
	_, err := utils.VerifyToken(authHeader)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	// Hapus token dari daftar token aktif
	delete(middleware.ActiveTokens, authHeader)

	return ctx.Status(200).JSON(fiber.Map{
		"message": "logout successful",
	})
}

func GetProfile(c *fiber.Ctx) error {
	admin := c.Locals("admin").(entity.Staff)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": admin,
	})
}
