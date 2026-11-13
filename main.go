package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AanFLenrir/campuscare-backend/config"
	"github.com/AanFLenrir/campuscare-backend/database"
	"github.com/AanFLenrir/campuscare-backend/route"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	// Initialize DB pool
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Initialize Fiber app
	app := fiber.New()

	// Register routes (will include middleware inside)
	route.Setup(app, cfg, db)

	// Graceful shutdown handling
	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	_ = app.Shutdown()
}
