package notification

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Create(tx *gorm.DB, notification models.Notification) error {
	if err := tx.Create(&notification).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	return nil
}

func (r Repository) FindByUser(userID uint, page int, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	baseQuery := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}

	offset := (page - 1) * limit
	if err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC, id DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error; err != nil {
		return nil, 0, fmt.Errorf("find notifications: %w", err)
	}

	return notifications, total, nil
}

func (r Repository) CountUnread(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}

	return count, nil
}

func (r Repository) FindByIDForUser(userID uint, notificationID uint) (models.Notification, error) {
	var notification models.Notification
	if err := r.db.
		Where("id = ? AND user_id = ?", notificationID, userID).
		First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Notification{}, ErrNotificationNotFound
		}
		return models.Notification{}, fmt.Errorf("find notification by id: %w", err)
	}

	return notification, nil
}

func (r Repository) Save(notification *models.Notification) error {
	if err := r.db.Save(notification).Error; err != nil {
		return fmt.Errorf("save notification: %w", err)
	}

	return nil
}

func nowUTC() time.Time {
	return time.Now().UTC()
}
