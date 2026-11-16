package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AanFLenrir/PetrivaBackend/config"
	"github.com/AanFLenrir/PetrivaBackend/database"
	"github.com/AanFLenrir/PetrivaBackend/route"
	"github.com/gofiber/fiber/v2"

	_ "github.com/AanFLenrir/PetrivaBackend/docs"
)

// @title PETRIVA API
// @version 1.0.0
// @description Pet Care Management System REST API built with Go, Fiber v2, PostgreSQL, and Clean Architecture.
// @termsOfService https://petriva.com/terms/

// @contact.name PETRIVA Engineering Team
// @contact.email team@petriva.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format "Bearer {token}"
func main() {
	cfg := config.Load()

	// Initialize DB pool
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "PETRIVA - Pet Care Management System v1.0.0",
	})

	// Register routes (will include middleware inside)
	route.Setup(app, cfg, db)

	// Graceful shutdown handling
	go func() {
		log.Printf("PETRIVA server listening on port :%s", cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down PETRIVA server gracefully...")
	_ = app.Shutdown()
}
