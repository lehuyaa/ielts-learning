package lessonactivity

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
	achievementmodule "ielts-learning/backend/internal/modules/achievement"
	activitymodule "ielts-learning/backend/internal/modules/activity"
	"ielts-learning/backend/internal/modules/grammarprogress"
	notificationmodule "ielts-learning/backend/internal/modules/notification"
	xpmodule "ielts-learning/backend/internal/modules/xp"
)

var (
	ErrLessonNotFound   = errors.New("lesson not found")
	ErrActivityNotFound = errors.New("activity not found")
)

type Repository struct {
	db                 *gorm.DB
	xpService          xpmodule.Service
	achievementService achievementmodule.Service
	dailyActivity      activitymodule.Service
}

func NewRepository(db *gorm.DB) Repository {
	xpService := xpmodule.NewService(xpmodule.NewRepository(db))
	notificationService := notificationmodule.NewService(notificationmodule.NewRepository(db))
	return Repository{
		db:                 db,
		xpService:          xpService,
		achievementService: achievementmodule.NewService(achievementmodule.NewRepository(db), xpService, notificationService),
		dailyActivity:      activitymodule.NewService(activitymodule.NewRepository(db)),
	}
}

// Transaction runs fn with a repository bound to one database transaction.
func (r Repository) Transaction(fn func(tx Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepository := r
		txRepository.db = tx
		return fn(txRepository)
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

func (r Repository) FindVocabulariesBySlugs(slugs []string) ([]models.Vocabulary, error) {
	if len(slugs) == 0 {
		return []models.Vocabulary{}, nil
	}

	var vocabularies []models.Vocabulary
	if err := r.db.Where("slug IN ?", slugs).Find(&vocabularies).Error; err != nil {
		return nil, fmt.Errorf("find vocabularies: %w", err)
	}

	return vocabularies, nil
}

// FindLessonProgressForUpdate returns the user's lesson progress (locked), or
// a new unsaved row when the user has not started the lesson.
func (r Repository) FindLessonProgressForUpdate(userID uint, lessonID uint) (models.UserLessonProgress, bool, error) {
	var progress models.UserLessonProgress
	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND lesson_id = ?", userID, lessonID).
		First(&progress).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.UserLessonProgress{UserID: userID, LessonID: lessonID}, false, nil
	}
	if err != nil {
		return models.UserLessonProgress{}, false, fmt.Errorf("find lesson progress: %w", err)
	}

	return progress, true, nil
}

func (r Repository) SaveLessonProgress(progress *models.UserLessonProgress) error {
	if err := r.db.Save(progress).Error; err != nil {
		return fmt.Errorf("save lesson progress: %w", err)
	}
	return nil
}

// UnlockNextLesson creates an UNLOCKED progress row for the next lesson of
// the topic, mirroring what a passed quiz does.
func (r Repository) UnlockNextLesson(userID uint, lesson models.Lesson, now time.Time) error {
	var nextLesson models.Lesson
	err := r.db.
		Where("topic_id = ? AND is_published = ? AND (order_index > ? OR (order_index = ? AND id > ?))", lesson.TopicID, true, lesson.OrderIndex, lesson.OrderIndex, lesson.ID).
		Order("order_index ASC, id ASC").
		First(&nextLesson).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("find next lesson: %w", err)
	}

	var count int64
	if err := r.db.Model(&models.UserLessonProgress{}).
		Where("user_id = ? AND lesson_id = ?", userID, nextLesson.ID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("find next lesson progress: %w", err)
	}
	if count > 0 {
		return nil
	}

	progress := models.UserLessonProgress{
		UserID:        userID,
		LessonID:      nextLesson.ID,
		Status:        models.LessonStatusUnlocked,
		LastStudiedAt: &now,
	}
	if err := r.db.Create(&progress).Error; err != nil {
		return fmt.Errorf("unlock next lesson: %w", err)
	}
	return nil
}

func (r Repository) AwardXP(input xpmodule.AwardInput) (bool, error) {
	result, err := r.xpService.AwardXP(r.db, input)
	if err != nil {
		return false, err
	}
	return result.EventCreated, nil
}

func (r Repository) RecordDailyActivity(userID uint, event activitymodule.Event) error {
	return r.dailyActivity.RecordActivityTx(r.db, userID, event)
}

func (r Repository) CheckAchievements(userID uint) error {
	_, err := r.achievementService.CheckAndUnlockAchievementsTx(r.db, userID)
	return err
}

func (r Repository) RecordGrammarAnswer(userID uint, skillTag string, scoreRatio float64, isCorrect bool, at time.Time) error {
	return grammarprogress.Record(r.db, userID, skillTag, scoreRatio, isCorrect, at)
}

func (r Repository) FindGrammarPoint(id uint) (models.GrammarPoint, bool, error) {
	var point models.GrammarPoint
	err := r.db.First(&point, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.GrammarPoint{}, false, nil
	}
	if err != nil {
		return models.GrammarPoint{}, false, fmt.Errorf("find grammar point: %w", err)
	}
	return point, true, nil
}

func (r Repository) FindPhonemes(symbols []string) ([]models.Phoneme, error) {
	var phonemes []models.Phoneme
	if err := r.db.Where("symbol IN ?", symbols).Find(&phonemes).Error; err != nil {
		return nil, fmt.Errorf("find phonemes: %w", err)
	}
	return phonemes, nil
}
