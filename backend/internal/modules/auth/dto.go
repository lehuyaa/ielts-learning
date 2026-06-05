package auth

import (
	"time"

	"ielts-learning/backend/internal/models"
)

type RegisterRequest struct {
	Email      string   `json:"email" binding:"required,email"`
	Password   string   `json:"password" binding:"required,min=8"`
	Name       string   `json:"name" binding:"required,min=2,max=255"`
	Username   *string  `json:"username" binding:"omitempty,min=3,max=80"`
	TargetBand *float64 `json:"targetBand" binding:"omitempty,gte=0,lte=9"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string       `json:"accessToken"`
	TokenType   string       `json:"tokenType"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID                   uint            `json:"id"`
	Email                string          `json:"email"`
	Name                 string          `json:"name"`
	Username             *string         `json:"username"`
	AvatarURL            string          `json:"avatarUrl"`
	Role                 models.UserRole `json:"role"`
	TargetBand           float64         `json:"targetBand"`
	CurrentBand          *float64        `json:"currentBand"`
	StartingBand         *float64        `json:"startingBand"`
	RecommendedBand      *float64        `json:"recommendedBand"`
	PlacementCompletedAt *time.Time      `json:"placementCompletedAt"`
	TotalXP              int             `json:"totalXp"`
	Level                int             `json:"level"`
	LevelTitle           string          `json:"levelTitle"`
	CurrentStreak        int             `json:"currentStreak"`
	LongestStreak        int             `json:"longestStreak"`
	CreatedAt            time.Time       `json:"createdAt"`
}

func toUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:                   user.ID,
		Email:                user.Email,
		Name:                 user.Name,
		Username:             user.Username,
		AvatarURL:            user.AvatarURL,
		Role:                 user.Role,
		TargetBand:           user.TargetBand,
		CurrentBand:          user.CurrentBand,
		StartingBand:         user.StartingBand,
		RecommendedBand:      user.RecommendedBand,
		PlacementCompletedAt: user.PlacementCompletedAt,
		TotalXP:              user.TotalXP,
		Level:                user.Level,
		LevelTitle:           levelTitle(user.Level),
		CurrentStreak:        user.CurrentStreak,
		LongestStreak:        user.LongestStreak,
		CreatedAt:            user.CreatedAt,
	}
}

func levelTitle(level int) string {
	switch {
	case level >= 20:
		return "Master"
	case level >= 10:
		return "Expert"
	case level >= 5:
		return "Builder"
	default:
		return "Starter"
	}
}
