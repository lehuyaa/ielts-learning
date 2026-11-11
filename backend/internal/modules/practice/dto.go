package practice

import (
	"encoding/json"
	"time"

	"ielts-learning/backend/internal/models"
)

type GenerateRequest struct {
	SkillTag string `json:"skillTag"`
}

// SetResponse is a practice set as the player shows it: exercises are
// numbered 1..n (exerciseId), answers so far allow resuming.
type SetResponse struct {
	ID       uint   `json:"id"`
	SkillTag string `json:"skillTag"`
	Label    string `json:"label"`
	// Source is AI (generated) or BANK (picked from lessons).
	Source    string         `json:"source"`
	Day       string         `json:"day"`
	Exercises []Exercise     `json:"exercises"`
	Answers   []AnswerResult `json:"answers"`
}

type Exercise struct {
	ID         uint                `json:"id"`
	Type       models.ExerciseType `json:"type"`
	Prompt     string              `json:"prompt"`
	Payload    json.RawMessage     `json:"payload" swaggertype:"object"`
	Points     int                 `json:"points"`
	OrderIndex int                 `json:"orderIndex"`
}

type AnswerRequest struct {
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
