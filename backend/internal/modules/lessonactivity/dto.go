package lessonactivity

import (
	"encoding/json"
	"time"

	"ielts-learning/backend/internal/models"
)

// ActivityStatusNotStarted is returned for activities without progress.
const ActivityStatusNotStarted = "NOT_STARTED"

type LessonActivitiesResponse struct {
	Lesson     LessonSummary      `json:"lesson"`
	Progress   ActivitiesProgress `json:"progress"`
	Activities []ActivityItem     `json:"activities"`
}

type LessonSummary struct {
	ID     uint                `json:"id"`
	Title  string              `json:"title"`
	Slug   string              `json:"slug"`
	Status models.LessonStatus `json:"status"`
}

type ActivitiesProgress struct {
	PassedRequired int `json:"passedRequired"`
	TotalRequired  int `json:"totalRequired"`
	// NextActivityID is the first unfinished activity, for resuming.
	NextActivityID *uint `json:"nextActivityId"`
}

type ActivityItem struct {
	ID             uint                `json:"id"`
	Slug           string              `json:"slug"`
	Type           models.ActivityType `json:"type"`
	Title          string              `json:"title"`
	InstructionsVI string              `json:"instructionsVi"`
	OrderIndex     int                 `json:"orderIndex"`
	IsRequired     bool                `json:"isRequired"`
	PassScore      int                 `json:"passScore"`
	XPReward       int                 `json:"xpReward"`
	ExerciseCount  int                 `json:"exerciseCount"`
	Status         string              `json:"status"`
	Score          *int                `json:"score"`
	BestScore      *int                `json:"bestScore"`
	Attempts       int                 `json:"attempts"`
	IsLocked       bool                `json:"isLocked"`
}

type ActivityDetailResponse struct {
	Activity  ActivityDetail     `json:"activity"`
	Exercises []ExerciseResponse `json:"exercises"`
	// Vocabularies are the word cards of a VOCABULARY activity, in content order.
	Vocabularies []VocabularyCard `json:"vocabularies"`
	// Answers already given in the current (unfinished) attempt, for resuming.
	Answers []AnswerResult `json:"answers"`
}

type ActivityDetail struct {
	ActivityItem
	LessonID       uint            `json:"lessonId"`
	Content        json.RawMessage `json:"content" swaggertype:"object"`
	GrammarPointID *uint           `json:"grammarPointId"`
	PassageID      *uint           `json:"passageId"`
	ScenarioID     *uint           `json:"scenarioId"`
	CurrentAttempt int             `json:"currentAttempt"`
}

// ExerciseResponse never contains the answer key.
type ExerciseResponse struct {
	ID         uint                `json:"id"`
	Type       models.ExerciseType `json:"type"`
	Prompt     string              `json:"prompt"`
	Payload    json.RawMessage     `json:"payload" swaggertype:"object"`
	Points     int                 `json:"points"`
	OrderIndex int                 `json:"orderIndex"`
}

type SubmitAnswerRequest struct {
	ExerciseID uint            `json:"exerciseId"`
	Answer     json.RawMessage `json:"answer" swaggertype:"object"`
}

type AnswerResult struct {
	ExerciseID    uint            `json:"exerciseId"`
	Attempt       int             `json:"attempt"`
	IsCorrect     bool            `json:"isCorrect"`
	ScoreRatio    float64         `json:"scoreRatio"`
	Answer        json.RawMessage `json:"answer" swaggertype:"object"`
	CorrectAnswer any             `json:"correctAnswer"`
	ExplanationVI string          `json:"explanationVi"`
	AnsweredAt    time.Time       `json:"answeredAt"`
}

type VocabularyCard struct {
	ID               uint   `json:"id"`
	Word             string `json:"word"`
	Slug             string `json:"slug"`
	IPA              string `json:"ipa"`
	AudioURL         string `json:"audioUrl"`
	PartOfSpeech     string `json:"partOfSpeech"`
	MeaningVI        string `json:"meaningVi"`
	MeaningEN        string `json:"meaningEn"`
	ExampleSentence  string `json:"exampleSentence"`
	ExampleMeaningVI string `json:"exampleMeaningVi"`
}

type CompleteActivityResponse struct {
	Activity       ActivityItem `json:"activity"`
	Score          int          `json:"score"`
	Passed         bool         `json:"passed"`
	CorrectCount   int          `json:"correctCount"`
	AnsweredCount  int          `json:"answeredCount"`
	TotalCount     int          `json:"totalCount"`
	XPAwarded      int          `json:"xpAwarded"`
	NextActivityID *uint        `json:"nextActivityId"`
	Lesson         LessonResult `json:"lesson"`
}

type LessonResult struct {
	ID             uint                `json:"id"`
	Status         models.LessonStatus `json:"status"`
	PassedRequired int                 `json:"passedRequired"`
	TotalRequired  int                 `json:"totalRequired"`
	// JustCompleted is true only for the activity that finished the lesson.
	JustCompleted bool `json:"justCompleted"`
	Score         *int `json:"score"`
	XPAwarded     int  `json:"xpAwarded"`
}
