package database

import (
	"countries-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type DB struct {
	Db *gorm.DB
}

var Database DB

func ConnectDB(dbName string) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Connected to database")

	db.Logger = logger.Default.LogMode(logger.Info)

	// TODO: add migrations
	err = db.AutoMigrate(&models.Country{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	Database.Db = db
}
