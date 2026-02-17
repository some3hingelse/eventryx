package handlers

import (
	"eventryx.api_service/internal/controllers"
	"eventryx.api_service/internal/database/models"
	"eventryx.api_service/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("liveness", controllers.Liveness)
	api := app.Group("/api/v1")

	api.Post("auth", controllers.Login)
	api.Post("users", middleware.UserAuthMiddleware(models.IsAdmin), controllers.AddUser)

	servicesApi := api.Group("services")

	servicesApi.Post("", middleware.UserAuthMiddleware(models.IsUser), controllers.RegisterService)
	servicesApi.Post(":id/tokens", middleware.UserAuthMiddleware(models.IsUser), controllers.CreateServiceToken)
	servicesApi.Post(":id/data", middleware.ServiceAuthMiddleware, controllers.SendData)
}
