// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/google/wire"
	"user-service/app/config"
	"user-service/app/delivery/handler"
	"user-service/app/repository"
	"user-service/app/route"
	"user-service/app/usecase"
)

// Injectors from wire.go:

func InitializedApp() *config.App {
	logger := config.NewLogger()
	server := config.NewGrpcServer(logger)
	viper := config.NewViper(logger)
	db := config.NewDatabase(viper)
	validate := config.NewValidator()
	userRepositoryImpl := repository.NewUserRepositoryImpl()
	userUseCaseImpl := usecase.NewUserUseCaseImpl(db, validate, userRepositoryImpl, logger)
	userServiceImpl := handler.NewUserServiceImpl(userUseCaseImpl)
	grpcServerRouter := route.NewGrpcServerRouter(server, userServiceImpl)
	app := config.NewApp(grpcServerRouter, viper)
	return app
}

// wire.go:

var userSet = wire.NewSet(repository.NewUserRepositoryImpl, wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)), usecase.NewUserUseCaseImpl, wire.Bind(new(usecase.UserUseCase), new(*usecase.UserUseCaseImpl)), handler.NewUserServiceImpl, wire.Bind(new(handler.UserService), new(*handler.UserServiceImpl)))