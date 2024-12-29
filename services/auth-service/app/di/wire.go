//go:build wireinject
// +build wireinject

package di

import (
	"auth-service/app/config"
	"auth-service/app/delivery/client"
	"auth-service/app/delivery/handler"
	"auth-service/app/helper"
	"auth-service/app/route"
	"auth-service/app/usecase"

	"github.com/google/wire"
)

var authSet = wire.NewSet(
	usecase.NewAuthUseCaseImpl,
	wire.Bind(new(usecase.AuthUseCase), new(*usecase.AuthUseCaseImpl)),
	handler.NewAuthServiceImpl,
	wire.Bind(new(handler.AuthService), new(*handler.AuthServiceImpl)),
)

func InitializedApp() *config.App {
	wire.Build(
		config.NewViper,
		config.NewApp,
		config.NewGrpcServer,
		config.NewValidator,
		config.NewLogger,
		helper.NewJWTHelperImpl,
		wire.Bind(new(helper.JWTHelper), new(*helper.JWTHelperImpl)),
		helper.NewOauthImpl,
		wire.Bind(new(helper.OAuth), new(*helper.OAuthImpl)),
		route.NewGrpcServerRouter,
		authSet,
		client.NewUserClient,
	)
	return nil
}