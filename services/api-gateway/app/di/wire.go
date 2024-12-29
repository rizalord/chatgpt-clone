//go:build wireinject
// +build wireinject

package di

import (
	"api-gateway/app/config"
	"api-gateway/app/delivery/client"
	"api-gateway/app/delivery/http"
	"api-gateway/app/delivery/http/middleware"
	"api-gateway/app/route"
	"api-gateway/app/usecase"

	"github.com/google/wire"
)

var authSet = wire.NewSet(
	usecase.NewAuthUseCaseImpl,
	wire.Bind(new(usecase.AuthUseCase), new(*usecase.AuthUseCaseImpl)),
	http.NewAuthController,
)

var chatSet = wire.NewSet(
	usecase.NewChatUseCaseImpl,
	wire.Bind(new(usecase.ChatUseCase), new(*usecase.ChatUseCaseImpl)),
	http.NewChatController,
)

var grpcClientSet = wire.NewSet(
	client.NewAuthClient,
	client.NewChatClient,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
	middleware.NewWebsocketMiddleware,
)

func InitializedApp() *config.App {
	wire.Build(
		config.NewViper,
		config.NewApp,
		config.NewFiber,
		config.NewValidator,
		config.NewLogger,
		route.NewHttpRouter,
		middlewareSet,
		authSet,
		chatSet,
		grpcClientSet,
	)
	return nil
}