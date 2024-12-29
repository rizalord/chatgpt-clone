package client

import (
	"api-gateway/app/model/proto"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Service proto.AuthServiceClient
	Conn *grpc.ClientConn
}

func NewAuthClient(config *viper.Viper) *AuthClient {
	authHost := config.GetString("GRPC_AUTH_SERVICE")
	secure := config.GetBool("GRPC_AUTH_SERVICE_SECURE")

	var opts []grpc.DialOption

	if authHost != "" {
		opts = append(opts, grpc.WithAuthority(authHost))
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
	
	authConn, err := grpc.NewClient(authHost, opts...)
	if err != nil {
		panic(fmt.Errorf("error connecting to auth service: %v", err))
	}

	userService := proto.NewAuthServiceClient(authConn)

	return &AuthClient{
		Service: userService,
		Conn: authConn,
	}
}