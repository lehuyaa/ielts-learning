package lesson

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var (
	ErrLessonNotFound = errors.New("lesson not found")
	ErrUserNotFound   = errors.New("user not found")
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

func (r Repository) FindLesson(lessonID uint) (models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.
		Preload("Topic.BandLevel").
		Where("is_published = ?", true).
		First(&lesson, lessonID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Lesson{}, ErrLessonNotFound
		}
		return models.Lesson{}, fmt.Errorf("find lesson: %w", err)
	}

	return lesson, nil
}

func (r Repository) FindTopicLessons(topicID uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.
		Where("topic_id = ? AND is_published = ?", topicID, true).
		Order("order_index ASC, id ASC").
		Find(&lessons).Error
	if err != nil {
		return nil, fmt.Errorf("find topic lessons: %w", err)
	}

	return lessons, nil
}

func (r Repository) FindLessonVocabularies(lessonID uint) ([]models.Vocabulary, error) {
	var words []models.Vocabulary
	err := r.db.Model(&models.Vocabulary{}).
		Joins("JOIN lesson_vocabularies ON lesson_vocabularies.vocabulary_id = vocabularies.id").
		Where("lesson_vocabularies.lesson_id = ? AND lesson_vocabularies.is_required = ?", lessonID, true).
		Order("lesson_vocabularies.order_index ASC, vocabularies.id ASC").
		Find(&words).Error
	if err != nil {
		return nil, fmt.Errorf("find lesson vocabularies: %w", err)
	}

	return words, nil
}

func (r Repository) FindLessonProgress(userID uint, lessonIDs []uint) ([]models.UserLessonProgress, error) {
	if len(lessonIDs) == 0 {
		return []models.UserLessonProgress{}, nil
	}

	var progress []models.UserLessonProgress
	err := r.db.
		Where("user_id = ? AND lesson_id IN ?", userID, lessonIDs).
		Find(&progress).Error
	if err != nil {
		return nil, fmt.Errorf("find lesson progress: %w", err)
	}

	return progress, nil
}

func (r Repository) FindVocabularyProgress(userID uint, vocabularyIDs []uint) ([]models.UserVocabularyProgress, error) {
	if len(vocabularyIDs) == 0 {
		return []models.UserVocabularyProgress{}, nil
	}

	var progress []models.UserVocabularyProgress
	err := r.db.
		Where("user_id = ? AND vocabulary_id IN ?", userID, vocabularyIDs).
		Find(&progress).Error
	if err != nil {
		return nil, fmt.Errorf("find vocabulary progress: %w", err)
	}

	return progress, nil
}

func (r Repository) StartLesson(userID uint, lesson models.Lesson, wordCount int, now time.Time) (models.UserLessonProgress, error) {
	var progress models.UserLessonProgress
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where("user_id = ? AND lesson_id = ?", userID, lesson.ID).
			First(&progress).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("find start progress: %w", err)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			progress = models.UserLessonProgress{
				UserID:        userID,
				LessonID:      lesson.ID,
				Status:        models.LessonStatusInProgress,
				StartedAt:     &now,
				LastStudiedAt: &now,
				TotalWords:    wordCount,
			}
			if err := tx.Create(&progress).Error; err != nil {
				return fmt.Errorf("create lesson progress: %w", err)
			}
			return nil
		}

		if progress.Status != models.LessonStatusCompleted {
			progress.Status = models.LessonStatusInProgress
		}
		if progress.StartedAt == nil {
			progress.StartedAt = &now
		}
		progress.LastStudiedAt = &now
		progress.TotalWords = wordCount

		if err := tx.Save(&progress).Error; err != nil {
			return fmt.Errorf("update lesson progress: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.UserLessonProgress{}, err
	}

	return progress, nil
}
