package config

import (
	"api-gateway/app/delivery/client"
	"api-gateway/app/route"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type App struct {
	Config *viper.Viper
	Fiber *fiber.App
	AuthClient *client.AuthClient
	ChatClient *client.ChatClient
}

func NewApp(
	config *viper.Viper,
	router *route.HttpRouter,
	authClient *client.AuthClient,
	chatClient *client.ChatClient,
) *App {
	return &App{
		Config: config,
		Fiber: router.Router,
		AuthClient: authClient,
		ChatClient: chatClient,
	}
}