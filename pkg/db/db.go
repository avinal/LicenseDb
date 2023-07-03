package db

import (
	"fmt"
	"log"

	"github.com/fossology/LicenseDb/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dburi := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", "localhost", "5432", "fossy", "fossology", "fossy")
	gormConfig := &gorm.Config{}
	database, err := gorm.Open(postgres.Open(dburi), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(&models.License{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	if err := database.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	DB = database
}
