package main

import (
	"auth-service/app/di"
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {
	app := di.InitializedApp()
	config := app.Config
	server := app.GrpcServerRouter.GrpcServer
	db, _ := app.Gorm.DB()

	defer db.Close()
	defer server.GracefulStop()

	port := config.GetString("PORT")
	if port == "" {
		port = "3003"
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Port must be a number: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	defer lis.Close()

	// Log the server address
	log.Printf("Server listening at %v", lis.Addr())
	if e := server.Serve(lis); e != nil {
		log.Fatalf("Failed to serve: %v", e)
	}
}
