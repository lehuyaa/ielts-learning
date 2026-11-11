package models

import (
	"time"

	"gorm.io/datatypes"
)

// Where a practice set's exercises come from (Task 74).
const (
	PracticeFromAI   = "AI"
	PracticeFromBank = "BANK"
)

// PracticeSet is a learner's extra practice for one weakness (skill tag),
// built once per day: AI-generated exercises or ones picked from lessons.
type PracticeSet struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"uniqueIndex:idx_practice_set_day;not null" json:"userId"`
	SkillTag    string         `gorm:"size:80;uniqueIndex:idx_practice_set_day;not null" json:"skillTag"`
	Day         string         `gorm:"size:10;uniqueIndex:idx_practice_set_day;not null" json:"day"`
	Source      string         `gorm:"type:varchar(10);not null" json:"source"`
	ItemsJSON   datatypes.JSON `gorm:"type:json;not null" json:"items"`
	ResultsJSON datatypes.JSON `gorm:"type:json" json:"results"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}
