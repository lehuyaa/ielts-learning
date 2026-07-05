package xp

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) FindUserForUpdate(tx *gorm.DB, userID uint) (models.User, error) {
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
		return models.User{}, fmt.Errorf("find user for xp update: %w", err)
	}

	return user, nil
}

func (r Repository) XPEventExists(tx *gorm.DB, userID uint, sourceType EventType, sourceID *uint) (bool, error) {
	var count int64
	query := tx.Model(&models.UserXPEvent{}).
		Where("user_id = ? AND source_type = ?", userID, string(sourceType))

	if sourceID == nil {
		query = query.Where("source_id IS NULL")
	} else {
		query = query.Where("source_id = ?", *sourceID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check xp event exists: %w", err)
	}

	return count > 0, nil
}

func (r Repository) CreateXPEvent(tx *gorm.DB, event models.UserXPEvent) error {
	if err := tx.Create(&event).Error; err != nil {
		return fmt.Errorf("create xp event: %w", err)
	}

	return nil
}

func (r Repository) SaveUser(tx *gorm.DB, user *models.User) error {
	if err := tx.Save(user).Error; err != nil {
		return fmt.Errorf("save xp user: %w", err)
	}

	return nil
}

func (r Repository) SumUserXP(tx *gorm.DB, userID uint) (int, error) {
	var result struct {
		Total int
	}

	if err := tx.Model(&models.UserXPEvent{}).
		Select("COALESCE(SUM(xp), 0) AS total").
		Where("user_id = ?", userID).
		Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("sum user xp: %w", err)
	}

	return result.Total, nil
}

func normalizeTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}

	return value
}
