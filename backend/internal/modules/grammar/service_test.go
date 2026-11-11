package grammar

import (
	"testing"
	"time"

	"ielts-learning/backend/internal/models"
)

func TestOrderForPractice(t *testing.T) {
	base := time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)
	exercises := []models.Exercise{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}
	latest := map[uint]models.UserExerciseAnswer{
		1: {ExerciseID: 1, IsCorrect: true, AnsweredAt: base.Add(2 * time.Hour)},
		2: {ExerciseID: 2, IsCorrect: false, AnsweredAt: base},
		4: {ExerciseID: 4, IsCorrect: false, AnsweredAt: base.Add(time.Hour)},
		5: {ExerciseID: 5, IsCorrect: true, AnsweredAt: base},
	}

	got := orderForPractice(exercises, latest)
	want := []uint{4, 2, 3, 5, 1} // recent mistakes, older mistakes, new, oldest correct, newest correct
	for i, id := range want {
		if got[i].ID != id {
			t.Fatalf("order = %v, want %v", ids(got), want)
		}
	}
}

func ids(exercises []models.Exercise) []uint {
	result := make([]uint, len(exercises))
	for i, exercise := range exercises {
		result[i] = exercise.ID
	}
	return result
}
