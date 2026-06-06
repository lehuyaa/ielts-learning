package lesson

import (
	"time"

	"ielts-learning/backend/internal/models"
)

type DetailResponse struct {
	Lesson       LessonResponse       `json:"lesson"`
	Topic        TopicResponse        `json:"topic"`
	BandLevel    BandLevelResponse    `json:"bandLevel"`
	Progress     ProgressResponse     `json:"progress"`
	Vocabularies []VocabularyResponse `json:"vocabularies"`
}

type LessonResponse struct {
	ID                 uint                `json:"id"`
	Title              string              `json:"title"`
	Slug               string              `json:"slug"`
	Description        string              `json:"description"`
	RequiredScore      int                 `json:"requiredScore"`
	EstimatedMinutes   int                 `json:"estimatedMinutes"`
	XPReward           int                 `json:"xpReward"`
	WordCount          int                 `json:"wordCount"`
	OrderIndex         int                 `json:"orderIndex"`
	Status             models.LessonStatus `json:"status"`
	ProgressPercentage int                 `json:"progressPercentage"`
	LockedReason       *string             `json:"lockedReason"`
}

type TopicResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Icon  string `json:"icon"`
	Emoji string `json:"emoji"`
	Color string `json:"color"`
}

type BandLevelResponse struct {
	ID        uint    `json:"id"`
	BandScore float64 `json:"bandScore"`
	Title     string  `json:"title"`
}

type ProgressResponse struct {
	Status             models.LessonStatus `json:"status"`
	Score              *int                `json:"score"`
	BestScore          *int                `json:"bestScore"`
	BestXP             int                 `json:"bestXp"`
	WordsLearned       int                 `json:"wordsLearned"`
	TotalWords         int                 `json:"totalWords"`
	ProgressPercentage int                 `json:"progressPercentage"`
	StartedAt          *time.Time          `json:"startedAt"`
	CompletedAt        *time.Time          `json:"completedAt"`
	LastStudiedAt      *time.Time          `json:"lastStudiedAt"`
}

type VocabularyResponse struct {
	ID              uint                    `json:"id"`
	Word            string                  `json:"word"`
	Slug            string                  `json:"slug"`
	IPA             string                  `json:"ipa"`
	AudioURL        string                  `json:"audioUrl"`
	PartOfSpeech    string                  `json:"partOfSpeech"`
	MeaningVI       string                  `json:"meaningVi"`
	MeaningEN       string                  `json:"meaningEn"`
	ShortDefinition string                  `json:"shortDefinition"`
	ExampleSentence string                  `json:"exampleSentence"`
	Difficulty      models.DifficultyLevel  `json:"difficulty"`
	TargetBand      *float64                `json:"targetBand"`
	Status          models.VocabularyStatus `json:"status"`
	ReviewCount     int                     `json:"reviewCount"`
	CorrectCount    int                     `json:"correctCount"`
	WrongCount      int                     `json:"wrongCount"`
	Learned         bool                    `json:"learned"`
	LearnedAt       *time.Time              `json:"learnedAt"`
	LastReviewedAt  *time.Time              `json:"lastReviewedAt"`
	NextReviewAt    *time.Time              `json:"nextReviewAt"`
}

type StartResponse struct {
	LessonID      uint                `json:"lessonId"`
	Status        models.LessonStatus `json:"status"`
	StartedAt     *time.Time          `json:"startedAt"`
	LastStudiedAt *time.Time          `json:"lastStudiedAt"`
}
