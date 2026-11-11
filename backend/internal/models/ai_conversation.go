package models

import (
	"time"

	"gorm.io/datatypes"
)

// AIConversationScenario is a conversation practice scenario: either owned by
// a user (custom) or a system scenario shared by everyone (UserID nil,
// IsSystem true), with the role the AI plays and the learner's task
// (docs/foundation-spec.md §3.1).
type AIConversationScenario struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// UserID is nil for system scenarios.
	UserID *uint `gorm:"uniqueIndex:idx_scenario_user_slug" json:"userId"`
	// Slug is unique per user (and among system scenarios), not globally, so
	// every user can have their own "travel", "job-interview", etc.
	Slug             string          `gorm:"size:150;uniqueIndex:idx_scenario_user_slug;not null" json:"slug"`
	Title            string          `gorm:"size:150;not null" json:"title"`
	Description      string          `gorm:"type:text;not null" json:"description"`
	SituationContext string          `gorm:"type:text" json:"situationContext"`
	Level            DifficultyLevel `gorm:"type:varchar(20);default:'BEGINNER';not null" json:"level"`
	Duration         string          `gorm:"size:50" json:"duration"`

	IsSystem           bool           `gorm:"default:false;not null;index" json:"isSystem"`
	CEFRLevel          *string        `gorm:"column:cefr_level;size:5" json:"cefrLevel"`
	Role               string         `gorm:"type:text" json:"role"`
	Task               string         `gorm:"type:text" json:"task"`
	TargetPatternsJSON datatypes.JSON `gorm:"type:json" json:"targetPatterns"`
	MaxTurns           int            `gorm:"default:8;not null" json:"maxTurns"`
	// OpeningLine is the AI's first message of a session.
	OpeningLine string `gorm:"type:text" json:"openingLine"`
	// PassageID is the model dialogue for this situation (shared with listening).
	PassageID  *uint  `gorm:"index" json:"passageId"`
	OrderIndex int    `gorm:"default:0" json:"orderIndex"`
	Group      string `gorm:"column:group_label;size:100" json:"group"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// AIConversationMessage is one turn (user or assistant). Messages of a
// session carry SessionID; older scenario-level chats have none.
type AIConversationMessage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ScenarioID uint      `gorm:"index;not null" json:"scenarioId"`
	SessionID  *uint     `gorm:"index" json:"sessionId"`
	Role       string    `gorm:"size:20;not null" json:"role"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time `json:"createdAt"`

	Scenario AIConversationScenario `gorm:"foreignKey:ScenarioID" json:"-"`
}

type SessionStatus string

const (
	SessionInProgress SessionStatus = "IN_PROGRESS"
	SessionFinished   SessionStatus = "FINISHED"
	SessionAbandoned  SessionStatus = "ABANDONED"
)

// AIConversationSession is one role-play attempt with its feedback report
// (docs/foundation-spec.md §3.4, §6.2).
type AIConversationSession struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"userId"`
	ScenarioID uint           `gorm:"index;not null" json:"scenarioId"`
	ActivityID *uint          `gorm:"index" json:"activityId"`
	Status     SessionStatus  `gorm:"type:varchar(20);default:'IN_PROGRESS';not null" json:"status"`
	TurnCount  int            `gorm:"default:0;not null" json:"turnCount"`
	HintCount  int            `gorm:"default:0;not null" json:"hintCount"`
	Score      *int           `json:"score"`
	TaskDone   bool           `gorm:"column:task_completed;default:false;not null" json:"taskCompleted"`
	ReportJSON datatypes.JSON `gorm:"type:json" json:"report"`
	StartedAt  time.Time      `json:"startedAt"`
	FinishedAt *time.Time     `json:"finishedAt"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}
