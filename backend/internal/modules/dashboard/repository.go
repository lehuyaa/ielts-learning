package dashboard

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

func (r Repository) FindUser(userID uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

func (r Repository) FindDailyActivity(userID uint, date string) (models.DailyActivity, bool, error) {
	var activity models.DailyActivity
	err := r.db.
		Where("user_id = ? AND date = ?", userID, date).
		First(&activity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DailyActivity{}, false, nil
		}
		return models.DailyActivity{}, false, fmt.Errorf("find daily activity: %w", err)
	}

	return activity, true, nil
}

func (r Repository) CountReviewDue(userID uint, now time.Time) (int, error) {
	var count int64
	err := r.db.Model(&models.UserVocabularyProgress{}).
		Where("user_id = ? AND next_review_at IS NOT NULL AND next_review_at <= ?", userID, now).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count review due: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountTotalWordsLearned(userID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.UserVocabularyProgress{}).
		Where(
			"user_id = ? AND status IN ?",
			userID,
			[]models.VocabularyStatus{
				models.VocabularyStatusLearning,
				models.VocabularyStatusReview,
				models.VocabularyStatusMastered,
			},
		).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count total words learned: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountMasteredWords(userID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.UserVocabularyProgress{}).
		Where("user_id = ? AND status = ?", userID, models.VocabularyStatusMastered).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count mastered words: %w", err)
	}

	return int(count), nil
}

func (r Repository) CountCompletedLessons(userID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.UserLessonProgress{}).
		Where("user_id = ? AND status = ?", userID, models.LessonStatusCompleted).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count completed lessons: %w", err)
	}

	return int(count), nil
}

func (r Repository) FindRecentXPEvents(userID uint, limit int) ([]models.UserXPEvent, error) {
	var events []models.UserXPEvent
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC, id DESC").
		Limit(limit).
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("find recent xp events: %w", err)
	}

	return events, nil
}
