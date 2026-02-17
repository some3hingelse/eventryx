package controllers

import (
	"context"
	"encoding/json"
	"time"

	"eventryx.api_service/internal/database/models"
	"eventryx.api_service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// SendData
// @Summary      Send data to event ingest service
// @Description  Method for send data to event ingest service
// @Tags         Services
// @Security	 Service
// @Produce      json
// @Param        request                body                                                    map[string]interface{}          true    "request"
// @Success      200                    {object}                                                map[string]interface{}
// @Failure      500                    {object}                                                map[string]interface{}
// @Router       /api/v1/services/data [post]
func SendData(c *fiber.Ctx) error {
	service := c.Locals("service").(models.Service)

	var requestData map[string]interface{}
	err := json.Unmarshal(c.Body(), &requestData)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "Technical troubles, please try again later"})
	}
	dataToSend := map[string]interface{}{
		"service_id": service.Id,
		"data":       requestData,
		"timestamp":  time.Now().UTC().Unix(),
	}
	jsonData, err := json.Marshal(dataToSend)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "Technical troubles, please try again later"})
	}

	err = utils.KafkaProducer.SendMessageToKafka(context.Background(), jsonData)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "Technical troubles, please try again later"})
	}

	return c.JSON(fiber.Map{"message": "OK!"})
}
