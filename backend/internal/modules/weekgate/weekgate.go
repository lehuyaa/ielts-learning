// Package weekgate links the weeks of a course (Task 75): the first lesson of
// a CORE topic opens only once the previous topic of the same course and
// track is finished — its last lesson, the week's checkpoint, is completed.
package weekgate

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

// LockedReason is shown on a week's first lesson while the previous week is open.
const LockedReason = "Pass the previous week's checkpoint to unlock this week."

// trackTopics lists the topics of the topic's course and track in roadmap
// order (bands, then topics).
func trackTopics(db *gorm.DB, topicID uint) ([]models.Topic, int, error) {
	var topic models.Topic
	if err := db.Preload("BandLevel").Take(&topic, topicID).Error; err != nil {
		return nil, -1, fmt.Errorf("find topic: %w", err)
	}
	var topics []models.Topic
	err := db.Model(&models.Topic{}).
		Joins("JOIN band_levels bl ON bl.id = topics.band_level_id").
		Where("bl.course_id = ? AND topics.track = ?", topic.BandLevel.CourseID, topic.Track).
		Order("bl.order_index ASC, bl.id ASC, topics.order_index ASC, topics.id ASC").
		Find(&topics).Error
	if err != nil {
		return nil, -1, fmt.Errorf("find track topics: %w", err)
	}
	for i, candidate := range topics {
		if candidate.ID == topic.ID {
			return topics, i, nil
		}
	}
	return topics, -1, nil
}

func lastLesson(db *gorm.DB, topicID uint) (models.Lesson, bool, error) {
	var lesson models.Lesson
	err := db.Where("topic_id = ? AND is_published = ?", topicID, true).
		Order("order_index DESC, id DESC").Take(&lesson).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return lesson, false, nil
	}
	if err != nil {
		return lesson, false, fmt.Errorf("find last lesson: %w", err)
	}
	return lesson, true, nil
}

// PreviousWeekDone reports whether the topic before this one in its track is
// finished (true for the first topic and for always-unlocked tracks).
func PreviousWeekDone(db *gorm.DB, userID uint, topic models.Topic) (bool, error) {
	if topic.IsAlwaysUnlocked() {
		return true, nil
	}
	topics, index, err := trackTopics(db, topic.ID)
	if err != nil || index <= 0 {
		return true, err
	}
	for i := index - 1; i >= 0; i-- {
		last, found, err := lastLesson(db, topics[i].ID)
		if err != nil {
			return false, err
		}
		if !found {
			continue // an empty topic does not block the path
		}
		var count int64
		if err := db.Model(&models.UserLessonProgress{}).
			Where("user_id = ? AND lesson_id = ? AND status = ?", userID, last.ID, models.LessonStatusCompleted).
			Count(&count).Error; err != nil {
			return false, fmt.Errorf("find lesson progress: %w", err)
		}
		return count > 0, nil
	}
	return true, nil
}

// NextWeekFirstLesson returns the first lesson of the topic after this one in
// its track, if any.
func NextWeekFirstLesson(db *gorm.DB, topicID uint) (models.Lesson, bool, error) {
	topics, index, err := trackTopics(db, topicID)
	if err != nil || index < 0 {
		return models.Lesson{}, false, err
	}
	for _, next := range topics[index+1:] {
		var lesson models.Lesson
		err := db.Where("topic_id = ? AND is_published = ?", next.ID, true).
			Order("order_index ASC, id ASC").Take(&lesson).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return lesson, false, fmt.Errorf("find next week lesson: %w", err)
		}
		return lesson, true, nil
	}
	return models.Lesson{}, false, nil
}
