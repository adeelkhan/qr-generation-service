package database

import (
	"github.com/adeelkhan/qr-service/internal/config"
	"github.com/adeelkhan/qr-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.User{}, &models.QRCode{}); err != nil {
		return nil, err
	}
	return db, nil
}

// ConnectSQLite opens an in-memory SQLite database for tests.
func ConnectSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqliteDriver(), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.User{}, &models.QRCode{}); err != nil {
		return nil, err
	}
	return db, nil
}
