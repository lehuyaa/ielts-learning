package lessonactivity

import (
	"testing"

	"ielts-learning/backend/internal/models"
)

func TestScoreAttempt(t *testing.T) {
	exercises := []models.Exercise{{ID: 1, Points: 10}, {ID: 2, Points: 10}, {ID: 3, Points: 20}, {ID: 4, Points: 0}}

	tests := []struct {
		name    string
		answers []models.UserExerciseAnswer
		want    attemptTally
	}{
		{name: "nothing answered", want: attemptTally{Score: 0}},
		{
			name: "all correct",
			answers: []models.UserExerciseAnswer{
				{ExerciseID: 1, IsCorrect: true, ScoreRatio: 1}, {ExerciseID: 2, IsCorrect: true, ScoreRatio: 1},
				{ExerciseID: 3, IsCorrect: true, ScoreRatio: 1}, {ExerciseID: 4, IsCorrect: true, ScoreRatio: 1},
			},
			want: attemptTally{Score: 100, Correct: 4, Answered: 4},
		},
		{
			name: "points weighted with partial credit, unanswered counts zero",
			answers: []models.UserExerciseAnswer{
				{ExerciseID: 1, IsCorrect: true, ScoreRatio: 1},
				{ExerciseID: 3, IsCorrect: false, ScoreRatio: 0.5},
			},
			// (10 + 20*0.5) / (10+10+20+1) = 20/41
			want: attemptTally{Score: 49, Correct: 1, Answered: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := scoreAttempt(exercises, tt.answers); got != tt.want {
				t.Errorf("scoreAttempt() = %+v, want %+v", got, tt.want)
			}
		})
	}

	if got := scoreAttempt(nil, nil); got.Score != 100 {
		t.Errorf("activity without exercises should score 100, got %d", got.Score)
	}
}
