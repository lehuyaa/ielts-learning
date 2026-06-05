package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	Email                string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Name                 string         `gorm:"size:255" json:"name"`
	Username             *string        `gorm:"size:80;uniqueIndex" json:"username"`
	PasswordHash         string         `gorm:"size:255;not null" json:"-"`
	AvatarURL            string         `gorm:"size:500" json:"avatarUrl"`
	Role                 UserRole       `gorm:"type:varchar(20);default:'USER';not null" json:"role"`
	TargetBand           float64        `gorm:"type:decimal(3,1);default:7.0;index" json:"targetBand"`
	CurrentBand          *float64       `gorm:"type:decimal(3,1);index" json:"currentBand"`
	StartingBand         *float64       `gorm:"type:decimal(3,1)" json:"startingBand"`
	RecommendedBand      *float64       `gorm:"type:decimal(3,1)" json:"recommendedBand"`
	PlacementCompletedAt *time.Time     `json:"placementCompletedAt"`
	TotalXP              int            `gorm:"default:0" json:"totalXp"`
	Level                int            `gorm:"default:1" json:"level"`
	CurrentStreak        int            `gorm:"default:0" json:"currentStreak"`
	LongestStreak        int            `gorm:"default:0" json:"longestStreak"`
	LastActiveAt         *time.Time     `json:"lastActiveAt"`
	Timezone             string         `gorm:"size:80;default:'UTC';not null" json:"timezone"`
	Locale               string         `gorm:"size:20;default:'en';not null" json:"locale"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}
