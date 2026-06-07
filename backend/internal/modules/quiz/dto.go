package quiz

import (
	"time"

	"ielts-learning/backend/internal/models"
)

type LessonResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	RequiredScore int    `json:"requiredScore"`
	TimeLimit     *int   `json:"timeLimitSeconds"`
}

type OptionResponse struct {
	ID    uint   `json:"id"`
	Label string `json:"label"`
	Text  string `json:"text"`
}

type QuestionResponse struct {
	ID      uint                    `json:"id"`
	Type    models.QuizQuestionType `json:"type"`
	Prompt  string                  `json:"prompt"`
	Points  int                     `json:"points"`
	Options []OptionResponse        `json:"options"`
}

type GetResponse struct {
	Lesson    LessonResponse     `json:"lesson"`
	Questions []QuestionResponse `json:"questions"`
}

type SubmitRequest struct {
	Answers []SubmitAnswerRequest `json:"answers" binding:"required"`
}

type SubmitAnswerRequest struct {
	QuestionID uint `json:"questionId" binding:"required"`
	OptionID   uint `json:"optionId" binding:"required"`
}

type CheckAnswerRequest struct {
	QuestionID uint `json:"questionId" binding:"required"`
	OptionID   uint `json:"optionId" binding:"required"`
}

type CheckAnswerResponse struct {
	QuestionID       uint   `json:"questionId"`
	SelectedOptionID uint   `json:"selectedOptionId"`
	CorrectOptionID  uint   `json:"correctOptionId"`
	IsCorrect        bool   `json:"isCorrect"`
	Explanation      string `json:"explanation"`
	EarnedPoints     int    `json:"earnedPoints"`
}

type SubmitResultResponse struct {
	QuestionID       uint   `json:"questionId"`
	SelectedOptionID *uint  `json:"selectedOptionId"`
	CorrectOptionID  uint   `json:"correctOptionId"`
	IsCorrect        bool   `json:"isCorrect"`
	Explanation      string `json:"explanation"`
	PointsAwarded    int    `json:"pointsAwarded"`
}

type SubmitResponse struct {
	AttemptID      uint                   `json:"attemptId"`
	LessonID       uint                   `json:"lessonId"`
	Score          int                    `json:"score"`
	RequiredScore  int                    `json:"requiredScore"`
	Passed         bool                   `json:"passed"`
	CorrectCount   int                    `json:"correctCount"`
	TotalQuestions int                    `json:"totalQuestions"`
	EarnedXP       int                    `json:"earnedXp"`
	CompletedAt    *time.Time             `json:"completedAt"`
	Results        []SubmitResultResponse `json:"results"`
}
