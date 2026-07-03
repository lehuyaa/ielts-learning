package profile

import "time"

type GetProfileResponse struct {
	User            UserProfileResponse      `json:"user"`
	Stats           ProfileStatsResponse     `json:"stats"`
	Achievements    []ProfileAchievementItem `json:"achievements"`
	ActivitySummary ActivitySummaryResponse  `json:"activitySummary"`
}

type UserProfileResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Username         *string   `json:"username"`
	Email            string    `json:"email"`
	AvatarURL        string    `json:"avatarUrl"`
	Initials         string    `json:"initials"`
	MemberSince      time.Time `json:"memberSince"`
	CurrentBand      *float64  `json:"currentBand"`
	TargetBand       float64   `json:"targetBand"`
	TotalXP          int       `json:"totalXp"`
	Level            int       `json:"level"`
	LevelTitle       string    `json:"levelTitle"`
	CurrentLevelXP   int       `json:"currentLevelXp"`
	NextLevelXP      int       `json:"nextLevelXp"`
	XPUntilNextLevel int       `json:"xpUntilNextLevel"`
	CurrentStreak    int       `json:"currentStreak"`
	LongestStreak    int       `json:"longestStreak"`
	Timezone         string    `json:"timezone"`
	Locale           string    `json:"locale"`
}

type ProfileStatsResponse struct {
	WordsLearned                int `json:"wordsLearned"`
	LessonsDone                 int `json:"lessonsDone"`
	MasteredWords               int `json:"masteredWords"`
	LearningWords               int `json:"learningWords"`
	NewWords                    int `json:"newWords"`
	VocabularyMasteryPercentage int `json:"vocabularyMasteryPercentage"`
}

type ProfileAchievementItem struct {
	ID               uint       `json:"id"`
	Code             string     `json:"code"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Icon             string     `json:"icon"`
	Unlocked         bool       `json:"unlocked"`
	UnlockedAt       *time.Time `json:"unlockedAt"`
	ProgressValue    int        `json:"progressValue"`
	RequirementValue int        `json:"requirementValue"`
}

type ActivitySummaryResponse struct {
	ActiveDays         int                         `json:"activeDays"`
	WordsLearnedLast84 int                         `json:"wordsLearnedLast84Days"`
	AverageWordsPerDay float64                     `json:"averageWordsPerDay"`
	RecentActivity     []ProfileRecentActivityItem `json:"recentActivity"`
}

type ProfileRecentActivityItem struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	XP          int       `json:"xp"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UpdateProfileRequest struct {
	Name       *string  `json:"name"`
	Username   *string  `json:"username"`
	TargetBand *float64 `json:"targetBand"`
	Timezone   *string  `json:"timezone"`
	Locale     *string  `json:"locale"`
}

type UpdateProfileResponse = UserProfileResponse
