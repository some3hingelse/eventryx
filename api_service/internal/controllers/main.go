package controllers

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"eventryx.api_service/internal/database/models"
	"eventryx.api_service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// SendData
// @Summary      Send data to event ingest service
// @Description  Method for send data to event ingest service
// @Tags         Services
// @Security	 User
// @Produce      json
// @Param        request                body                                                    map[string]interface{}                             true    "request"
// @Param        id		                path                                                    int                             true    "request"
// @Success      200                    {object}                                                map[string]interface{}
// @Failure      500                    {object}                                                map[string]interface{}
// @Router       /api/v1/services/{id}/data [post]
func SendData(c *fiber.Ctx) error {
	serviceId, _ := strconv.Atoi(c.Params("id"))
	service := models.Service{Id: &serviceId}
	if !service.Exists() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Service does not exist"})
	}

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
