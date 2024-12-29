package route

import (
	"user-service/app/delivery/handler"
	"user-service/app/model/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServerRouter struct {
	GrpcServer *grpc.Server
}

func NewGrpcServerRouter(
	grpcServer *grpc.Server,
	userService handler.UserService,
) *GrpcServerRouter {
	proto.RegisterUserServiceServer(grpcServer, userService)
	reflection.Register(grpcServer)

	return &GrpcServerRouter{
		GrpcServer: grpcServer,
	}
}