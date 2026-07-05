package profile

import (
	"errors"
	"math"
	"strings"
	"time"

	"ielts-learning/backend/internal/models"
	xpmodule "ielts-learning/backend/internal/modules/xp"
)

var ErrUsernameAlreadyUsed = errors.New("username already used")

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) GetProfile(userID uint) (GetProfileResponse, error) {
	user, err := s.repository.FindUser(userID)
	if err != nil {
		return GetProfileResponse{}, err
	}

	lessonsDone, err := s.repository.CountLessonsDone(userID)
	if err != nil {
		return GetProfileResponse{}, err
	}
	wordsLearned, err := s.repository.CountWordsByStatuses(
		userID,
		models.VocabularyStatusLearning,
		models.VocabularyStatusReview,
		models.VocabularyStatusMastered,
	)
	if err != nil {
		return GetProfileResponse{}, err
	}
	masteredWords, err := s.repository.CountWordsByStatuses(userID, models.VocabularyStatusMastered)
	if err != nil {
		return GetProfileResponse{}, err
	}
	learningWords, err := s.repository.CountWordsByStatuses(
		userID,
		models.VocabularyStatusLearning,
		models.VocabularyStatusReview,
	)
	if err != nil {
		return GetProfileResponse{}, err
	}
	newWords, err := s.repository.CountWordsByStatuses(userID, models.VocabularyStatusNew)
	if err != nil {
		return GetProfileResponse{}, err
	}
	trackedWords, err := s.repository.CountTrackedWords(userID)
	if err != nil {
		return GetProfileResponse{}, err
	}
	passedQuizzes, err := s.repository.CountPassedQuizzes(userID)
	if err != nil {
		return GetProfileResponse{}, err
	}
	band7Lessons, err := s.repository.CountCompletedLessonsByMinimumBand(userID, 7.0)
	if err != nil {
		return GetProfileResponse{}, err
	}
	achievements, userAchievements, err := s.repository.FindAchievements(userID)
	if err != nil {
		return GetProfileResponse{}, err
	}
	recentXPEvents, err := s.repository.FindRecentXPEvents(userID, 5)
	if err != nil {
		return GetProfileResponse{}, err
	}
	dailyActivities, err := s.repository.FindDailyActivitiesSince(userID, time.Now().UTC().AddDate(0, 0, -83))
	if err != nil {
		return GetProfileResponse{}, err
	}

	stats := ProfileStatsResponse{
		WordsLearned:                wordsLearned,
		LessonsDone:                 lessonsDone,
		MasteredWords:               masteredWords,
		LearningWords:               learningWords,
		NewWords:                    newWords,
		VocabularyMasteryPercentage: percentage(masteredWords, trackedWords),
	}

	progress := achievementProgress{
		completedLessons: lessonsDone,
		learnedWords:     wordsLearned,
		passedQuizzes:    passedQuizzes,
		currentStreak:    user.CurrentStreak,
		band7Lessons:     band7Lessons,
	}

	return GetProfileResponse{
		User:            toUserProfileResponse(user),
		Stats:           stats,
		Achievements:    toAchievementItems(achievements, userAchievements, progress),
		ActivitySummary: toActivitySummary(dailyActivities, recentXPEvents),
	}, nil
}

func (s Service) UpdateProfile(userID uint, req UpdateProfileRequest) (UpdateProfileResponse, error) {
	user, err := s.repository.FindUser(userID)
	if err != nil {
		return UpdateProfileResponse{}, err
	}

	req = NormalizeUpdateProfileRequest(req)

	if req.Username != nil {
		exists, err := s.repository.UsernameExistsForOtherUser(userID, *req.Username)
		if err != nil {
			return UpdateProfileResponse{}, err
		}
		if exists {
			return UpdateProfileResponse{}, ErrUsernameAlreadyUsed
		}
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Username != nil {
		user.Username = req.Username
	}
	if req.TargetBand != nil {
		user.TargetBand = *req.TargetBand
	}
	if req.Timezone != nil {
		user.Timezone = *req.Timezone
	}
	if req.Locale != nil {
		user.Locale = *req.Locale
	}

	if err := s.repository.SaveUser(&user); err != nil {
		return UpdateProfileResponse{}, err
	}

	return toUserProfileResponse(user), nil
}

type achievementProgress struct {
	completedLessons int
	learnedWords     int
	passedQuizzes    int
	currentStreak    int
	band7Lessons     int
}

func toUserProfileResponse(user models.User) UserProfileResponse {
	level := maxInt(user.Level, 1)
	currentLevelFloor := xpmodule.CurrentLevelFloor(level)
	nextLevelXP := xpmodule.XPPerLevel
	currentLevelXP := user.TotalXP - currentLevelFloor
	if currentLevelXP < 0 {
		currentLevelXP = 0
	}
	if currentLevelXP > nextLevelXP {
		currentLevelXP = nextLevelXP
	}
	xpUntilNextLevel := xpmodule.NextLevelTotalXP(level) - user.TotalXP
	if xpUntilNextLevel < 0 {
		xpUntilNextLevel = 0
	}

	return UserProfileResponse{
		ID:               user.ID,
		Name:             user.Name,
		Username:         user.Username,
		Email:            user.Email,
		AvatarURL:        user.AvatarURL,
		Initials:         initialsFromName(user.Name, user.Email),
		MemberSince:      user.CreatedAt,
		CurrentBand:      user.CurrentBand,
		TargetBand:       user.TargetBand,
		TotalXP:          user.TotalXP,
		Level:            level,
		LevelTitle:       levelTitle(level),
		CurrentLevelXP:   currentLevelXP,
		NextLevelXP:      nextLevelXP,
		XPUntilNextLevel: xpUntilNextLevel,
		CurrentStreak:    user.CurrentStreak,
		LongestStreak:    user.LongestStreak,
		Timezone:         user.Timezone,
		Locale:           user.Locale,
	}
}

func toAchievementItems(
	achievements []models.Achievement,
	userAchievements []models.UserAchievement,
	progress achievementProgress,
) []ProfileAchievementItem {
	unlockedByAchievementID := make(map[uint]models.UserAchievement, len(userAchievements))
	for _, userAchievement := range userAchievements {
		unlockedByAchievementID[userAchievement.AchievementID] = userAchievement
	}

	items := make([]ProfileAchievementItem, 0, len(achievements))
	for _, achievement := range achievements {
		userAchievement, unlocked := unlockedByAchievementID[achievement.ID]
		var unlockedAt *time.Time
		if unlocked {
			value := userAchievement.UnlockedAt
			unlockedAt = &value
		}

		items = append(items, ProfileAchievementItem{
			ID:               achievement.ID,
			Code:             achievement.Code,
			Title:            achievement.Title,
			Description:      achievement.Description,
			Icon:             achievement.Icon,
			Unlocked:         unlocked,
			UnlockedAt:       unlockedAt,
			ProgressValue:    calculateAchievementProgress(achievement.RequirementType, progress),
			RequirementValue: achievement.RequirementValue,
		})
	}

	return items
}

func calculateAchievementProgress(requirementType string, progress achievementProgress) int {
	switch requirementType {
	case "completed_lessons":
		return progress.completedLessons
	case "learned_words":
		return progress.learnedWords
	case "passed_quizzes":
		return progress.passedQuizzes
	case "current_streak":
		return progress.currentStreak
	case "band_7_lessons":
		return progress.band7Lessons
	default:
		return 0
	}
}

func toActivitySummary(
	activities []models.DailyActivity,
	recentXPEvents []models.UserXPEvent,
) ActivitySummaryResponse {
	activeDays := 0
	wordsLearnedLast84 := 0
	for _, activity := range activities {
		wordsLearnedLast84 += activity.WordsLearned
		if activity.WordsLearned > 0 || activity.WordsReviewed > 0 || activity.QuizzesTaken > 0 || activity.LessonsDone > 0 || activity.XPEarned > 0 {
			activeDays++
		}
	}

	recentActivity := make([]ProfileRecentActivityItem, 0, len(recentXPEvents))
	for _, event := range recentXPEvents {
		title, description := summarizeXPEvent(event)
		recentActivity = append(recentActivity, ProfileRecentActivityItem{
			ID:          event.ID,
			Type:        event.SourceType,
			Title:       title,
			Description: description,
			XP:          event.XP,
			CreatedAt:   event.CreatedAt,
		})
	}

	averageWordsPerDay := math.Round((float64(wordsLearnedLast84)/84.0)*10) / 10

	return ActivitySummaryResponse{
		ActiveDays:         activeDays,
		WordsLearnedLast84: wordsLearnedLast84,
		AverageWordsPerDay: averageWordsPerDay,
		RecentActivity:     recentActivity,
	}
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
	for index, part := range parts {
		if part == "" {
			continue
		}
		parts[index] = strings.ToUpper(part[:1]) + part[1:]
	}

	return strings.Join(parts, " ")
}

func initialsFromName(name string, email string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) >= 2 {
		return strings.ToUpper(parts[0][:1] + parts[1][:1])
	}
	if len(parts) == 1 && len(parts[0]) > 0 {
		return strings.ToUpper(parts[0][:1])
	}
	if strings.TrimSpace(email) != "" {
		return strings.ToUpper(email[:1])
	}

	return "U"
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

func percentage(numerator int, denominator int) int {
	if denominator <= 0 {
		return 0
	}

	return int(math.Round((float64(numerator) / float64(denominator)) * 100))
}

func maxInt(value int, minimum int) int {
	if value < minimum {
		return minimum
	}

	return value
}
