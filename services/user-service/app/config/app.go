package config

import (
	"user-service/app/route"

	"github.com/spf13/viper"
)

type App struct {
	GrpcServerRouter *route.GrpcServerRouter
	Config *viper.Viper
}

func NewApp(
	grpcServerRouter *route.GrpcServerRouter,
	config *viper.Viper,
) *App {
	return &App{
		GrpcServerRouter: grpcServerRouter,
		Config: config,
	}
}