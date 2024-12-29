package config

import (
	"auth-service/app/delivery/client"
	"auth-service/app/route"

	"github.com/spf13/viper"
)

type App struct {
	Config *viper.Viper
	GrpcServerRouter *route.GrpcServerRouter
	UserClient *client.UserClient
}

func NewApp(
	config *viper.Viper,
	grpcServerRouter *route.GrpcServerRouter,
	userClient *client.UserClient,
) *App {
	return &App{
		Config: config,
		GrpcServerRouter: grpcServerRouter,
		UserClient: userClient,
	}
}