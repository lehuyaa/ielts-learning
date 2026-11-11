package models

import (
	"time"

	"gorm.io/datatypes"
)

// AI usage kinds counted per user per day (docs/foundation-spec.md §6.7).
const (
	AIUsageShortAnswer = "SHORT_ANSWER"
	AIUsageWriting     = "WRITING"
	// AIUsagePronunciation counts Azure pronunciation assessments (Task 71).
	AIUsagePronunciation = "PRONUNCIATION"
	// AIUsagePractice counts AI-generated practice sets (Task 74).
	AIUsagePractice = "PRACTICE"
)

// AIUsageDaily counts AI calls per user, day (YYYY-MM-DD, server time) and kind.
type AIUsageDaily struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_ai_usage_user_day_kind;not null" json:"userId"`
	Day       string    `gorm:"uniqueIndex:idx_ai_usage_user_day_kind;size:10;not null" json:"day"`
	Kind      string    `gorm:"uniqueIndex:idx_ai_usage_user_day_kind;size:32;not null" json:"kind"`
	Count     int       `gorm:"not null;default:0" json:"count"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (AIUsageDaily) TableName() string { return "ai_usage_dailies" }

// AIShortAnswerGrading caches the AI verdict for one exercise and one
// normalised answer, so the same answer is never graded twice.
type AIShortAnswerGrading struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ExerciseID uint           `gorm:"uniqueIndex:idx_ai_short_answer;not null" json:"exerciseId"`
	AnswerHash string         `gorm:"uniqueIndex:idx_ai_short_answer;size:64;not null" json:"answerHash"`
	AnswerText string         `gorm:"type:text;not null" json:"answerText"`
	ResultJSON datatypes.JSON `gorm:"type:json;not null" json:"result"`
	Model      string         `gorm:"size:64" json:"model"`
	CreatedAt  time.Time      `json:"createdAt"`
}
