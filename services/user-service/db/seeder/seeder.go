package main

import (
	"fmt"
	"user-service/app/config"
	"user-service/app/helper"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenConnection() *gorm.DB {
	// Read environment variables
	logger := logrus.New()	
	config := config.NewViper(logger)
	dbName := config.GetString("DATABASE_NAME")
	dbUser := config.GetString("DATABASE_USER")
	dbPass := config.GetString("DATABASE_PASS")
	dbHost := config.GetString("DATABASE_HOST")
	dbPort := config.GetString("DATABASE_PORT")
	dbSSLMode := config.GetString("DATABASE_SSLMODE")

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error connecting to database: %v", err))
	}

	return db
}

func SeedData(db *gorm.DB) {
	password, err := helper.HashPassword("secretpassword")
	if err != nil {
		panic(fmt.Errorf("error hashing password: %v", err))
	}

	// Seed users
	err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Example User", "user@example.com", password).Error
	if err != nil {
		panic(fmt.Errorf("error seeding users: %v", err))
	}
}

func main() {
	db := OpenConnection()
	fmt.Println("Database connected")
	SeedData(db)
	fmt.Println("Data seeded")
}