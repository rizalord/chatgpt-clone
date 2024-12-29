package main

import (
	"api-gateway/app/di"
	"fmt"
	"log"
	"strconv"
)

func main() {
	app := di.InitializedApp()
	config := app.Config

	defer app.AuthClient.Conn.Close()
	defer app.ChatClient.Conn.Close()

	port := config.GetString("PORT")
	if port == "" {
		port = "3004"
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Port must be a number: %v", err)
	}

	err := app.Fiber.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
