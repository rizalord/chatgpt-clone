package config

import (
	"api-gateway/app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config * viper.Viper) *fiber.App {
	// Setup configuration for Fiber
	var app = fiber.New(fiber.Config{
		ErrorHandler: NewErrorHandler(),
	})

	// Setup routes
	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		if e, ok := err.(*model.RequestError); ok {
			return ctx.Status(e.Code).JSON(e)
		}

		return ctx.Status(err.(*fiber.Error).Code).JSON(fiber.Map{
			"message": err.(*fiber.Error).Message,
		})
	}
}