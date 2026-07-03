package dashboard

import "time"

type SummaryResponse struct {
	User           UserSummaryResponse      `json:"user"`
	TodayProgress  TodayProgressResponse    `json:"todayProgress"`
	LearningStats  LearningStatsResponse    `json:"learningStats"`
	TargetBand     TargetBandResponse       `json:"targetBand"`
	XP             XPResponse               `json:"xp"`
	ReviewDue      ReviewDueResponse        `json:"reviewDue"`
	RecentActivity []RecentActivityResponse `json:"recentActivity"`
}

type UserSummaryResponse struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	AvatarURL     string   `json:"avatarUrl"`
	TargetBand    float64  `json:"targetBand"`
	CurrentBand   *float64 `json:"currentBand"`
	CurrentStreak int      `json:"currentStreak"`
}

type TodayProgressResponse struct {
	Date          string `json:"date"`
	WordsLearned  int    `json:"wordsLearned"`
	WordsReviewed int    `json:"wordsReviewed"`
	QuizzesTaken  int    `json:"quizzesTaken"`
	LessonsDone   int    `json:"lessonsDone"`
	XPEarned      int    `json:"xpEarned"`
}

type LearningStatsResponse struct {
	TotalWordsLearned int `json:"totalWordsLearned"`
	LessonsCompleted  int `json:"lessonsCompleted"`
	MasteryPercentage int `json:"masteryPercentage"`
}

type TargetBandResponse struct {
	StartingBand       *float64 `json:"startingBand"`
	CurrentBand        *float64 `json:"currentBand"`
	TargetBand         float64  `json:"targetBand"`
	ProgressPercentage int      `json:"progressPercentage"`
}

type XPResponse struct {
	TotalXP            int    `json:"totalXp"`
	Level              int    `json:"level"`
	LevelTitle         string `json:"levelTitle"`
	XPToNextLevel      int    `json:"xpToNextLevel"`
	ProgressPercentage int    `json:"progressPercentage"`
}

type ReviewDueResponse struct {
	Count int `json:"count"`
}

type RecentActivityResponse struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	XP          int       `json:"xp"`
	CreatedAt   time.Time `json:"createdAt"`
}
