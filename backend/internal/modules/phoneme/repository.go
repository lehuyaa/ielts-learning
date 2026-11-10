package phoneme

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var ErrPhonemeNotFound = errors.New("phoneme not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

// FindAll returns every sound card in chart order.
func (r Repository) FindAll() ([]models.Phoneme, error) {
	var phonemes []models.Phoneme
	if err := r.db.Order("order_index ASC, id ASC").Find(&phonemes).Error; err != nil {
		return nil, fmt.Errorf("find phonemes: %w", err)
	}
	return phonemes, nil
}

func (r Repository) FindProgress(userID uint) (map[uint]models.UserPhonemeProgress, error) {
	var rows []models.UserPhonemeProgress
	if err := r.db.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("find phoneme progress: %w", err)
	}
	result := make(map[uint]models.UserPhonemeProgress, len(rows))
	for _, row := range rows {
		result[row.PhonemeID] = row
	}
	return result, nil
}

// FindLessons lists published PRONUNCIATION activities whose
// content.phonemes contains the symbol.
func (r Repository) FindLessons(symbol string) ([]LessonLink, error) {
	links := []LessonLink{}
	err := r.db.Raw(`SELECT l.id AS lesson_id, l.title AS lesson_title, t.title AS topic_title,
			a.id AS activity_id, a.title AS activity_title
		FROM lesson_activities a
		JOIN lessons l ON l.id = a.lesson_id AND l.deleted_at IS NULL
		JOIN topics t ON t.id = l.topic_id AND t.deleted_at IS NULL
		WHERE a.deleted_at IS NULL AND a.type = ? AND a.status = ?
			AND JSON_CONTAINS(JSON_EXTRACT(a.content_json, '$.phonemes'), JSON_QUOTE(?))
		ORDER BY t.order_index, l.order_index, a.order_index`,
		models.ActivityPronunciation, models.ContentPublished, symbol).Scan(&links).Error
	if err != nil {
		return nil, fmt.Errorf("find phoneme lessons: %w", err)
	}
	return links, nil
}

// FindPracticeTopics lists VOCABULARY-track topics (always unlocked) that
// contain pronunciation activities.
func (r Repository) FindPracticeTopics() ([]PracticeTopic, error) {
	topics := []PracticeTopic{}
	err := r.db.Raw(`SELECT t.id, t.title, t.description, COUNT(DISTINCT l.id) AS lesson_count
		FROM topics t
		JOIN lessons l ON l.topic_id = t.id AND l.deleted_at IS NULL
		JOIN lesson_activities a ON a.lesson_id = l.id AND a.deleted_at IS NULL
		WHERE t.deleted_at IS NULL AND t.track = ? AND a.type = ? AND a.status = ?
		GROUP BY t.id, t.title, t.description, t.order_index
		ORDER BY t.order_index, t.id`,
		models.TopicTrackVocabulary, models.ActivityPronunciation, models.ContentPublished).Scan(&topics).Error
	if err != nil {
		return nil, fmt.Errorf("find practice topics: %w", err)
	}
	return topics, nil
}
