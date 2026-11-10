package lessonactivity

import (
	"context"
	"testing"

	"ielts-learning/backend/internal/models"
)

func TestSanitizeVerdict(t *testing.T) {
	cases := []struct {
		in   shortAnswerVerdict
		want float64
	}{
		{shortAnswerVerdict{IsCorrect: true, ScoreRatio: 0.2}, 1},
		{shortAnswerVerdict{IsCorrect: false, ScoreRatio: 0.7}, 0.5},
		{shortAnswerVerdict{IsCorrect: false, ScoreRatio: 0.3}, 0},
		{shortAnswerVerdict{IsCorrect: false, ScoreRatio: 1}, 0.5},
	}
	for _, c := range cases {
		if got := sanitizeVerdict(c.in).ScoreRatio; got != c.want {
			t.Errorf("sanitize(%+v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestResolveKeepsMatchResultWithoutAI(t *testing.T) {
	exercise := models.Exercise{ID: 1, Type: models.ExerciseTransform}
	matched := GradeResult{IsCorrect: true, ScoreRatio: 1}

	var grader *AIGrader
	outcome := grader.Resolve(context.Background(), 1, exercise, []byte(`{"text":"x"}`), matched)
	if outcome.GradedBy != GradedByMatch || outcome.AIStatus != "" || !outcome.IsCorrect {
		t.Fatalf("matched answer changed: %+v", outcome)
	}

	needsAI := GradeResult{IsCorrect: false, NeedsAI: true}
	outcome = grader.Resolve(context.Background(), 1, exercise, []byte(`{"text":"x"}`), needsAI)
	if outcome.AIStatus != AIStatusUnavailable || outcome.IsCorrect {
		t.Fatalf("nil grader should report UNAVAILABLE: %+v", outcome)
	}

	unconfigured := NewAIGrader(nil, nil, "m", 10)
	outcome = unconfigured.Resolve(context.Background(), 1, exercise, []byte(`{"text":"x"}`), needsAI)
	if outcome.AIStatus != AIStatusUnavailable {
		t.Fatalf("unconfigured grader should report UNAVAILABLE: %+v", outcome)
	}
}

func TestPromptHelpers(t *testing.T) {
	if transformInstructionEN("NEGATIVE") != "make it negative" {
		t.Fatal("unexpected instruction text")
	}
	if got := quoteAll([]string{"I work in Hanoi."}); got[0] != `"I work in Hanoi."` {
		t.Fatalf("quoteAll = %v", got)
	}
}
