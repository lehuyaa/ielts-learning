package models

import "time"

// AIConversationScenario is a user-owned conversation practice scenario
// (title, description, optional situation context, level, duration) used by
// the AI Conversation feature. Every user gets a default set seeded the
// first time they list their scenarios (see aiconversation.Service).
type AIConversationScenario struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"uniqueIndex:idx_scenario_user_slug;not null" json:"userId"`
	// Slug is unique per user, not globally, so every user can have their
	// own "travel", "job-interview", etc. without colliding with others.
	Slug             string          `gorm:"size:150;uniqueIndex:idx_scenario_user_slug;not null" json:"slug"`
	Title            string          `gorm:"size:150;not null" json:"title"`
	Description      string          `gorm:"type:text;not null" json:"description"`
	SituationContext string          `gorm:"type:text" json:"situationContext"`
	Level            DifficultyLevel `gorm:"type:varchar(20);default:'BEGINNER';not null" json:"level"`
	Duration         string          `gorm:"size:50" json:"duration"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
