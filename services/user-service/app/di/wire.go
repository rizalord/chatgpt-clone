//go:build wireinject
// +build wireinject

package di

import (
	"user-service/app/config"
	"user-service/app/delivery/handler"
	"user-service/app/repository"
	"user-service/app/route"
	"user-service/app/usecase"

	"github.com/google/wire"
)

var userSet = wire.NewSet(
	repository.NewUserRepositoryImpl,
	wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)),
	usecase.NewUserUseCaseImpl,
	wire.Bind(new(usecase.UserUseCase), new(*usecase.UserUseCaseImpl)),
	handler.NewUserServiceImpl,
	wire.Bind(new(handler.UserService), new(*handler.UserServiceImpl)),
)

func InitializedApp() *config.App {
	wire.Build(
		config.NewViper,
		config.NewDatabase,
		config.NewApp,
		config.NewGrpcServer,
		config.NewValidator,
		config.NewLogger,
		route.NewGrpcServerRouter,
		userSet,
	)
	return nil
}