package handlers

import (
	"eventryx.api_service/config"
	"eventryx.api_service/internal/controllers"
	"eventryx.api_service/internal/middleware"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("liveness", controllers.Liveness)

	api := app.Group("/api/v1")

	api.Post("auth", controllers.Login)

	api.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.Config.TokenSecret),
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		},
	}))
	api.Use(middleware.UserAuthMiddleware)
	api.Post("services", controllers.RegisterService)
	api.Post("services/:id/data", controllers.SendData)
	api.Use(middleware.AdminAuthMiddleware)
	api.Post("users", controllers.AddUser)
}
