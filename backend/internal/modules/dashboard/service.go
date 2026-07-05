package dashboard

import (
	"errors"
	"math"
	"strings"
	"time"

	"ielts-learning/backend/internal/models"
	xpmodule "ielts-learning/backend/internal/modules/xp"
)

var ErrUserNotFound = errors.New("user not found")

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) GetSummary(userID uint) (SummaryResponse, error) {
	user, err := s.repository.FindUser(userID)
	if err != nil {
		return SummaryResponse{}, err
	}

	now := time.Now().UTC()
	userLocation := loadLocation(user.Timezone)
	userNow := now.In(userLocation)
	todayDate := userNow.Format("2006-01-02")

	todayActivity, foundTodayActivity, err := s.repository.FindDailyActivity(userID, todayDate)
	if err != nil {
		return SummaryResponse{}, err
	}

	reviewDueCount, err := s.repository.CountReviewDue(userID, now)
	if err != nil {
		return SummaryResponse{}, err
	}

	totalWordsLearned, err := s.repository.CountTotalWordsLearned(userID)
	if err != nil {
		return SummaryResponse{}, err
	}

	masteredWords, err := s.repository.CountMasteredWords(userID)
	if err != nil {
		return SummaryResponse{}, err
	}

	lessonsCompleted, err := s.repository.CountCompletedLessons(userID)
	if err != nil {
		return SummaryResponse{}, err
	}

	recentXPEvents, err := s.repository.FindRecentXPEvents(userID, 5)
	if err != nil {
		return SummaryResponse{}, err
	}

	return SummaryResponse{
		User: UserSummaryResponse{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			AvatarURL:     user.AvatarURL,
			TargetBand:    user.TargetBand,
			CurrentBand:   user.CurrentBand,
			CurrentStreak: user.CurrentStreak,
		},
		TodayProgress:  toTodayProgress(todayDate, todayActivity, foundTodayActivity),
		LearningStats:  toLearningStats(totalWordsLearned, lessonsCompleted, masteredWords),
		TargetBand:     toTargetBand(user),
		XP:             toXP(user),
		ReviewDue:      ReviewDueResponse{Count: reviewDueCount},
		RecentActivity: toRecentActivity(recentXPEvents),
	}, nil
}

func toTodayProgress(date string, activity models.DailyActivity, found bool) TodayProgressResponse {
	if !found {
		return TodayProgressResponse{
			Date: date,
		}
	}

	return TodayProgressResponse{
		Date:          date,
		WordsLearned:  activity.WordsLearned,
		WordsReviewed: activity.WordsReviewed,
		QuizzesTaken:  activity.QuizzesTaken,
		LessonsDone:   activity.LessonsDone,
		XPEarned:      activity.XPEarned,
	}
}

func toLearningStats(totalWordsLearned int, lessonsCompleted int, masteredWords int) LearningStatsResponse {
	return LearningStatsResponse{
		TotalWordsLearned: totalWordsLearned,
		LessonsCompleted:  lessonsCompleted,
		MasteryPercentage: percentage(masteredWords, totalWordsLearned),
	}
}

func toTargetBand(user models.User) TargetBandResponse {
	startingBand := user.StartingBand
	currentBand := user.CurrentBand

	return TargetBandResponse{
		StartingBand:       startingBand,
		CurrentBand:        currentBand,
		TargetBand:         user.TargetBand,
		ProgressPercentage: calculateTargetBandProgress(user),
	}
}

func toXP(user models.User) XPResponse {
	level := maxInt(user.Level, 1)
	currentLevelFloor := xpmodule.CurrentLevelFloor(level)
	nextLevelTotalXP := xpmodule.NextLevelTotalXP(level)
	xpToNextLevel := nextLevelTotalXP - user.TotalXP
	if xpToNextLevel < 0 {
		xpToNextLevel = 0
	}

	progressWithinLevel := user.TotalXP - currentLevelFloor
	levelSpan := nextLevelTotalXP - currentLevelFloor

	return XPResponse{
		TotalXP:            user.TotalXP,
		Level:              level,
		LevelTitle:         levelTitle(level),
		XPToNextLevel:      xpToNextLevel,
		ProgressPercentage: percentage(progressWithinLevel, levelSpan),
	}
}

func toRecentActivity(events []models.UserXPEvent) []RecentActivityResponse {
	activities := make([]RecentActivityResponse, 0, len(events))
	for _, event := range events {
		title, description := summarizeXPEvent(event)
		activities = append(activities, RecentActivityResponse{
			ID:          event.ID,
			Type:        event.SourceType,
			Title:       title,
			Description: description,
			XP:          event.XP,
			CreatedAt:   event.CreatedAt,
		})
	}

	return activities
}

func summarizeXPEvent(event models.UserXPEvent) (string, string) {
	if strings.TrimSpace(event.Description) != "" {
		return titleFromSourceType(event.SourceType), event.Description
	}

	return titleFromSourceType(event.SourceType), "XP earned from learning activity"
}

func titleFromSourceType(sourceType string) string {
	switch sourceType {
	case string(xpmodule.EventLessonCompleted), "LESSON_COMPLETION":
		return "Lesson Completion"
	case string(xpmodule.EventFlashcardReview):
		return "Flashcard Review"
	case string(xpmodule.EventQuizCorrect):
		return "Quiz Correct"
	case "QUIZ_SUCCESS", "QUIZ_COMPLETION":
		return "Quiz Success"
	case string(xpmodule.EventAchievementUnlocked):
		return "Achievement Unlocked"
	default:
		return humanizeSourceType(sourceType)
	}
}

func humanizeSourceType(sourceType string) string {
	if sourceType == "" {
		return "Learning Activity"
	}

	parts := strings.Split(strings.ToLower(sourceType), "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}

	return strings.Join(parts, " ")
}

func calculateTargetBandProgress(user models.User) int {
	if user.CurrentBand == nil {
		return 0
	}

	currentBand := *user.CurrentBand
	if currentBand <= 0 || user.TargetBand <= 0 {
		return 0
	}

	if user.StartingBand != nil && user.TargetBand > *user.StartingBand {
		return clampPercentage(
			int(math.Round(((currentBand - *user.StartingBand) / (user.TargetBand - *user.StartingBand)) * 100)),
		)
	}

	return clampPercentage(int(math.Round((currentBand / user.TargetBand) * 100)))
}

func loadLocation(name string) *time.Location {
	if strings.TrimSpace(name) == "" {
		return time.UTC
	}

	location, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}

	return location
}

func percentage(part int, total int) int {
	if total <= 0 || part <= 0 {
		return 0
	}

	return clampPercentage(int(math.Round((float64(part) / float64(total)) * 100)))
}

func clampPercentage(value int) int {
	switch {
	case value < 0:
		return 0
	case value > 100:
		return 100
	default:
		return value
	}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
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
