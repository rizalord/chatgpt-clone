//go:build wireinject
// +build wireinject

package di

import (
	"auth-service/app/config"
	"auth-service/app/delivery/handler"
	"auth-service/app/repository"
	"auth-service/app/route"
	"auth-service/app/usecase"

	"github.com/google/wire"
)

var chatSet = wire.NewSet(
	repository.NewChatRepositoryImpl,
	wire.Bind(new(repository.ChatRepository), new(*repository.ChatRepositoryImpl)),
	usecase.NewChatUseCaseImpl,
	wire.Bind(new(usecase.ChatUseCase), new(*usecase.ChatUseCaseImpl)),
)

var messageSet = wire.NewSet(
	repository.NewMessageRepositoryImpl,
	wire.Bind(new(repository.MessageRepository), new(*repository.MessageRepositoryImpl)),
	usecase.NewMessageUseCaseImpl,
	wire.Bind(new(usecase.MessageUseCase), new(*usecase.MessageUseCaseImpl)),
)

var serviceSet = wire.NewSet(
	handler.NewChatServiceImpl,
	wire.Bind(new(handler.ChatService), new(*handler.ChatServiceImpl)),
)

func InitializedApp() *config.App {
	wire.Build(
		config.NewViper,
		config.NewDatabase,
		config.NewApp,
		config.NewGrpcServer,
		config.NewValidator,
		config.NewLogger,
		config.NewGemini,
		route.NewGrpcServerRouter,
		chatSet,
		messageSet,
		serviceSet,
	)
	return nil
}