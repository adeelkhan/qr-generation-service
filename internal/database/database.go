package database

import (
	"github.com/adeelkhan/qr-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect opens a database connection based on the provided configuration.
// It prefers PostgreSQL when a DATABASE_URL is set, and falls back to SQLite.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DatabaseURL != "" {
		return gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	}
	return gorm.Open(sqlite.Open("qr-service.db"), &gorm.Config{})
}
