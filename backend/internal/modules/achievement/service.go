package achievement

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
	notificationmodule "ielts-learning/backend/internal/modules/notification"
	xpmodule "ielts-learning/backend/internal/modules/xp"
)

const (
	CodeFirstLesson      = "FIRST_LESSON"
	CodeStreak7          = "STREAK_7"
	CodeWordCollector100 = "WORD_COLLECTOR_100"
	CodeEducationMaster  = "EDUCATION_MASTER"
)

var ErrAchievementNotFound = errors.New("achievement not found")

type UnlockResult struct {
	Code      string
	Unlocked  bool
	CreatedAt time.Time
}

type Service struct {
	repository          Repository
	xpService           xpmodule.Service
	notificationService notificationmodule.Service
}

func NewService(repository Repository, xpService xpmodule.Service, notificationService notificationmodule.Service) Service {
	return Service{
		repository:          repository,
		xpService:           xpService,
		notificationService: notificationService,
	}
}

func (s Service) CheckAndUnlockAchievements(userID uint) ([]UnlockResult, error) {
	var results []UnlockResult
	err := s.repository.db.Transaction(func(tx *gorm.DB) error {
		var err error
		results, err = s.CheckAndUnlockAchievementsTx(tx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s Service) CheckAndUnlockAchievementsTx(tx *gorm.DB, userID uint) ([]UnlockResult, error) {
	achievements, err := s.repository.FindActiveAchievementsByCodes(tx, []string{
		CodeFirstLesson,
		CodeStreak7,
		CodeWordCollector100,
		CodeEducationMaster,
	})
	if err != nil {
		return nil, err
	}

	user, err := s.repository.FindUser(tx, userID)
	if err != nil {
		return nil, err
	}

	completedLessons, err := s.repository.CountCompletedLessons(tx, userID)
	if err != nil {
		return nil, err
	}

	learnedWords, err := s.repository.CountLearnedWords(tx, userID)
	if err != nil {
		return nil, err
	}

	educationLessonCount, err := s.repository.CountLessonsForTopicSlug(tx, "education")
	if err != nil {
		return nil, err
	}

	completedEducationLessons, err := s.repository.CountCompletedLessonsForTopicSlug(tx, userID, "education")
	if err != nil {
		return nil, err
	}

	progressByCode := map[string]int{
		CodeFirstLesson:      completedLessons,
		CodeStreak7:          user.CurrentStreak,
		CodeWordCollector100: learnedWords,
		CodeEducationMaster:  completedEducationLessons,
	}

	results := make([]UnlockResult, 0, len(achievements))
	for _, achievement := range achievements {
		if !eligibleForUnlock(achievement.Code, completedLessons, user.CurrentStreak, learnedWords, educationLessonCount, completedEducationLessons) {
			continue
		}

		unlocked, unlockedAt, err := s.unlockAchievementTx(tx, userID, achievement, progressByCode[achievement.Code])
		if err != nil {
			return nil, err
		}
		if unlocked {
			results = append(results, UnlockResult{
				Code:      achievement.Code,
				Unlocked:  true,
				CreatedAt: unlockedAt,
			})
		}
	}

	return results, nil
}

func (s Service) UnlockAchievement(userID uint, achievementCode string) (bool, error) {
	var unlocked bool
	err := s.repository.db.Transaction(func(tx *gorm.DB) error {
		achievement, err := s.repository.FindAchievementByCode(tx, achievementCode)
		if err != nil {
			return err
		}

		progressValue, err := s.progressValueForAchievement(tx, userID, achievement.Code)
		if err != nil {
			return err
		}

		unlocked, _, err = s.unlockAchievementTx(tx, userID, achievement, progressValue)
		return err
	})
	if err != nil {
		return false, err
	}

	return unlocked, nil
}

func (s Service) HasAchievement(userID uint, achievementCode string) (bool, error) {
	var hasAchievement bool
	err := s.repository.db.Transaction(func(tx *gorm.DB) error {
		achievement, err := s.repository.FindAchievementByCode(tx, achievementCode)
		if err != nil {
			return err
		}

		hasAchievement, err = s.repository.HasAchievement(tx, userID, achievement.ID)
		return err
	})
	if err != nil {
		return false, err
	}

	return hasAchievement, nil
}

func (s Service) unlockAchievementTx(tx *gorm.DB, userID uint, achievement models.Achievement, progressValue int) (bool, time.Time, error) {
	unlockedAt := nowUTC()
	created, err := s.repository.CreateUserAchievement(tx, models.UserAchievement{
		UserID:        userID,
		AchievementID: achievement.ID,
		ProgressValue: progressValue,
		IsSeen:        false,
		UnlockedAt:    unlockedAt,
	})
	if err != nil {
		return false, time.Time{}, err
	}
	if !created {
		return false, time.Time{}, nil
	}

	if achievement.XPReward > 0 {
		sourceID := achievement.ID
		if _, err := s.xpService.AwardXP(tx, xpmodule.AwardInput{
			UserID:           userID,
			SourceType:       xpmodule.EventAchievementUnlocked,
			SourceID:         &sourceID,
			XP:               achievement.XPReward,
			Description:      fmt.Sprintf("Unlocked achievement: %s", achievement.Title),
			AwardedAt:        unlockedAt,
			PreventDuplicate: true,
			TouchLastActive:  false,
		}); err != nil {
			return false, time.Time{}, fmt.Errorf("award achievement xp: %w", err)
		}
	}

	if err := s.notificationService.CreateNotificationTx(tx, notificationmodule.CreateInput{
		UserID:  userID,
		Type:    models.NotificationAchievement,
		Title:   achievement.Title,
		Message: achievement.Description,
		Metadata: map[string]string{
			"actionUrl": "/profile",
		},
	}); err != nil {
		return false, time.Time{}, err
	}

	return true, unlockedAt, nil
}

func (s Service) progressValueForAchievement(tx *gorm.DB, userID uint, achievementCode string) (int, error) {
	user, err := s.repository.FindUser(tx, userID)
	if err != nil {
		return 0, err
	}

	switch achievementCode {
	case CodeFirstLesson:
		return s.repository.CountCompletedLessons(tx, userID)
	case CodeStreak7:
		return user.CurrentStreak, nil
	case CodeWordCollector100:
		return s.repository.CountLearnedWords(tx, userID)
	case CodeEducationMaster:
		return s.repository.CountCompletedLessonsForTopicSlug(tx, userID, "education")
	default:
		return 0, ErrAchievementNotFound
	}
}

func eligibleForUnlock(code string, completedLessons int, currentStreak int, learnedWords int, educationLessonCount int, completedEducationLessons int) bool {
	switch code {
	case CodeFirstLesson:
		return completedLessons >= 1
	case CodeStreak7:
		return currentStreak >= 7
	case CodeWordCollector100:
		return learnedWords >= 100
	case CodeEducationMaster:
		return educationLessonCount > 0 && completedEducationLessons >= educationLessonCount
	default:
		return false
	}
}
