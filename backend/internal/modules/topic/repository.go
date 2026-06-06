package topic

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var ErrTopicNotFound = errors.New("topic not found")

type VocabularyCount struct {
	LessonID uint
	Count    int
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

func (r Repository) FindTopic(topicID uint) (models.Topic, error) {
	var topic models.Topic
	err := r.db.
		Preload("BandLevel").
		First(&topic, topicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Topic{}, ErrTopicNotFound
		}
		return models.Topic{}, fmt.Errorf("find topic: %w", err)
	}

	return topic, nil
}

func (r Repository) FindLessons(topicID uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.
		Where("topic_id = ? AND is_published = ?", topicID, true).
		Order("order_index ASC, id ASC").
		Find(&lessons).Error
	if err != nil {
		return nil, fmt.Errorf("find lessons: %w", err)
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
		return nil, fmt.Errorf("find user lesson progress: %w", err)
	}

	return progress, nil
}

func (r Repository) CountVocabularyByLesson(lessonIDs []uint) (map[uint]int, error) {
	counts := make(map[uint]int)
	if len(lessonIDs) == 0 {
		return counts, nil
	}

	var rows []VocabularyCount
	err := r.db.
		Model(&models.LessonVocabulary{}).
		Select("lesson_id, COUNT(*) AS count").
		Where("lesson_id IN ? AND is_required = ?", lessonIDs, true).
		Group("lesson_id").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("count lesson vocabularies: %w", err)
	}

	for _, row := range rows {
		counts[row.LessonID] = row.Count
	}

	return counts, nil
}
