package achievement

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

func (r Repository) FindActiveAchievementsByCodes(tx *gorm.DB, codes []string) ([]models.Achievement, error) {
	var achievements []models.Achievement
	if err := tx.
		Where("is_active = ? AND code IN ?", true, codes).
		Order("sort_order ASC, id ASC").
		Find(&achievements).Error; err != nil {
		return nil, fmt.Errorf("find active achievements: %w", err)
	}

	return achievements, nil
}

func (r Repository) FindAchievementByCode(tx *gorm.DB, code string) (models.Achievement, error) {
	var achievement models.Achievement
	if err := tx.Where("code = ? AND is_active = ?", code, true).First(&achievement).Error; err != nil {
		return models.Achievement{}, fmt.Errorf("find achievement by code: %w", err)
	}

	return achievement, nil
}

func (r Repository) HasAchievement(tx *gorm.DB, userID uint, achievementID uint) (bool, error) {
	var count int64
	if err := tx.Model(&models.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check user achievement: %w", err)
	}

	return count > 0, nil
}

func (r Repository) CreateUserAchievement(tx *gorm.DB, userAchievement models.UserAchievement) (bool, error) {
	result := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "achievement_id"},
		},
		DoNothing: true,
	}).Create(&userAchievement)
	if result.Error != nil {
		return false, fmt.Errorf("create user achievement: %w", result.Error)
	}

	return result.RowsAffected > 0, nil
}

func (r Repository) CreateNotification(tx *gorm.DB, notification models.Notification) error {
	if err := tx.Create(&notification).Error; err != nil {
		return fmt.Errorf("create achievement notification: %w", err)
	}

	return nil
}

func (r Repository) FindUser(tx *gorm.DB, userID uint) (models.User, error) {
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		return models.User{}, fmt.Errorf("find achievement user: %w", err)
	}

	return user, nil
}

func (r Repository) CountCompletedLessons(tx *gorm.DB, userID uint) (int, error) {
	var count int64
	if err := tx.Model(&models.UserLessonProgress{}).
		Where("user_id = ? AND status = ?", userID, models.LessonStatusCompleted).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count completed lessons: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountLearnedWords(tx *gorm.DB, userID uint) (int, error) {
	var count int64
	if err := tx.Model(&models.UserVocabularyProgress{}).
		Where("user_id = ? AND status IN ?", userID, []models.VocabularyStatus{
			models.VocabularyStatusLearning,
			models.VocabularyStatusReview,
			models.VocabularyStatusMastered,
		}).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count learned words: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountLessonsForTopicSlug(tx *gorm.DB, topicSlug string) (int, error) {
	var count int64
	if err := tx.Model(&models.Lesson{}).
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Where("topics.slug = ? AND lessons.is_published = ?", topicSlug, true).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count lessons for topic slug: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountCompletedLessonsForTopicSlug(tx *gorm.DB, userID uint, topicSlug string) (int, error) {
	var count int64
	if err := tx.Model(&models.UserLessonProgress{}).
		Joins("JOIN lessons ON lessons.id = user_lesson_progresses.lesson_id").
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Where("user_lesson_progresses.user_id = ? AND user_lesson_progresses.status = ? AND topics.slug = ?",
			userID, models.LessonStatusCompleted, topicSlug).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count completed lessons for topic slug: %w", err)
	}

	return int(count), nil
}

func nowUTC() time.Time {
	return time.Now().UTC()
}
