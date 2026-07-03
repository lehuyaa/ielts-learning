package profile

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) FindUser(userID uint) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

func (r Repository) UsernameExistsForOtherUser(userID uint, username string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).
		Where("username = ? AND id <> ?", username, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check username exists for other user: %w", err)
	}

	return count > 0, nil
}

func (r Repository) SaveUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("save user: %w", err)
	}

	return nil
}

func (r Repository) CountLessonsDone(userID uint) (int, error) {
	var count int64
	if err := r.db.Model(&models.UserLessonProgress{}).
		Where("user_id = ? AND status = ?", userID, models.LessonStatusCompleted).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count lessons done: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountWordsByStatuses(userID uint, statuses ...models.VocabularyStatus) (int, error) {
	var count int64
	if err := r.db.Model(&models.UserVocabularyProgress{}).
		Where("user_id = ? AND status IN ?", userID, statuses).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count words by statuses: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountTrackedWords(userID uint) (int, error) {
	var count int64
	if err := r.db.Model(&models.UserVocabularyProgress{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count tracked words: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountPassedQuizzes(userID uint) (int, error) {
	var count int64
	if err := r.db.Model(&models.UserQuizAttempt{}).
		Where("user_id = ? AND passed = ?", userID, true).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count passed quizzes: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountCompletedLessonsByMinimumBand(userID uint, minBand float64) (int, error) {
	var count int64
	err := r.db.Model(&models.UserLessonProgress{}).
		Joins("JOIN lessons ON lessons.id = user_lesson_progresses.lesson_id").
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("user_lesson_progresses.user_id = ? AND user_lesson_progresses.status = ? AND band_levels.band_score >= ?",
			userID, models.LessonStatusCompleted, minBand).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count completed lessons by minimum band: %w", err)
	}

	return int(count), nil
}

func (r Repository) FindAchievements(userID uint) ([]models.Achievement, []models.UserAchievement, error) {
	var achievements []models.Achievement
	if err := r.db.
		Where("is_active = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&achievements).Error; err != nil {
		return nil, nil, fmt.Errorf("find achievements: %w", err)
	}

	var userAchievements []models.UserAchievement
	if err := r.db.
		Where("user_id = ?", userID).
		Find(&userAchievements).Error; err != nil {
		return nil, nil, fmt.Errorf("find user achievements: %w", err)
	}

	return achievements, userAchievements, nil
}

func (r Repository) FindRecentXPEvents(userID uint, limit int) ([]models.UserXPEvent, error) {
	var events []models.UserXPEvent
	if err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC, id DESC").
		Limit(limit).
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("find recent xp events: %w", err)
	}

	return events, nil
}

func (r Repository) FindDailyActivitiesSince(userID uint, from time.Time) ([]models.DailyActivity, error) {
	var activities []models.DailyActivity
	if err := r.db.
		Where("user_id = ? AND date >= ?", userID, from.Format("2006-01-02")).
		Order("date ASC").
		Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("find daily activities since: %w", err)
	}

	return activities, nil
}
