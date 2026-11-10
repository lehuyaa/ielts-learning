package grammar

import (
	"encoding/json"
	"time"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/lessonactivity"
)

const practiceSize = 10

type ListResponse struct {
	Items []GrammarPointItem `json:"items"`
}

type GrammarPointItem struct {
	ID            uint             `json:"id"`
	Slug          string           `json:"slug"`
	CEFRLevel     string           `json:"cefrLevel"`
	Title         string           `json:"title"`
	TitleVI       string           `json:"titleVi"`
	OrderIndex    int              `json:"orderIndex"`
	ExerciseCount int              `json:"exerciseCount"`
	Progress      ProgressResponse `json:"progress"`
}

type ProgressResponse struct {
	// Mastery is the average score (0-100) of the last 20 answers.
	Mastery         int        `json:"mastery"`
	CorrectCount    int        `json:"correctCount"`
	WrongCount      int        `json:"wrongCount"`
	Practised       bool       `json:"practised"`
	LastPracticedAt *time.Time `json:"lastPracticedAt"`
}

type DetailResponse struct {
	GrammarPoint lessonactivity.GrammarPointCard `json:"grammarPoint"`
	Progress     ProgressResponse                `json:"progress"`
	// ExerciseCount counts practice exercises available to the learner.
	ExerciseCount int `json:"exerciseCount"`
}

type PracticeResponse struct {
	Exercises []PracticeExercise `json:"exercises"`
}

// PracticeExercise never contains the answer key.
type PracticeExercise struct {
	ID         uint                `json:"id"`
	Type       models.ExerciseType `json:"type"`
	Prompt     string              `json:"prompt"`
	Payload    json.RawMessage     `json:"payload" swaggertype:"object"`
	Points     int                 `json:"points"`
	OrderIndex int                 `json:"orderIndex"`
	// LastResult is "WRONG", "CORRECT" or "NEW" for the learner's latest answer.
	LastResult string `json:"lastResult"`
}

type PracticeAnswerRequest struct {
	ExerciseID uint            `json:"exerciseId"`
	Answer     json.RawMessage `json:"answer" swaggertype:"object"`
}

type PracticeAnswerResponse struct {
	ExerciseID    uint                       `json:"exerciseId"`
	Attempt       int                        `json:"attempt"`
	IsCorrect     bool                       `json:"isCorrect"`
	ScoreRatio    float64                    `json:"scoreRatio"`
	Answer        json.RawMessage            `json:"answer" swaggertype:"object"`
	CorrectAnswer any                        `json:"correctAnswer"`
	ExplanationVI string                     `json:"explanationVi"`
	AnsweredAt    time.Time                  `json:"answeredAt"`
	GradedBy      string                     `json:"gradedBy,omitempty"`
	AIFeedback    *lessonactivity.AIFeedback `json:"aiFeedback,omitempty"`
	AIStatus      string                     `json:"aiStatus,omitempty"`
	Progress      ProgressResponse           `json:"progress"`
}
