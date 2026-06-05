package models

import "time"

type Achievement struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Code             string    `gorm:"size:100;uniqueIndex;not null" json:"code"`
	Title            string    `gorm:"size:255;not null" json:"title"`
	Description      string    `gorm:"type:text" json:"description"`
	Icon             string    `gorm:"size:100" json:"icon"`
	Category         string    `gorm:"size:50" json:"category"`
	RequirementType  string    `gorm:"size:50" json:"requirementType"`
	RequirementValue int       `gorm:"default:0" json:"requirementValue"`
	XPReward         int       `gorm:"default:0" json:"xpReward"`
	SortOrder        int       `gorm:"default:0" json:"sortOrder"`
	IsActive         bool      `gorm:"default:true" json:"isActive"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type UserAchievement struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex:idx_user_achievement;not null" json:"userId"`
	AchievementID uint      `gorm:"uniqueIndex:idx_user_achievement;not null" json:"achievementId"`
	ProgressValue int       `gorm:"default:0" json:"progressValue"`
	IsSeen        bool      `gorm:"default:false" json:"isSeen"`
	UnlockedAt    time.Time `json:"unlockedAt"`

	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Achievement Achievement `gorm:"foreignKey:AchievementID" json:"achievement,omitempty"`
}

type UserXPEvent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"userId"`
	SourceType  string    `gorm:"size:50;not null" json:"sourceType"`
	SourceID    *uint     `json:"sourceId"`
	XP          int       `gorm:"not null" json:"xp"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type DailyChallenge struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:100;uniqueIndex;not null" json:"code"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Metric      string    `gorm:"size:50;not null" json:"metric"`
	TargetValue int       `gorm:"not null" json:"targetValue"`
	XPReward    int       `gorm:"default:0" json:"xpReward"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserDailyChallenge struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `gorm:"uniqueIndex:idx_user_challenge_day;not null" json:"userId"`
	DailyChallengeID uint       `gorm:"uniqueIndex:idx_user_challenge_day;not null" json:"dailyChallengeId"`
	Date             time.Time  `gorm:"type:date;uniqueIndex:idx_user_challenge_day;not null" json:"date"`
	ProgressValue    int        `gorm:"default:0" json:"progressValue"`
	TargetValue      int        `gorm:"not null" json:"targetValue"`
	CompletedAt      *time.Time `json:"completedAt"`
	ClaimedAt        *time.Time `json:"claimedAt"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`

	User           User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DailyChallenge DailyChallenge `gorm:"foreignKey:DailyChallengeID" json:"dailyChallenge,omitempty"`
}

type Notification struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"index:idx_user_notification_read;index:idx_user_notification_created;not null" json:"userId"`
	Type      NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Title     string           `gorm:"size:255;not null" json:"title"`
	Body      string           `gorm:"type:text" json:"body"`
	ActionURL string           `gorm:"size:500" json:"actionUrl"`
	ReadAt    *time.Time       `gorm:"index:idx_user_notification_read" json:"readAt"`
	CreatedAt time.Time        `gorm:"index:idx_user_notification_created" json:"createdAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
