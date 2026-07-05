package xp

import (
	"time"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

const XPPerLevel = 200

type EventType string

const (
	EventFlashcardReview     EventType = "FLASHCARD_REVIEW"
	EventQuizCorrect         EventType = "QUIZ_CORRECT"
	EventLessonCompleted     EventType = "LESSON_COMPLETED"
	EventAchievementUnlocked EventType = "ACHIEVEMENT_UNLOCKED"
)

type AwardInput struct {
	UserID           uint
	SourceType       EventType
	SourceID         *uint
	XP               int
	Description      string
	AwardedAt        time.Time
	PreventDuplicate bool
	TouchLastActive  bool
}

type AwardResult struct {
	User         models.User
	EventCreated bool
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) AwardXP(tx *gorm.DB, input AwardInput) (AwardResult, error) {
	user, err := s.repository.FindUserForUpdate(tx, input.UserID)
	if err != nil {
		return AwardResult{}, err
	}

	if input.XP <= 0 {
		return AwardResult{User: user, EventCreated: false}, nil
	}

	awardedAt := normalizeTime(input.AwardedAt)

	if input.PreventDuplicate {
		exists, err := s.repository.XPEventExists(tx, input.UserID, input.SourceType, input.SourceID)
		if err != nil {
			return AwardResult{}, err
		}
		if exists {
			return AwardResult{User: user, EventCreated: false}, nil
		}
	}

	event := models.UserXPEvent{
		UserID:      input.UserID,
		SourceType:  string(input.SourceType),
		SourceID:    input.SourceID,
		XP:          input.XP,
		Description: input.Description,
		CreatedAt:   awardedAt,
	}
	if err := s.repository.CreateXPEvent(tx, event); err != nil {
		return AwardResult{}, err
	}

	user.TotalXP += input.XP
	user.Level = LevelForTotalXP(user.TotalXP)
	if input.TouchLastActive {
		user.LastActiveAt = &awardedAt
	}

	if err := s.repository.SaveUser(tx, &user); err != nil {
		return AwardResult{}, err
	}

	return AwardResult{User: user, EventCreated: true}, nil
}

func (s Service) RecalculateLevel(tx *gorm.DB, userID uint) (models.User, error) {
	user, err := s.repository.FindUserForUpdate(tx, userID)
	if err != nil {
		return models.User{}, err
	}

	totalXP, err := s.repository.SumUserXP(tx, userID)
	if err != nil {
		return models.User{}, err
	}

	user.TotalXP = totalXP
	user.Level = LevelForTotalXP(totalXP)

	if err := s.repository.SaveUser(tx, &user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func LevelForTotalXP(totalXP int) int {
	if totalXP < 0 {
		totalXP = 0
	}

	return totalXP/XPPerLevel + 1
}

func CurrentLevelFloor(level int) int {
	level = maxInt(level, 1)

	return (level - 1) * XPPerLevel
}

func NextLevelTotalXP(level int) int {
	level = maxInt(level, 1)

	return level * XPPerLevel
}

func maxInt(value int, minimum int) int {
	if value < minimum {
		return minimum
	}

	return value
}
