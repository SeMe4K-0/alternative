package main

import (
	"backend/app"
	"backend/config"
	"log"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// @title Comet Tracker API
// @version 1.0
// @description This is the API for the Don't Look Up Hackathon project.
// @description It allows users to track comets by submitting observations and calculating their orbits.

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name session_id
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: could not load .env file")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg := config.Load()

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	if err := application.Run(cfg.Server.Port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
