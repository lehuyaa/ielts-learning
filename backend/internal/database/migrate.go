package database

import (
	"fmt"

	"ielts-learning/backend/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(models.AutoMigrateModels()...); err != nil {
		return fmt.Errorf("auto migrate models: %w", err)
	}

	return nil
}
