package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"user-service/app/di"
)

func main() {
	app := di.InitializedApp()
	config := app.Config
	server := app.GrpcServerRouter.GrpcServer

	port := config.GetString("PORT")
	if port == "" {
		port = "3001"
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Port must be a number: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	defer server.GracefulStop()

	log.Printf("Server listening at %v", lis.Addr())
	if e := server.Serve(lis); e != nil {
		log.Fatalf("Failed to serve: %v", e)
	}
}
