package main

import (
	"flag"
	"fmt"
	"log"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/config"
	"ielts-learning/backend/internal/database"
	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/seeds"
)

func main() {
	promoteAdminEmail := flag.String("promote-admin", "", "Email of an existing user to promote to ADMIN role, instead of running the seed")
	onlyFoundation := flag.Bool("only-foundation", false, "Seed only the A1 Foundation band level (topics, lessons, quizzes from imported words), without demo data")
	flag.Parse()

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if *promoteAdminEmail != "" {
		if err := promoteAdmin(db, *promoteAdminEmail); err != nil {
			log.Fatalf("failed to promote user: %v", err)
		}
		log.Printf("promoted %s to ADMIN", *promoteAdminEmail)
		return
	}

	if *onlyFoundation {
		if err := seeds.RunFoundation(db); err != nil {
			log.Fatalf("failed to seed A1 foundation: %v", err)
		}
		log.Println("A1 foundation seed completed")
		return
	}

	if err := seeds.Run(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	log.Println("seed completed")
}

func promoteAdmin(db *gorm.DB, email string) error {
	result := db.Model(&models.User{}).
		Where("email = ?", email).
		Update("role", models.UserRoleAdmin)
	if result.Error != nil {
		return fmt.Errorf("update user role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with email %q", email)
	}

	return nil
}
