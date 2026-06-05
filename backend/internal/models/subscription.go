package models

import (
	"time"

	"gorm.io/datatypes"
)

type SubscriptionPlan struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Code            string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	PriceCents      int            `gorm:"default:0" json:"priceCents"`
	Currency        string         `gorm:"size:10;default:'USD';not null" json:"currency"`
	BillingInterval string         `gorm:"size:20;default:'MONTH';not null" json:"billingInterval"`
	MaxWords        *int           `json:"maxWords"`
	MaxTopics       *int           `json:"maxTopics"`
	MaxUsers        int            `gorm:"default:1" json:"maxUsers"`
	FeaturesJSON    datatypes.JSON `gorm:"type:json" json:"features,omitempty"`
	IsPopular       bool           `gorm:"default:false" json:"isPopular"`
	IsActive        bool           `gorm:"default:true" json:"isActive"`
	SortOrder       int            `gorm:"default:0" json:"sortOrder"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

type UserSubscription struct {
	ID                  uint               `gorm:"primaryKey" json:"id"`
	UserID              uint               `gorm:"index;not null" json:"userId"`
	SubscriptionPlanID  uint               `gorm:"index;not null" json:"subscriptionPlanId"`
	Status              SubscriptionStatus `gorm:"type:varchar(30);default:'ACTIVE';not null" json:"status"`
	StartedAt           time.Time          `gorm:"not null" json:"startedAt"`
	TrialEndsAt         *time.Time         `json:"trialEndsAt"`
	CurrentPeriodEndsAt *time.Time         `json:"currentPeriodEndsAt"`
	CanceledAt          *time.Time         `json:"canceledAt"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`

	User             User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SubscriptionPlan SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID" json:"subscriptionPlan,omitempty"`
}

type AIExampleRequest struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"userId"`
	VocabularyID  *uint     `gorm:"index" json:"vocabularyId"`
	TopicID       *uint     `gorm:"index" json:"topicId"`
	Prompt        string    `gorm:"type:text;not null" json:"prompt"`
	GeneratedText string    `gorm:"type:text" json:"generatedText"`
	Status        string    `gorm:"size:30;default:'PENDING';not null" json:"status"`
	ErrorMessage  string    `gorm:"type:text" json:"errorMessage"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Vocabulary *Vocabulary `gorm:"foreignKey:VocabularyID" json:"vocabulary,omitempty"`
	Topic      *Topic      `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
}
