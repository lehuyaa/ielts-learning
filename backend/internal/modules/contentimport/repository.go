package contentimport

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
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

// FindBandLevel finds a course level by course slug and CEFR level.
func (r Repository) FindBandLevel(courseSlug string, cefrLevel string) (models.BandLevel, error) {
	var band models.BandLevel
	err := r.db.
		Joins("JOIN courses ON courses.id = band_levels.course_id AND courses.deleted_at IS NULL").
		Where("courses.slug = ? AND band_levels.cefr_level = ?", courseSlug, cefrLevel).
		First(&band).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.BandLevel{}, err
	}
	if err != nil {
		return models.BandLevel{}, fmt.Errorf("find band level: %w", err)
	}
	return band, nil
}

func (r Repository) ExistingVocabularySlugs(slugs []string) (map[string]bool, error) {
	result := make(map[string]bool, len(slugs))
	if len(slugs) == 0 {
		return result, nil
	}
	var found []string
	if err := r.db.Model(&models.Vocabulary{}).Where("slug IN ?", slugs).Pluck("slug", &found).Error; err != nil {
		return nil, fmt.Errorf("find vocabulary slugs: %w", err)
	}
	for _, slug := range found {
		result[slug] = true
	}
	return result, nil
}

func (r Repository) FindTopic(bandLevelID uint, slug string) (models.Topic, bool, error) {
	var topic models.Topic
	err := r.db.Where("band_level_id = ? AND slug = ?", bandLevelID, slug).First(&topic).Error
	return topic, found(err), ignoreNotFound(err, "find topic")
}

func (r Repository) FindLesson(topicID uint, slug string) (models.Lesson, bool, error) {
	var lesson models.Lesson
	err := r.db.Where("topic_id = ? AND slug = ?", topicID, slug).First(&lesson).Error
	return lesson, found(err), ignoreNotFound(err, "find lesson")
}

func (r Repository) FindActivity(lessonID uint, slug string) (models.LessonActivity, bool, error) {
	var activity models.LessonActivity
	err := r.db.Where("lesson_id = ? AND slug = ?", lessonID, slug).First(&activity).Error
	return activity, found(err), ignoreNotFound(err, "find activity")
}

func (r Repository) ExerciseOrderIndexes(activityID uint) (map[int]uint, error) {
	var rows []models.Exercise
	if err := r.db.Select("id", "order_index").Where("activity_id = ?", activityID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("find exercises: %w", err)
	}
	result := make(map[int]uint, len(rows))
	for _, row := range rows {
		result[row.OrderIndex] = row.ID
	}
	return result, nil
}

func (r Repository) Save(value any) error {
	if err := r.db.Save(value).Error; err != nil {
		return fmt.Errorf("save %T: %w", value, err)
	}
	return nil
}

// SaveActivity saves an activity. GORM skips zero values that have a column
// default on insert (is_required=false would become true), so those fields
// are written explicitly afterwards.
func (r Repository) SaveActivity(activity *models.LessonActivity) error {
	// Captured first: on insert GORM fills skipped zero fields with their
	// tag default in the struct.
	flags := map[string]any{
		"is_required": activity.IsRequired,
		"pass_score":  activity.PassScore,
		"xp_reward":   activity.XPReward,
	}
	if err := r.Save(activity); err != nil {
		return err
	}
	activity.IsRequired = flags["is_required"].(bool)
	activity.PassScore = flags["pass_score"].(int)
	activity.XPReward = flags["xp_reward"].(int)
	err := r.db.Model(&models.LessonActivity{ID: activity.ID}).Updates(flags).Error
	if err != nil {
		return fmt.Errorf("save activity flags: %w", err)
	}
	return nil
}

// SaveExercise inserts a new exercise or overwrites an existing one in place
// (keeping its ID and created_at, so learner answers stay linked).
func (r Repository) SaveExercise(exercise *models.Exercise) error {
	if exercise.ID == 0 {
		if err := r.db.Create(exercise).Error; err != nil {
			return fmt.Errorf("create exercise: %w", err)
		}
		return nil
	}
	err := r.db.Model(&models.Exercise{ID: exercise.ID}).
		Select("type", "prompt", "payload_json", "answer_json", "explanation_vi", "skill_tag", "points", "order_index").
		Updates(exercise).Error
	if err != nil {
		return fmt.Errorf("update exercise: %w", err)
	}
	return nil
}

// DeleteExercisesAfter removes exercises that are no longer in the bundle.
func (r Repository) DeleteExercisesAfter(activityID uint, lastOrderIndex int) error {
	err := r.db.Where("activity_id = ? AND order_index > ?", activityID, lastOrderIndex).Delete(&models.Exercise{}).Error
	if err != nil {
		return fmt.Errorf("delete removed exercises: %w", err)
	}
	return nil
}

// UpdateCourseTotals refreshes the cached topic/lesson/word counts of a course.
func (r Repository) UpdateCourseTotals(courseID uint) error {
	var topics, lessons, words int64
	if err := r.db.Model(&models.Topic{}).
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).Count(&topics).Error; err != nil {
		return fmt.Errorf("count topics: %w", err)
	}
	if err := r.db.Model(&models.Lesson{}).
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).Count(&lessons).Error; err != nil {
		return fmt.Errorf("count lessons: %w", err)
	}
	if err := r.db.Model(&models.LessonVocabulary{}).
		Joins("JOIN lessons ON lessons.id = lesson_vocabularies.lesson_id").
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).Count(&words).Error; err != nil {
		return fmt.Errorf("count words: %w", err)
	}

	err := r.db.Model(&models.Course{}).Where("id = ?", courseID).Updates(map[string]any{
		"total_topics":  topics,
		"total_lessons": lessons,
		"total_words":   words,
	}).Error
	if err != nil {
		return fmt.Errorf("update course totals: %w", err)
	}
	return nil
}

func found(err error) bool {
	return err == nil
}

func ignoreNotFound(err error, action string) error {
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return fmt.Errorf("%s: %w", action, err)
}

func (r Repository) FindGrammarPoint(slug string) (models.GrammarPoint, bool, error) {
	var point models.GrammarPoint
	err := r.db.Where("slug = ?", slug).First(&point).Error
	return point, found(err), ignoreNotFound(err, "find grammar point")
}

func (r Repository) GrammarPointIDs(slugs []string) (map[string]uint, error) {
	result := map[string]uint{}
	if len(slugs) == 0 {
		return result, nil
	}
	var points []models.GrammarPoint
	if err := r.db.Select("id", "slug").Where("slug IN ?", slugs).Find(&points).Error; err != nil {
		return nil, fmt.Errorf("find grammar points: %w", err)
	}
	for _, point := range points {
		result[point.Slug] = point.ID
	}
	return result, nil
}

func (r Repository) FindPhoneme(symbol string) (models.Phoneme, bool, error) {
	var phoneme models.Phoneme
	err := r.db.Where("symbol = ?", symbol).First(&phoneme).Error
	return phoneme, found(err), ignoreNotFound(err, "find phoneme")
}

func (r Repository) ExistingPhonemeSymbols(symbols []string) (map[string]bool, error) {
	result := map[string]bool{}
	var rows []string
	if err := r.db.Model(&models.Phoneme{}).Where("symbol IN ?", symbols).Pluck("symbol", &rows).Error; err != nil {
		return nil, fmt.Errorf("find phonemes: %w", err)
	}
	for _, symbol := range rows {
		result[symbol] = true
	}
	return result, nil
}

func (r Repository) FindPassage(slug string) (models.Passage, bool, error) {
	var passage models.Passage
	err := r.db.Where("slug = ?", slug).First(&passage).Error
	return passage, found(err), ignoreNotFound(err, "find passage")
}

func (r Repository) PassageIDs(slugs []string) (map[string]uint, error) {
	result := map[string]uint{}
	if len(slugs) == 0 {
		return result, nil
	}
	var passages []models.Passage
	if err := r.db.Select("id", "slug").Where("slug IN ?", slugs).Find(&passages).Error; err != nil {
		return nil, fmt.Errorf("find passages: %w", err)
	}
	for _, passage := range passages {
		result[passage.Slug] = passage.ID
	}
	return result, nil
}

func (r Repository) FindSystemScenario(slug string) (models.AIConversationScenario, bool, error) {
	var scenario models.AIConversationScenario
	err := r.db.Where("slug = ? AND is_system = ?", slug, true).First(&scenario).Error
	return scenario, found(err), ignoreNotFound(err, "find scenario")
}

func (r Repository) SystemScenarioIDs(slugs []string) (map[string]uint, error) {
	result := map[string]uint{}
	if len(slugs) == 0 {
		return result, nil
	}
	var scenarios []models.AIConversationScenario
	if err := r.db.Select("id", "slug").Where("slug IN ? AND is_system = ?", slugs, true).Find(&scenarios).Error; err != nil {
		return nil, fmt.Errorf("find scenarios: %w", err)
	}
	for _, scenario := range scenarios {
		result[scenario.Slug] = scenario.ID
	}
	return result, nil
}
