package controllers

import (
	"time"

	"eventryx.api_service/config"
	"eventryx.api_service/internal/database/models"
	"eventryx.api_service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type userAuthInfo struct {
	Username string `json:"username" validate:"required,min=3,max=32" message:"Username must be a valid string with length from 3 to 32"`
	Password string `json:"password" validate:"required,min=12,max=64" message:"Password must be a valid string with length from 12 to 64"`
}

// Login
// @Summary      Login
// @Description  Method for login
// @Tags         Auth
// @Produce      json
// @Param        request                body                                                    userAuthInfo                             true    "request"
// @Success      200                    {object}                                                map[string]interface{}
// @Failure      500                    {object}                                                map[string]interface{}
// @Router       /api/v1/auth [post]
func Login(c *fiber.Ctx) error {
	request := new(userAuthInfo)

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

	user := &models.User{Name: &request.Username}
	if !user.Get() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Username or password is incorrect",
		})
	}

	if !utils.CheckPassword(request.Password, *user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Username or password is incorrect",
		})
	}

	accessClaims := jwt.MapClaims{
		"id":   user.Id,
		"name": user.Name,
		"role": user.Role,
		"exp": time.Now().
			Add(time.Hour * time.Duration(config.Config.AccessTokenLifespan)).
			Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.Id,
		"exp": time.Now().
			Add(time.Hour * time.Duration(config.Config.RefreshTokenLifespan)).
			Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	access, err := accessToken.SignedString([]byte(config.Config.TokenSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	refresh, err := refreshToken.SignedString([]byte(config.Config.TokenSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"access":  access,
		"refresh": refresh,
	})
}

// AddUser
// @Summary      Add user
// @Description  Method for adding user
// @Tags         Admin
// @Produce      json
// @Security     User
// @Param        request                body                                                    userAuthInfo                             true    "request"
// @Success      200                    {object}                                                map[string]interface{}
// @Failure      500                    {object}                                                map[string]interface{}
// @Router       /api/v1/users [post]
func AddUser(c *fiber.Ctx) error {
	request := new(userAuthInfo)

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

	user := models.User{Name: &request.Username}
	if user.Exists() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User with this name already exists"})
	}

	user.Password = &request.Password
	user.Role = models.IsUser

	if err := user.Create(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "Technical troubles, please try again later"})
	}

	return c.JSON(fiber.Map{"message": "User successfully added"})
}
