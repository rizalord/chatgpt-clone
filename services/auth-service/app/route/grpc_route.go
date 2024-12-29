package route

import (
	"auth-service/app/delivery/handler"
	"auth-service/app/model/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServerRouter struct {
	GrpcServer *grpc.Server
}

func NewGrpcServerRouter(
	grpcServer *grpc.Server,
	authService handler.AuthService,
) *GrpcServerRouter {
	proto.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	return &GrpcServerRouter{
		GrpcServer: grpcServer,
	}
}