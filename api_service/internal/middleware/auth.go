package middleware

import (
	"fmt"
	"strings"

	"eventryx.api_service/config"
	"eventryx.api_service/internal/database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UserAuthMiddleware(role models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.Config.TokenSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := int(claims["id"].(float64))

		user := models.User{Id: &userId}
		if !user.Get() {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		}

		if user.Role != role {
			return c.Status(403).JSON(fiber.Map{"message": "Forbidden"})
		}

		c.Locals("user", user)
		return c.Next()
	}
}

func ServiceAuthMiddleware(c *fiber.Ctx) error {
	authHeader := strings.Split(c.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	token := models.ServiceToken{Value: &authHeader[1]}
	if !token.Get() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	c.Locals("service", *token.Service)
	return c.Next()
}
