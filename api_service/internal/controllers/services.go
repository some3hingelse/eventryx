package controllers

import (
	"eventryx.api_service/internal/database/models"
	"eventryx.api_service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type registerServiceRequest struct {
	Name string `json:"name" validate:"required"`
}

// RegisterService
// @Summary      Register service
// @Description  Method for service registration
// @Tags         Services
// @Security	 User
// @Produce      json
// @Param        request                body                                                    registerServiceRequest                             true    "request"
// @Success      200                    {object}                                                map[string]interface{}
// @Failure      500                    {object}                                                map[string]interface{}
// @Router       /api/v1/services [post]
func RegisterService(c *fiber.Ctx) error {
	request := new(registerServiceRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.Validator.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": utils.ValidationErrorsToMap(err, request),
		})
	}

	service := models.Service{Name: &request.Name}
	if service.Exists() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Service with this name already exists"})
	}

	service.OwnerId = c.Locals("user").(models.User).Id

	if err := service.Create(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "Technical troubles, please try again later"})
	}

	return c.JSON(service)
}
