package models

import (
	"time"

	"gorm.io/datatypes"
)

// PhonemeKind groups sound cards on the IPA chart.
type PhonemeKind string

const (
	PhonemeVowel     PhonemeKind = "VOWEL"
	PhonemeDiphthong PhonemeKind = "DIPHTHONG"
	PhonemeConsonant PhonemeKind = "CONSONANT"
)

// Phoneme is a sound card (docs/foundation-spec.md §3.3). Symbol is the IPA
// as written in our vocabulary data (American English, Oxford NAmE style).
type Phoneme struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Symbol           string         `gorm:"size:10;uniqueIndex;not null" json:"symbol"`
	Kind             PhonemeKind    `gorm:"type:varchar(20);not null" json:"kind"`
	Voiced           *bool          `json:"voiced"`
	NameVI           string         `gorm:"column:name_vi;size:255" json:"nameVi"`
	ExampleWordsJSON datatypes.JSON `gorm:"type:json" json:"exampleWords"`
	MinimalPairsJSON datatypes.JSON `gorm:"type:json" json:"minimalPairs"`
	MouthVI          string         `gorm:"column:mouth_vi;type:text" json:"mouthVi"`
	MistakesVI       string         `gorm:"column:mistakes_vi;type:text" json:"mistakesVi"`
	AudioURL         string         `gorm:"size:500" json:"audioUrl"`
	ImageURL         *string        `gorm:"size:500" json:"imageUrl"`
	OrderIndex       int            `gorm:"default:0" json:"orderIndex"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

// UserPhonemeProgress is a learner's pronunciation accuracy for one sound
// (rolling average from pronunciation scoring, Task 71).
type UserPhonemeProgress struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"uniqueIndex:idx_user_phoneme;not null" json:"userId"`
	PhonemeID     uint       `gorm:"uniqueIndex:idx_user_phoneme;index;not null" json:"phonemeId"`
	Accuracy      int        `gorm:"default:0;not null" json:"accuracy"`
	Attempts      int        `gorm:"default:0;not null" json:"attempts"`
	LastPracticed *time.Time `gorm:"column:last_practiced_at" json:"lastPracticedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
