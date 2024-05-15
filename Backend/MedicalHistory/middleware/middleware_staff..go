package middleware

import (
	"github.com/gofiber/fiber/v2"
	"medicalHistory/utils"

	"log"
)

func StaffAuth(ctx *fiber.Ctx) error {
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

	role := int(claims["role"].(float64))
	log.Println(role)

	if role != 1 && role != 2 {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden access",
		})
	}

	ctx.Locals("staffInfo", claims)

	return ctx.Next()

}
