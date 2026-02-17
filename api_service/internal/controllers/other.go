package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Liveness
// @Summary      Liveness probe
// @Description  Liveness probe
// @Tags         HealthChecks
// @Produce      json
// @Success      200                    {object}                                                map[string]interface{}
// @Failure      500                    {object}                                                map[string]interface{}
// @Router       /liveness [get]
func Liveness(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
