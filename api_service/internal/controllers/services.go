package controllers

import (
	"strconv"
	"time"

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

// CreateServiceToken
// @Summary      Create service token
// @Description  Method for service token creating
// @Tags         Services
// @Security	 User
// @Produce      json
// @Param		 expires_at				formData												time.Time					    false	"expires at"
// @Param        id		                path                                                    int                             true    "service id"
// @Success      200                    {object}                                                map[string]interface{}
// @Failure      500                    {object}                                                map[string]interface{}
// @Router       /api/v1/services/{id}/tokens [post]
func CreateServiceToken(c *fiber.Ctx) error {
	expiresAt, _ := time.Parse("2006-01-02T15:04:05Z", c.FormValue("expires_at"))

	serviceId, _ := strconv.Atoi(c.Params("id"))
	service := models.Service{Id: &serviceId}
	if !service.Exists() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Service does not exist"})
	}

	serviceToken := models.ServiceToken{ServiceId: &serviceId}
	if !expiresAt.IsZero() {
		serviceToken.ExpiresAt = &expiresAt
	}
	
	err := serviceToken.Create()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "Technical troubles, please try again later"})
	}

	return c.JSON(serviceToken)
}
