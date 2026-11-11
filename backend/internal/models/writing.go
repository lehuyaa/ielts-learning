package models

import (
	"time"

	"gorm.io/datatypes"
)

type WritingStatus string

const (
	WritingGraded          WritingStatus = "GRADED"
	WritingRewriteRequired WritingStatus = "REWRITE_REQUIRED"
	WritingCompleted       WritingStatus = "COMPLETED"
)

// WritingSubmission is one text written for a WRITING activity, its AI
// report (docs/foundation-spec.md §6.3) and the learner's rewrite.
type WritingSubmission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index:idx_writing_user_activity;not null" json:"userId"`
	ActivityID  uint           `gorm:"index:idx_writing_user_activity;index;not null" json:"activityId"`
	Text        string         `gorm:"type:text;not null" json:"text"`
	Status      WritingStatus  `gorm:"type:varchar(20);not null" json:"status"`
	Score       int            `gorm:"not null;default:0" json:"score"`
	ReportJSON  datatypes.JSON `gorm:"type:json" json:"report"`
	RewriteText *string        `gorm:"type:text" json:"rewriteText"`
	// RewriteReportJSON holds the mistakes still found in the last rewrite.
	RewriteReportJSON datatypes.JSON `gorm:"type:json" json:"rewriteReport"`
	RewriteCount      int            `gorm:"not null;default:0" json:"rewriteCount"`
	CompletedAt       *time.Time     `json:"completedAt"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

type MistakeSource string

const (
	MistakeFromExercise      MistakeSource = "EXERCISE"
	MistakeFromWriting       MistakeSource = "WRITING"
	MistakeFromSpeaking      MistakeSource = "SPEAKING"
	MistakeFromPronunciation MistakeSource = "PRONUNCIATION"
)

type MistakeStatus string

const (
	MistakeOpen      MistakeStatus = "OPEN"
	MistakeReviewing MistakeStatus = "REVIEWING"
	MistakeFixed     MistakeStatus = "FIXED"
)

// LearnerMistake is one distinct mistake of a learner, counted when it
// repeats (docs/foundation-spec.md §3.4). Reviewed in Task 73.
type LearnerMistake struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	UserID        uint          `gorm:"uniqueIndex:idx_learner_mistake;not null" json:"userId"`
	Source        MistakeSource `gorm:"type:varchar(20);not null" json:"source"`
	SkillTag      string        `gorm:"size:80;uniqueIndex:idx_learner_mistake;index;not null" json:"skillTag"`
	WrongText     string        `gorm:"size:500;not null" json:"wrongText"`
	CorrectText   string        `gorm:"size:500;uniqueIndex:idx_learner_mistake;not null" json:"correctText"`
	ExplanationVI string        `gorm:"column:explanation_vi;type:text" json:"explanationVi"`
	Occurrences   int           `gorm:"not null;default:1" json:"occurrences"`
	Status        MistakeStatus `gorm:"type:varchar(20);not null;default:'OPEN'" json:"status"`
	VocabularyID  *uint         `json:"vocabularyId"`
	// Review queue (Task 73): REVIEWING mistakes are shown as cards; three
	// correct reviews in a row mark them FIXED.
	ReviewStreak int        `gorm:"not null;default:0" json:"reviewStreak"`
	IntervalDays int        `gorm:"not null;default:0" json:"intervalDays"`
	NextReviewAt *time.Time `gorm:"index" json:"nextReviewAt"`
	FixedAt      *time.Time `json:"fixedAt"`
	LastSeenAt   time.Time  `json:"lastSeenAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
