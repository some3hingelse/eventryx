package middleware

import (
	"eventryx.api_service/internal/database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UserAuthMiddleware(c *fiber.Ctx) error {
	userJwt := c.Locals("user").(*jwt.Token)
	claims := userJwt.Claims.(jwt.MapClaims)
	userId := int(claims["id"].(float64))

	user := models.User{Id: &userId}
	if !user.Get() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	c.Locals("user", user)
	return c.Next()
}

func AdminAuthMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	if user.Role != models.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden"})
	}

	return c.Next()
}
