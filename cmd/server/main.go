package main

import (
	"log"

	"github.com/adeelkhan/qr-service/internal/auth"
	"github.com/adeelkhan/qr-service/internal/config"
	"github.com/adeelkhan/qr-service/internal/database"
	"github.com/adeelkhan/qr-service/internal/middleware"
	"github.com/adeelkhan/qr-service/internal/qr"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	authSvc := auth.NewService(db, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)

	qrSvc := qr.NewService(db)
	qrHandler := qr.NewHandler(qrSvc)

	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "web/static")

	// Template routes
	r.GET("/", func(c *gin.Context) { c.Redirect(302, "/login") })
	r.GET("/login", func(c *gin.Context) { c.HTML(200, "login.html", nil) })
	r.GET("/register", func(c *gin.Context) { c.HTML(200, "register.html", nil) })
	r.GET("/dashboard", func(c *gin.Context) { c.HTML(200, "dashboard.html", nil) })

	// API routes
	api := r.Group("/api/v1")
	authHandler.RegisterRoutes(api)
	qrHandler.RegisterRoutes(api, middleware.JWT(cfg.JWTSecret))

	log.Printf("server running on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
