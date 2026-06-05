package models

import (
	"time"

	"gorm.io/gorm"
)

type Lesson struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	TopicID              uint           `gorm:"uniqueIndex:idx_topic_lesson_slug;index;index:idx_topic_lesson_order;not null" json:"topicId"`
	Title                string         `gorm:"size:255;not null" json:"title"`
	Slug                 string         `gorm:"size:255;uniqueIndex:idx_topic_lesson_slug;not null" json:"slug"`
	Description          string         `gorm:"type:text" json:"description"`
	RequiredScore        int            `gorm:"default:80" json:"requiredScore"`
	EstimatedMinutes     int            `gorm:"default:10" json:"estimatedMinutes"`
	BandMin              *float64       `gorm:"type:decimal(3,1)" json:"bandMin"`
	BandMax              *float64       `gorm:"type:decimal(3,1)" json:"bandMax"`
	XPReward             int            `gorm:"default:0" json:"xpReward"`
	QuizTimeLimitSeconds *int           `json:"quizTimeLimitSeconds"`
	OrderIndex           int            `gorm:"default:0;index:idx_topic_lesson_order" json:"orderIndex"`
	IsPublished          bool           `gorm:"default:true" json:"isPublished"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	Topic              Topic              `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
	LessonVocabularies []LessonVocabulary `gorm:"foreignKey:LessonID" json:"lessonVocabularies,omitempty"`
	QuizQuestions      []QuizQuestion     `gorm:"foreignKey:LessonID" json:"quizQuestions,omitempty"`
}

type LessonVocabulary struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	LessonID     uint      `gorm:"uniqueIndex:idx_lesson_vocab;not null" json:"lessonId"`
	VocabularyID uint      `gorm:"uniqueIndex:idx_lesson_vocab;not null" json:"vocabularyId"`
	OrderIndex   int       `gorm:"default:0" json:"orderIndex"`
	IsRequired   bool      `gorm:"default:true" json:"isRequired"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	Lesson     Lesson     `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Vocabulary Vocabulary `gorm:"foreignKey:VocabularyID" json:"vocabulary,omitempty"`
}
