package main

import (
	"eventryx.api_service/config"
	_ "eventryx.api_service/docs"
	"eventryx.api_service/internal/database"
	"eventryx.api_service/internal/handlers"
	"eventryx.api_service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func init() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	database.InitConnectionString(
		config.Config.DbHost, config.Config.DbUsername, config.Config.DbPassword,
		config.Config.DbPort, config.Config.DbName,
	)
	database.CreateConnection()
	utils.CreateKafkaProducer(config.Config.KafkaBootstrapServers, config.Config.KafkaTopic)
}

// @title Eventryx API
// @version 1.0
// @description API for Eventryx.

// @securityDefinitions.apikey User
// @in header
// @name Authorization
// @description description

// @securityDefinitions.apikey Admin
// @in header
// @name Authorization
// @description description
func main() {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	app.Use(logger.New(logger.Config{Format: "${time} | ${status} | ${method} | ${path} | ${latency}\n"}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	handlers.RegisterRoutes(app)
	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}
}
