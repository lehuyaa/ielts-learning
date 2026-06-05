package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Title        string         `gorm:"size:255;not null" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	Slug         string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	IsPublished  bool           `gorm:"default:false" json:"isPublished"`
	OrderIndex   int            `gorm:"default:0" json:"orderIndex"`
	BandMin      float64        `gorm:"type:decimal(3,1);default:5.0" json:"bandMin"`
	BandMax      float64        `gorm:"type:decimal(3,1);default:8.5" json:"bandMax"`
	TotalWords   int            `gorm:"default:0" json:"totalWords"`
	TotalLessons int            `gorm:"default:0" json:"totalLessons"`
	TotalTopics  int            `gorm:"default:0" json:"totalTopics"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	BandLevels []BandLevel `gorm:"foreignKey:CourseID" json:"bandLevels,omitempty"`
}

type BandLevel struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CourseID    uint           `gorm:"index;not null" json:"courseId"`
	BandScore   float64        `gorm:"type:decimal(3,1);not null" json:"bandScore"`
	MinScore    *float64       `gorm:"type:decimal(3,1)" json:"minScore"`
	MaxScore    *float64       `gorm:"type:decimal(3,1)" json:"maxScore"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	StatusLabel *string        `gorm:"size:50" json:"statusLabel"`
	OrderIndex  int            `gorm:"default:0" json:"orderIndex"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Course Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Topics []Topic `gorm:"foreignKey:BandLevelID" json:"topics,omitempty"`
}

type Topic struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	BandLevelID uint           `gorm:"uniqueIndex:idx_band_level_topic_slug;index;not null" json:"bandLevelId"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Slug        string         `gorm:"size:255;uniqueIndex:idx_band_level_topic_slug;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"size:100" json:"icon"`
	Emoji       string         `gorm:"size:20" json:"emoji"`
	Color       string         `gorm:"size:30" json:"color"`
	OrderIndex  int            `gorm:"default:0" json:"orderIndex"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	BandLevel BandLevel `gorm:"foreignKey:BandLevelID" json:"bandLevel,omitempty"`
	Lessons   []Lesson  `gorm:"foreignKey:TopicID" json:"lessons,omitempty"`
}
