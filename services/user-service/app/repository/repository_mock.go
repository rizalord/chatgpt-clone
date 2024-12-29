package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open sqlmock database: %v", err)
    }

    dialector := postgres.New(postgres.Config{
        Conn:       mockDB,
        DriverName: "postgres",
    })

    db, err := gorm.Open(dialector, &gorm.Config{
        SkipDefaultTransaction: true,
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        t.Fatalf("failed to open gorm database: %v", err)
    }

    return db, mock
}