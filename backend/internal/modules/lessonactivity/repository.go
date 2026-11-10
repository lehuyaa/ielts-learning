package lessonactivity

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

var (
	ErrLessonNotFound   = errors.New("lesson not found")
	ErrActivityNotFound = errors.New("activity not found")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Transaction(fn func(tx Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(Repository{db: tx})
	})
}

func (r Repository) FindLesson(lessonID uint) (models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.
		Preload("Topic").
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

	var rows []models.UserLessonProgress
	if err := r.db.Where("user_id = ? AND lesson_id IN ?", userID, lessonIDs).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("find lesson progress: %w", err)
	}

	return rows, nil
}

func (r Repository) FindActivity(activityID uint) (models.LessonActivity, error) {
	var activity models.LessonActivity
	err := r.db.
		Where("status = ?", models.ContentPublished).
		First(&activity, activityID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LessonActivity{}, ErrActivityNotFound
		}
		return models.LessonActivity{}, fmt.Errorf("find activity: %w", err)
	}

	return activity, nil
}

func (r Repository) FindLessonActivities(lessonID uint) ([]models.LessonActivity, error) {
	var activities []models.LessonActivity
	err := r.db.
		Where("lesson_id = ? AND status = ?", lessonID, models.ContentPublished).
		Order("order_index ASC, id ASC").
		Find(&activities).Error
	if err != nil {
		return nil, fmt.Errorf("find lesson activities: %w", err)
	}

	return activities, nil
}

type exerciseCount struct {
	ActivityID uint
	Total      int
}

func (r Repository) CountExercises(activityIDs []uint) (map[uint]int, error) {
	result := make(map[uint]int, len(activityIDs))
	if len(activityIDs) == 0 {
		return result, nil
	}

	var rows []exerciseCount
	err := r.db.Model(&models.Exercise{}).
		Select("activity_id, COUNT(*) AS total").
		Where("activity_id IN ?", activityIDs).
		Group("activity_id").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("count exercises: %w", err)
	}
	for _, row := range rows {
		result[row.ActivityID] = row.Total
	}

	return result, nil
}

func (r Repository) FindExercises(activityID uint) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.
		Where("activity_id = ?", activityID).
		Order("order_index ASC, id ASC").
		Find(&exercises).Error
	if err != nil {
		return nil, fmt.Errorf("find exercises: %w", err)
	}

	return exercises, nil
}

func (r Repository) FindActivityProgress(userID uint, activityIDs []uint) ([]models.UserActivityProgress, error) {
	if len(activityIDs) == 0 {
		return []models.UserActivityProgress{}, nil
	}

	var rows []models.UserActivityProgress
	if err := r.db.Where("user_id = ? AND activity_id IN ?", userID, activityIDs).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("find activity progress: %w", err)
	}

	return rows, nil
}

// FindActivityProgressForUpdate locks the row inside a transaction so two
// concurrent answers cannot start two attempts.
func (r Repository) FindActivityProgressForUpdate(userID uint, activityID uint) (models.UserActivityProgress, bool, error) {
	var progress models.UserActivityProgress
	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND activity_id = ?", userID, activityID).
		First(&progress).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.UserActivityProgress{}, false, nil
	}
	if err != nil {
		return models.UserActivityProgress{}, false, fmt.Errorf("find activity progress: %w", err)
	}

	return progress, true, nil
}

func (r Repository) SaveActivityProgress(progress *models.UserActivityProgress) error {
	if err := r.db.Save(progress).Error; err != nil {
		return fmt.Errorf("save activity progress: %w", err)
	}
	return nil
}

func (r Repository) FindAnswers(userID uint, exerciseIDs []uint, attempt int) ([]models.UserExerciseAnswer, error) {
	if len(exerciseIDs) == 0 || attempt == 0 {
		return []models.UserExerciseAnswer{}, nil
	}

	var rows []models.UserExerciseAnswer
	err := r.db.
		Where("user_id = ? AND exercise_id IN ? AND activity_attempt = ?", userID, exerciseIDs, attempt).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("find exercise answers: %w", err)
	}

	return rows, nil
}

// UpsertAnswer keeps the latest answer per exercise and attempt.
func (r Repository) UpsertAnswer(answer *models.UserExerciseAnswer) error {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "exercise_id"}, {Name: "activity_attempt"}},
		DoUpdates: clause.AssignmentColumns([]string{"answer_json", "is_correct", "score_ratio", "answered_at", "updated_at"}),
	}).Create(answer).Error
	if err != nil {
		return fmt.Errorf("save exercise answer: %w", err)
	}
	return nil
}
