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

type ChatClient struct {
	Service proto.ChatServiceClient
	Conn *grpc.ClientConn
}

func NewChatClient(config *viper.Viper) *ChatClient {
	chatHost := config.GetString("GRPC_CHAT_SERVICE")
	secure := config.GetBool("GRPC_CHAT_SERVICE_SECURE")

	var opts []grpc.DialOption

	if chatHost != "" {
		opts = append(opts, grpc.WithAuthority(chatHost))
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

	chatConn, err := grpc.NewClient(chatHost, opts...)
	if err != nil {
		panic(fmt.Errorf("error connecting to chat service: %v", err))
	}

	chatService := proto.NewChatServiceClient(chatConn)

	return &ChatClient{
		Service: chatService,
		Conn: chatConn,
	}
}