package main

import (
	"backend/config"
	"backend/models"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables from OS")
	}

	cfg := config.Load()
	dbConfig := cfg.Database

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	log.Println("Connecting to the database...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection successful.")
	log.Println("Running migrations...")

	err = db.AutoMigrate(
		&models.User{},
		&models.PasswordResetToken{},
		&models.Comet{},
		&models.Observation{},
		&models.OrbitalCalculation{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	log.Println("✅ Database migration completed successfully!")
}