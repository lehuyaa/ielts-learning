package roadmap

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var ErrRoadmapNotFound = errors.New("roadmap not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) FindPublishedCourse() (models.Course, error) {
	var course models.Course
	err := r.db.
		Where("is_published = ?", true).
		Order("order_index ASC, id ASC").
		First(&course).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Course{}, ErrRoadmapNotFound
		}
		return models.Course{}, fmt.Errorf("find published course: %w", err)
	}

	return course, nil
}

func (r Repository) FindBandLevels(courseID uint) ([]models.BandLevel, error) {
	var bands []models.BandLevel
	err := r.db.
		Where("course_id = ?", courseID).
		Order("order_index ASC, band_score ASC, id ASC").
		Find(&bands).Error
	if err != nil {
		return nil, fmt.Errorf("find band levels: %w", err)
	}

	return bands, nil
}

func (r Repository) FindTopics(bandLevelIDs []uint) ([]models.Topic, error) {
	if len(bandLevelIDs) == 0 {
		return []models.Topic{}, nil
	}

	var topics []models.Topic
	err := r.db.
		Where("band_level_id IN ?", bandLevelIDs).
		Order("order_index ASC, id ASC").
		Find(&topics).Error
	if err != nil {
		return nil, fmt.Errorf("find topics: %w", err)
	}

	return topics, nil
}

func (r Repository) FindLessons(topicIDs []uint) ([]models.Lesson, error) {
	if len(topicIDs) == 0 {
		return []models.Lesson{}, nil
	}

	var lessons []models.Lesson
	err := r.db.
		Where("topic_id IN ? AND is_published = ?", topicIDs, true).
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
