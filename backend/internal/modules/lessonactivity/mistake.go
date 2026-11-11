package lessonactivity

import (
	"strings"
	"time"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/mistakes"
)

// ExerciseMistake builds the mistake-log row for a wrong answer (Task 73).
// The skill tag falls back to the activity type (e.g. "reading").
func ExerciseMistake(userID uint, exercise models.Exercise, activityType string, submitted []byte, result GradeOutcome, now time.Time) (models.LearnerMistake, bool) {
	if result.IsCorrect || exercise.Type == models.ExerciseSpeak {
		return models.LearnerMistake{}, false
	}
	wrong, correct, ok := mistakes.DescribeExercise(exercise.Type, exercise.PayloadJSON, exercise.AnswerJSON, submitted)
	if !ok || strings.TrimSpace(correct) == "" {
		return models.LearnerMistake{}, false
	}
	explanation := exercise.ExplanationVI
	if result.AIFeedback != nil {
		if suggested := strings.TrimSpace(result.AIFeedback.Suggested); suggested != "" {
			correct = suggested
		}
		if result.AIFeedback.ExplanationVI != "" {
			explanation = result.AIFeedback.ExplanationVI
		}
	}
	if strings.EqualFold(strings.TrimSpace(wrong), strings.TrimSpace(correct)) {
		return models.LearnerMistake{}, false
	}
	tag := strings.TrimSpace(exercise.SkillTag)
	if tag == "" {
		tag = strings.ToLower(activityType)
	}
	if tag == "" {
		tag = "exercise"
	}
	return models.LearnerMistake{
		UserID:        userID,
		Source:        models.MistakeFromExercise,
		SkillTag:      tag,
		WrongText:     wrong,
		CorrectText:   correct,
		ExplanationVI: explanation,
		LastSeenAt:    now,
	}, true
}
