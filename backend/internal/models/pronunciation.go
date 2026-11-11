package models

import (
	"time"

	"gorm.io/datatypes"
)

// Where a pronunciation attempt was recorded (docs/foundation-spec.md §3.4).
const (
	PronunciationWord      = "WORD"
	PronunciationShadowing = "SHADOWING"
	PronunciationReadAloud = "READ_ALOUD"
	PronunciationExercise  = "EXERCISE"
)

// PronunciationAttempt is one scored recording (Task 71). The audio itself
// is scored and discarded.
type PronunciationAttempt struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"userId"`
	ReferenceText string         `gorm:"size:500;not null" json:"referenceText"`
	Source        string         `gorm:"type:varchar(30);not null" json:"source"`
	ExerciseID    *uint          `gorm:"index" json:"exerciseId"`
	Accuracy      int            `gorm:"not null" json:"accuracy"`
	Fluency       int            `gorm:"not null" json:"fluency"`
	Completeness  int            `gorm:"not null" json:"completeness"`
	Pronunciation int            `gorm:"not null" json:"pronunciation"`
	WordsJSON     datatypes.JSON `gorm:"type:json" json:"words"`
	CreatedAt     time.Time      `json:"createdAt"`
}
