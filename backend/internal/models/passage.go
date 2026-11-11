package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PassageGenre is the kind of text or recording (docs/foundation-spec.md §3.3).
type PassageGenre string

const (
	GenreDialogue  PassageGenre = "DIALOGUE"
	GenreMessage   PassageGenre = "MESSAGE"
	GenreEmail     PassageGenre = "EMAIL"
	GenreNotice    PassageGenre = "NOTICE"
	GenreMenu      PassageGenre = "MENU"
	GenreSchedule  PassageGenre = "SCHEDULE"
	GenreArticle   PassageGenre = "ARTICLE"
	GenreProfile   PassageGenre = "PROFILE"
	GenreMonologue PassageGenre = "MONOLOGUE"
)

// Passage is a text and/or recording used by LISTENING and READING
// activities. SegmentsJSON holds the sentences, each with either its own
// audio file or start/end times inside AudioURL, which makes sentence replay,
// dictation and shadowing possible.
type Passage struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Slug         string         `gorm:"size:150;uniqueIndex;not null" json:"slug"`
	CEFRLevel    string         `gorm:"column:cefr_level;size:5;index" json:"cefrLevel"`
	Genre        PassageGenre   `gorm:"type:varchar(30);not null" json:"genre"`
	Title        string         `gorm:"size:255;not null" json:"title"`
	BodyText     string         `gorm:"type:text;not null" json:"bodyText"`
	AudioURL     *string        `gorm:"size:500" json:"audioUrl"`
	SegmentsJSON datatypes.JSON `gorm:"type:json" json:"segments"`
	GlossaryJSON datatypes.JSON `gorm:"type:json" json:"glossary"`
	Source       string         `gorm:"size:255" json:"source"`
	WordCount    int            `gorm:"default:0;not null" json:"wordCount"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
