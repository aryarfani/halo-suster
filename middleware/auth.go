package middleware

import (
	"eniqilo-store/config"
	"eniqilo-store/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// Auth protect routes
func Auth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(config.GetConfig("JWT_SECRET"))},
		SuccessHandler: jwtSuccess,
		ErrorHandler:   jwtError,
	})
}

func jwtSuccess(c *fiber.Ctx) error {
	claims, err := utils.ExtractToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Locals("user_id", claims.UserID)
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing or malformed JWT",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid or expired JWT",
		"data":    nil,
	})
}
