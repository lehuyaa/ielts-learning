package main

import (
	"log"

	"ielts-learning/backend/internal/config"
	"ielts-learning/backend/internal/database"
	"ielts-learning/backend/seeds"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := seeds.Run(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	log.Println("seed completed")
}
