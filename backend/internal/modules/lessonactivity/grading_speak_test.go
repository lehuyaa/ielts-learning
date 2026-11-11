package lessonactivity

import (
	"testing"

	"ielts-learning/backend/internal/models"
)

func TestGradeSpeak(t *testing.T) {
	result, err := Grade(models.ExerciseSpeak, []byte(`{"text":"Hi"}`), []byte(`{"minAccuracy":70}`), []byte(`{"attemptId":5}`))
	if err != nil || result.PronunciationAttempt != 5 {
		t.Fatalf("attempt answer: %+v %v", result, err)
	}
	if r := gradeSpeak(75, 70); !r.IsCorrect || r.ScoreRatio != 1 {
		t.Errorf("75/70: %+v", r)
	}
	if r := gradeSpeak(60, 70); r.IsCorrect || r.ScoreRatio != 0.5 {
		t.Errorf("60/70: %+v", r)
	}
	if r := gradeSpeak(40, 70); r.IsCorrect || r.ScoreRatio != 0 {
		t.Errorf("40/70: %+v", r)
	}
	if r, err := Grade(models.ExerciseSpeak, []byte(`{}`), []byte(`{}`), []byte(`{"selfCheck":true}`)); err != nil || !r.IsCorrect {
		t.Errorf("self-check: %+v %v", r, err)
	}
}
