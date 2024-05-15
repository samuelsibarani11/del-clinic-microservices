package middleware

import (
	"appointment/utils"
	"github.com/gofiber/fiber/v2"
)

var ActiveTokens = make(map[string]bool)

func Auth(ctx *fiber.Ctx) error {
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

	ctx.Locals("staffInfo", claims)

	return ctx.Next()

}
