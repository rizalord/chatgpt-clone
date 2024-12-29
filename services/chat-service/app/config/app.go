package config

import (
	"auth-service/app/route"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type App struct {
	GrpcServerRouter *route.GrpcServerRouter
	Config *viper.Viper
	Gorm *gorm.DB
}

func NewApp(
	grpcServerRouter *route.GrpcServerRouter,
	config *viper.Viper,
	gormDB *gorm.DB,
) *App {
	return &App{
		GrpcServerRouter: grpcServerRouter,
		Config: config,
		Gorm: gormDB,
	}
}