package client

import (
	"auth-service/app/model/proto"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	Service proto.UserServiceClient
	Conn *grpc.ClientConn
}

func NewUserClient(config *viper.Viper) *UserClient {
	userHost := config.GetString("GRPC_USER_SERVICE")
	secure := config.GetBool("GRPC_USER_SERVICE_SECURE")

	var opts []grpc.DialOption

	if userHost != "" {
		opts = append(opts, grpc.WithAuthority(userHost))
	}

	if !secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			panic(fmt.Errorf("error loading system root CA pool: %v", err))
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	userConn, err := grpc.NewClient(userHost, opts...)
	if err != nil {
		panic(fmt.Errorf("error connecting to user service: %v", err))
	}

	userService := proto.NewUserServiceClient(userConn)

	return &UserClient{
		Service: userService,
		Conn: userConn,
	}
}