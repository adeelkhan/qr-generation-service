package database

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// sqliteDriver returns a fresh isolated in-memory SQLite database.
// Each call gets a unique name so parallel tests don't share state.
func sqliteDriver() gorm.Dialector {
	return sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.New().String()))
}
