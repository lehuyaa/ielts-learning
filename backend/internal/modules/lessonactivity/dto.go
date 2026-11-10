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
