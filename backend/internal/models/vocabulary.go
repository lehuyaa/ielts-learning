package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Vocabulary struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	Word             string          `gorm:"size:255;index;not null" json:"word"`
	Slug             string          `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	IPA              string          `gorm:"size:255" json:"ipa"`
	AudioURL         string          `gorm:"size:500" json:"audioUrl"`
	PartOfSpeech     string          `gorm:"size:100" json:"partOfSpeech"`
	MeaningVI        string          `gorm:"type:text;not null" json:"meaningVi"`
	MeaningEN        string          `gorm:"type:text" json:"meaningEn"`
	ShortDefinition  string          `gorm:"type:text" json:"shortDefinition"`
	ExampleSentence  string          `gorm:"type:text" json:"exampleSentence"`
	ExampleMeaningVI string          `gorm:"type:text" json:"exampleMeaningVi"`
	ExampleSource    string          `gorm:"size:255" json:"exampleSource"`
	SynonymsJSON     datatypes.JSON  `gorm:"type:json" json:"synonyms,omitempty"`
	AntonymsJSON     datatypes.JSON  `gorm:"type:json" json:"antonyms,omitempty"`
	CollocationsJSON datatypes.JSON  `gorm:"type:json" json:"collocations,omitempty"`
	IELTSUsage       string          `gorm:"type:text" json:"ieltsUsage"`
	Difficulty       DifficultyLevel `gorm:"type:varchar(30);default:'INTERMEDIATE';index" json:"difficulty"`
	TargetBand       *float64        `gorm:"type:decimal(3,1);index" json:"targetBand"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}
