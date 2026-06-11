package main

import (
	"log"
	"github.com/adeelkhan/qr-service/internal/config"
	"github.com/adeelkhan/qr-service/internal/database"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	_ = db
	log.Println("connected to database")
}
