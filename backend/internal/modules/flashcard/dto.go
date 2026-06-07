package flashcard

import (
	"time"

	"ielts-learning/backend/internal/models"
)

type LessonFlashcardsResponse struct {
	Lesson   LessonResponse      `json:"lesson"`
	Progress SessionProgress     `json:"progress"`
	Cards    []FlashcardResponse `json:"cards"`
}

type LessonResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	TopicID    uint   `json:"topicId"`
	TopicTitle string `json:"topicTitle"`
	BandLabel  string `json:"bandLabel"`
}

type SessionProgress struct {
	Done      int `json:"done"`
	Remaining int `json:"remaining"`
	Total     int `json:"total"`
}

type DueReviewsResponse struct {
	Cards []FlashcardResponse `json:"cards"`
	Count int                 `json:"count"`
}

type FlashcardResponse struct {
	ID              uint                    `json:"id"`
	VocabularyID    uint                    `json:"vocabularyId"`
	Word            string                  `json:"word"`
	Slug            string                  `json:"slug"`
	IPA             string                  `json:"ipa"`
	AudioURL        string                  `json:"audioUrl"`
	PartOfSpeech    string                  `json:"partOfSpeech"`
	MeaningVI       string                  `json:"meaningVi"`
	MeaningEN       string                  `json:"meaningEn"`
	ExampleSentence string                  `json:"exampleSentence"`
	Synonyms        []string                `json:"synonyms"`
	Collocations    []string                `json:"collocations"`
	IELTSUsage      string                  `json:"ieltsUsage"`
	Difficulty      models.DifficultyLevel  `json:"difficulty"`
	Band            *float64                `json:"band"`
	TargetBand      *float64                `json:"targetBand"`
	TopicTitle      string                  `json:"topicTitle"`
	Status          models.VocabularyStatus `json:"status"`
	ReviewCount     int                     `json:"reviewCount"`
	NextReviewAt    *time.Time              `json:"nextReviewAt"`
}

type ReviewRequest struct {
	VocabularyID uint                   `json:"vocabularyId"`
	Rating       models.FlashcardRating `json:"rating"`
}

type ReviewResponse struct {
	VocabularyID   uint                    `json:"vocabularyId"`
	Status         models.VocabularyStatus `json:"status"`
	ReviewCount    int                     `json:"reviewCount"`
	CorrectCount   int                     `json:"correctCount"`
	WrongCount     int                     `json:"wrongCount"`
	LastRating     models.FlashcardRating  `json:"lastRating"`
	LastReviewedAt *time.Time              `json:"lastReviewedAt"`
	NextReviewAt   *time.Time              `json:"nextReviewAt"`
	MasteryScore   int                     `json:"masteryScore"`
	XPAwarded      int                     `json:"xpAwarded"`
	TotalXP        int                     `json:"totalXp"`
}
