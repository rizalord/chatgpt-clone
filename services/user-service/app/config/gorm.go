package config

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(config *viper.Viper) *gorm.DB {
	dbName := config.GetString("DATABASE_NAME")
	dbUser := config.GetString("DATABASE_USER")
	dbPass := config.GetString("DATABASE_PASS")
	dbHost := config.GetString("DATABASE_HOST")
	dbPort := config.GetString("DATABASE_PORT")
	dbSSLMode := config.GetString("DATABASE_SSLMODE")
	dbLogging := config.GetBool("DATABASE_LOGGING")

	dbLog := logger.Default.LogMode(logger.Info)
	if !dbLogging {
		dbLog = logger.Default
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: 	true,
		Logger: 			   		dbLog,
	})
	if err != nil {
		panic(fmt.Errorf("error connecting to database: %v", err))
	}

	return db
}