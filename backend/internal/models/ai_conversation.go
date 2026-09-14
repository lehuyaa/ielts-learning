package models

import "time"

// AIConversationScenario is a user-owned conversation practice scenario
// (title, description, optional situation context, level, duration) used by
// the AI Conversation feature.
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

// AIConversationMessage is one turn (user or assistant) in a scenario's chat
// history, persisted so a session can be resumed later.
type AIConversationMessage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ScenarioID uint      `gorm:"index;not null" json:"scenarioId"`
	Role       string    `gorm:"size:20;not null" json:"role"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time `json:"createdAt"`

	Scenario AIConversationScenario `gorm:"foreignKey:ScenarioID" json:"-"`
}
