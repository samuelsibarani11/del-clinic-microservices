package middleware

import (
	"github.com/gofiber/fiber/v2"
	"user/utils"
)

var ActiveTokens = make(map[string]bool)

func Auth(ctx *fiber.Ctx) error {

	// Menerima token dari header Authorization
	token := ctx.Get("Authorization")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Mendekode token untuk mendapatkan klaim
	claims, err := utils.DecodeToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated err",
		})
	}

	// Memeriksa apakah token ada dalam daftar token aktif
	if _, ok := ActiveTokens[token]; !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Menyimpan informasi pengguna ke dalam konteks lokal
	ctx.Locals("staffInfo", claims)

	return ctx.Next()

}
