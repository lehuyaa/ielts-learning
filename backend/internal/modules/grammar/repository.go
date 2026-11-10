package grammar

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

var ErrGrammarPointNotFound = errors.New("grammar point not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) DB() *gorm.DB {
	return r.db
}

func (r Repository) FindGrammarPoints(cefrLevel string) ([]models.GrammarPoint, error) {
	query := r.db.Order("cefr_level ASC, order_index ASC, id ASC")
	if cefrLevel != "" {
		query = query.Where("cefr_level = ?", cefrLevel)
	}
	var points []models.GrammarPoint
	if err := query.Find(&points).Error; err != nil {
		return nil, fmt.Errorf("find grammar points: %w", err)
	}
	return points, nil
}

func (r Repository) FindGrammarPoint(slug string) (models.GrammarPoint, error) {
	var point models.GrammarPoint
	err := r.db.Where("slug = ?", slug).First(&point).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.GrammarPoint{}, ErrGrammarPointNotFound
	}
	if err != nil {
		return models.GrammarPoint{}, fmt.Errorf("find grammar point: %w", err)
	}
	return point, nil
}

func (r Repository) FindProgress(userID uint) (map[uint]models.UserGrammarProgress, error) {
	var rows []models.UserGrammarProgress
	if err := r.db.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("find grammar progress: %w", err)
	}
	result := make(map[uint]models.UserGrammarProgress, len(rows))
	for _, row := range rows {
		result[row.GrammarPointID] = row
	}
	return result, nil
}

// availableExercisesQuery selects published exercises with a skill tag from
// lessons the learner has started (any lesson progress row).
func (r Repository) availableExercisesQuery(userID uint) *gorm.DB {
	return r.db.Model(&models.Exercise{}).
		Joins("JOIN lesson_activities ON lesson_activities.id = exercises.activity_id AND lesson_activities.deleted_at IS NULL AND lesson_activities.status = ?", models.ContentPublished).
		Joins("JOIN lessons ON lessons.id = lesson_activities.lesson_id AND lessons.deleted_at IS NULL AND lessons.is_published = ?", true).
		Joins("JOIN user_lesson_progresses ON user_lesson_progresses.lesson_id = lessons.id AND user_lesson_progresses.user_id = ?", userID)
}

type tagCount struct {
	SkillTag string
	Total    int
}

func (r Repository) CountAvailableExercisesByTag(userID uint) (map[string]int, error) {
	var rows []tagCount
	err := r.availableExercisesQuery(userID).
		Select("exercises.skill_tag AS skill_tag, COUNT(*) AS total").
		Where("exercises.skill_tag LIKE ?", "grammar:%").
		Group("exercises.skill_tag").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("count practice exercises: %w", err)
	}
	result := make(map[string]int, len(rows))
	for _, row := range rows {
		result[row.SkillTag] = row.Total
	}
	return result, nil
}

func (r Repository) FindAvailableExercises(userID uint, skillTag string) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.availableExercisesQuery(userID).
		Where("exercises.skill_tag = ?", skillTag).
		Order("lessons.id ASC, exercises.order_index ASC").
		Find(&exercises).Error
	if err != nil {
		return nil, fmt.Errorf("find practice exercises: %w", err)
	}
	return exercises, nil
}

// LatestAnswers returns each exercise's most recent answer by the user,
// across lesson attempts and practice.
func (r Repository) LatestAnswers(userID uint, exerciseIDs []uint) (map[uint]models.UserExerciseAnswer, error) {
	result := map[uint]models.UserExerciseAnswer{}
	if len(exerciseIDs) == 0 {
		return result, nil
	}
	var rows []models.UserExerciseAnswer
	err := r.db.Where("user_id = ? AND exercise_id IN ?", userID, exerciseIDs).
		Order("answered_at ASC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("find answers: %w", err)
	}
	for _, row := range rows {
		result[row.ExerciseID] = row
	}
	return result, nil
}

// SavePracticeAnswer stores a practice answer as attempt 0 (outside lessons).
func (r Repository) SavePracticeAnswer(answer *models.UserExerciseAnswer) error {
	answer.ActivityAttempt = 0
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "exercise_id"}, {Name: "activity_attempt"}},
		DoUpdates: clause.AssignmentColumns([]string{"answer_json", "is_correct", "score_ratio", "answered_at", "updated_at"}),
	}).Create(answer).Error
	if err != nil {
		return fmt.Errorf("save practice answer: %w", err)
	}
	return nil
}

func (r Repository) FindProgressFor(userID uint, grammarPointID uint) (models.UserGrammarProgress, bool, error) {
	var progress models.UserGrammarProgress
	err := r.db.Where("user_id = ? AND grammar_point_id = ?", userID, grammarPointID).First(&progress).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.UserGrammarProgress{}, false, nil
	}
	if err != nil {
		return models.UserGrammarProgress{}, false, fmt.Errorf("find grammar progress: %w", err)
	}
	return progress, true, nil
}
