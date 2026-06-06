package topic

import "ielts-learning/backend/internal/models"

type Response struct {
	Topic     TopicResponse     `json:"topic"`
	BandLevel BandLevelResponse `json:"bandLevel"`
	Summary   SummaryResponse   `json:"summary"`
	Lessons   []LessonResponse  `json:"lessons"`
}

type TopicResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Icon        string `json:"icon"`
	Emoji       string `json:"emoji"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type BandLevelResponse struct {
	ID        uint    `json:"id"`
	BandScore float64 `json:"bandScore"`
	Title     string  `json:"title"`
}

type SummaryResponse struct {
	ProgressPercentage int `json:"progressPercentage"`
	CompletedLessons   int `json:"completedLessons"`
	TotalLessons       int `json:"totalLessons"`
	TotalXP            int `json:"totalXP"`
}

type LessonResponse struct {
	ID                 uint                `json:"id"`
	Title              string              `json:"title"`
	Slug               string              `json:"slug"`
	Description        string              `json:"description"`
	WordCount          int                 `json:"wordCount"`
	EstimatedMinutes   int                 `json:"estimatedMinutes"`
	XPReward           int                 `json:"xpReward"`
	Status             models.LessonStatus `json:"status"`
	ProgressPercentage int                 `json:"progressPercentage"`
	LockedReason       *string             `json:"lockedReason"`
}
