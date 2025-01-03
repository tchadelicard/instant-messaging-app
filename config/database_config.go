package config

import (
	"fmt"
	"log"
	"os"

	"instant-messaging-app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	// Model migrations
	err = DB.AutoMigrate(&models.User{}, &models.Message{})
	if err != nil {
		log.Fatalf("Error during model migration: %v", err)
	}

	log.Println("Database connection and migration successful!")
}