package vocabulary

import (
	"time"

	"ielts-learning/backend/internal/models"
)

type ListQuery struct {
	Q          string
	Difficulty *models.DifficultyLevel
	TargetBand *float64
	Status     *models.VocabularyStatus
	Page       int
	Limit      int
}

type ListResponse struct {
	Items      []ListItemResponse `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type ListItemResponse struct {
	ID              uint                    `json:"id"`
	Word            string                  `json:"word"`
	Slug            string                  `json:"slug"`
	IPA             string                  `json:"ipa"`
	PartOfSpeech    string                  `json:"partOfSpeech"`
	MeaningVI       string                  `json:"meaningVi"`
	MeaningEN       string                  `json:"meaningEn"`
	ShortDefinition string                  `json:"shortDefinition"`
	Difficulty      models.DifficultyLevel  `json:"difficulty"`
	TargetBand      *float64                `json:"targetBand"`
	Topic           *TopicResponse          `json:"topic,omitempty"`
	Progress        ProgressSummaryResponse `json:"progress"`
	Status          models.VocabularyStatus `json:"status"`
	MasteryScore    int                     `json:"masteryScore"`
}

type DetailResponse struct {
	ID               uint                   `json:"id"`
	Word             string                 `json:"word"`
	Slug             string                 `json:"slug"`
	IPA              string                 `json:"ipa"`
	AudioURL         string                 `json:"audioUrl"`
	PartOfSpeech     string                 `json:"partOfSpeech"`
	MeaningVI        string                 `json:"meaningVi"`
	MeaningEN        string                 `json:"meaningEn"`
	PrimaryMeaning   string                 `json:"primaryMeaning"`
	SecondaryMeaning string                 `json:"secondaryMeaning"`
	ExampleSentences []string               `json:"exampleSentences"`
	Synonyms         []string               `json:"synonyms"`
	Antonyms         []string               `json:"antonyms"`
	Collocations     []string               `json:"collocations"`
	IELTSUsage       string                 `json:"ieltsUsage"`
	RelatedForms     []RelatedFormResponse  `json:"relatedForms"`
	Difficulty       models.DifficultyLevel `json:"difficulty"`
	TargetBand       *float64               `json:"targetBand"`
	Frequency        string                 `json:"frequency"`
	Rating           int                    `json:"rating"`
	Topic            *TopicResponse         `json:"topic,omitempty"`
	BandLevel        *BandLevelResponse     `json:"bandLevel,omitempty"`
	Progress         ProgressDetailResponse `json:"progress"`
	UserProgress     ProgressDetailResponse `json:"userProgress"`
	MasteryScore     int                    `json:"masteryScore"`
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

type RelatedFormResponse struct {
	Word         string `json:"word"`
	PartOfSpeech string `json:"partOfSpeech"`
}

type ProgressSummaryResponse struct {
	Status       models.VocabularyStatus `json:"status"`
	ReviewCount  int                     `json:"reviewCount"`
	NextReviewAt *time.Time              `json:"nextReviewAt"`
	MasteryScore int                     `json:"masteryScore"`
}

type ProgressDetailResponse struct {
	Status         models.VocabularyStatus `json:"status"`
	ReviewCount    int                     `json:"reviewCount"`
	CorrectCount   int                     `json:"correctCount"`
	WrongCount     int                     `json:"wrongCount"`
	NextReviewAt   *time.Time              `json:"nextReviewAt"`
	LastReviewedAt *time.Time              `json:"lastReviewedAt"`
	LearnedAt      *time.Time              `json:"learnedAt"`
	MasteryScore   int                     `json:"masteryScore"`
}
