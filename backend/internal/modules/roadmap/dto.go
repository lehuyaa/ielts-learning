package roadmap

import "ielts-learning/backend/internal/models"

type Response struct {
	Course     CourseResponse      `json:"course"`
	Summary    SummaryResponse     `json:"summary"`
	BandLevels []BandLevelResponse `json:"bandLevels"`
}

type CourseResponse struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	Slug         string  `json:"slug"`
	BandMin      float64 `json:"bandMin"`
	BandMax      float64 `json:"bandMax"`
	TotalWords   int     `json:"totalWords"`
	TotalLessons int     `json:"totalLessons"`
	TotalTopics  int     `json:"totalTopics"`
}

type SummaryResponse struct {
	TopicsCompleted int      `json:"topicsCompleted"`
	TotalTopics     int      `json:"totalTopics"`
	CurrentBand     *float64 `json:"currentBand"`
	WordsMastered   int      `json:"wordsMastered"`
	CurrentStreak   int      `json:"currentStreak"`
}

type BandLevelResponse struct {
	ID                 uint                `json:"id"`
	BandScore          float64             `json:"bandScore"`
	Title              string              `json:"title"`
	Description        string              `json:"description"`
	Status             models.LessonStatus `json:"status"`
	ProgressPercentage int                 `json:"progressPercentage"`
	TopicsCompleted    int                 `json:"topicsCompleted"`
	TotalTopics        int                 `json:"totalTopics"`
	Topics             []TopicResponse     `json:"topics"`
}

type TopicResponse struct {
	ID                 uint                `json:"id"`
	Title              string              `json:"title"`
	Slug               string              `json:"slug"`
	Emoji              string              `json:"emoji"`
	Color              string              `json:"color"`
	Status             models.LessonStatus `json:"status"`
	LessonsCompleted   int                 `json:"lessonsCompleted"`
	TotalLessons       int                 `json:"totalLessons"`
	ProgressPercentage int                 `json:"progressPercentage"`
	Lessons            []LessonResponse    `json:"lessons"`
}

type LessonResponse struct {
	ID               uint                `json:"id"`
	Title            string              `json:"title"`
	Slug             string              `json:"slug"`
	Status           models.LessonStatus `json:"status"`
	RequiredScore    int                 `json:"requiredScore"`
	EstimatedMinutes int                 `json:"estimatedMinutes"`
	XPReward         int                 `json:"xpReward"`
	Score            *int                `json:"score"`
	BestScore        *int                `json:"bestScore"`
}
