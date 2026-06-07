package flashcard

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrLessonNotFound     = errors.New("lesson not found")
	ErrVocabularyNotFound = errors.New("vocabulary not found")
)

type VocabularyTopicContext struct {
	VocabularyID uint
	TopicTitle   string
}

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

func (r Repository) FindLessonVocabularies(lessonID uint) ([]models.LessonVocabulary, error) {
	var rows []models.LessonVocabulary
	err := r.db.
		Preload("Vocabulary").
		Where("lesson_id = ? AND is_required = ?", lessonID, true).
		Order("order_index ASC, id ASC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("find lesson vocabularies: %w", err)
	}

	return rows, nil
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

func (r Repository) FindDueProgress(userID uint, query DueReviewQuery, now time.Time) ([]models.UserVocabularyProgress, error) {
	db := r.db.
		Model(&models.UserVocabularyProgress{}).
		Distinct("user_vocabulary_progresses.*").
		Preload("Vocabulary").
		Where("user_vocabulary_progresses.user_id = ?", userID).
		Where("user_vocabulary_progresses.next_review_at IS NOT NULL").
		Where("user_vocabulary_progresses.next_review_at <= ?", now)

	if query.TopicID != nil {
		db = db.
			Joins("JOIN lesson_vocabularies ON lesson_vocabularies.vocabulary_id = user_vocabulary_progresses.vocabulary_id").
			Joins("JOIN lessons ON lessons.id = lesson_vocabularies.lesson_id").
			Where("lessons.topic_id = ?", *query.TopicID)
	}

	var progress []models.UserVocabularyProgress
	err := db.
		Order("user_vocabulary_progresses.next_review_at ASC, user_vocabulary_progresses.id ASC").
		Limit(query.Limit).
		Find(&progress).Error
	if err != nil {
		return nil, fmt.Errorf("find due progress: %w", err)
	}

	return progress, nil
}

func (r Repository) FindTopicContext(vocabularyIDs []uint) (map[uint]VocabularyTopicContext, error) {
	contexts := make(map[uint]VocabularyTopicContext)
	if len(vocabularyIDs) == 0 {
		return contexts, nil
	}

	var rows []VocabularyTopicContext
	err := r.db.
		Table("lesson_vocabularies").
		Select(`
			lesson_vocabularies.vocabulary_id AS vocabulary_id,
			topics.title AS topic_title
		`).
		Joins("JOIN lessons ON lessons.id = lesson_vocabularies.lesson_id").
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Where("lesson_vocabularies.vocabulary_id IN ?", vocabularyIDs).
		Where("lesson_vocabularies.is_required = ?", true).
		Order("lesson_vocabularies.vocabulary_id ASC, lessons.order_index ASC, lesson_vocabularies.order_index ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("find topic context: %w", err)
	}

	for _, row := range rows {
		if _, exists := contexts[row.VocabularyID]; exists {
			continue
		}
		contexts[row.VocabularyID] = row
	}

	return contexts, nil
}

type ReviewUpdate struct {
	Progress  models.UserVocabularyProgress
	XPAwarded int
}

func (r Repository) SaveReview(
	user models.User,
	vocabularyID uint,
	now time.Time,
	update func(models.UserVocabularyProgress, bool) ReviewUpdate,
) (models.UserVocabularyProgress, models.User, int, error) {
	var progress models.UserVocabularyProgress
	updatedUser := user
	xpAwarded := 0

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var vocabulary models.Vocabulary
		if err := tx.First(&vocabulary, vocabularyID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrVocabularyNotFound
			}
			return fmt.Errorf("find review vocabulary: %w", err)
		}

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND vocabulary_id = ?", user.ID, vocabularyID).
			First(&progress).Error
		exists := true
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("find review progress: %w", err)
			}
			exists = false
			progress = models.UserVocabularyProgress{
				UserID:       user.ID,
				VocabularyID: vocabularyID,
				Status:       models.VocabularyStatusNew,
			}
		}

		result := update(progress, exists)
		progress = result.Progress
		xpAwarded = result.XPAwarded

		if exists {
			if err := tx.Save(&progress).Error; err != nil {
				return fmt.Errorf("update review progress: %w", err)
			}
		} else {
			if err := tx.Create(&progress).Error; err != nil {
				return fmt.Errorf("create review progress: %w", err)
			}
		}

		if xpAwarded > 0 {
			sourceID := vocabularyID
			xpEvent := models.UserXPEvent{
				UserID:      user.ID,
				SourceType:  "FLASHCARD_REVIEW",
				SourceID:    &sourceID,
				XP:          xpAwarded,
				Description: "Flashcard review",
				CreatedAt:   now,
			}
			if err := tx.Create(&xpEvent).Error; err != nil {
				return fmt.Errorf("create xp event: %w", err)
			}

			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&updatedUser, user.ID).Error; err != nil {
				return fmt.Errorf("find xp user: %w", err)
			}
			updatedUser.TotalXP += xpAwarded
			updatedUser.Level = updatedUser.TotalXP/200 + 1
			updatedUser.LastActiveAt = &now
			if err := tx.Save(&updatedUser).Error; err != nil {
				return fmt.Errorf("update xp user: %w", err)
			}
		}

		if err := upsertDailyActivity(tx, user.ID, activityDate(now, user.Timezone), xpAwarded); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.UserVocabularyProgress{}, models.User{}, 0, err
	}

	return progress, updatedUser, xpAwarded, nil
}

func upsertDailyActivity(tx *gorm.DB, userID uint, date time.Time, xpAwarded int) error {
	var activity models.DailyActivity
	err := tx.
		Where("user_id = ? AND date = ?", userID, date).
		First(&activity).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("find daily activity: %w", err)
		}

		activity = models.DailyActivity{
			UserID:        userID,
			Date:          date,
			WordsReviewed: 1,
			XPEarned:      xpAwarded,
		}
		if err := tx.Create(&activity).Error; err != nil {
			return fmt.Errorf("create daily activity: %w", err)
		}
		return nil
	}

	activity.WordsReviewed++
	activity.XPEarned += xpAwarded
	if err := tx.Save(&activity).Error; err != nil {
		return fmt.Errorf("update daily activity: %w", err)
	}

	return nil
}

func activityDate(now time.Time, timezone string) time.Time {
	location := time.UTC
	if timezone != "" {
		if loadedLocation, err := time.LoadLocation(timezone); err == nil {
			location = loadedLocation
		}
	}

	local := now.In(location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
}
