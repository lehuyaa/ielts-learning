package lessonactivity

import (
	"testing"

	"ielts-learning/backend/internal/models"
)

func TestDrawExercises(t *testing.T) {
	var pool []models.Exercise
	for i := 1; i <= 30; i++ {
		tag := "grammar:g1-to-be"
		switch {
		case i > 24:
			tag = "listening:dictation"
		case i > 20:
			tag = "sound:iː-ɪ"
		}
		pool = append(pool, models.Exercise{ID: uint(i), SkillTag: tag})
	}
	first := drawExercises(append([]models.Exercise{}, pool...), 15, 42)
	again := drawExercises(append([]models.Exercise{}, pool...), 15, 42)
	if len(first) != 15 {
		t.Fatalf("drew %d", len(first))
	}
	kinds := map[string]int{}
	seen := map[uint]bool{}
	for i := range first {
		if first[i].ID != again[i].ID {
			t.Fatal("same seed must give the same draw")
		}
		if seen[first[i].ID] {
			t.Fatal("duplicate item")
		}
		seen[first[i].ID] = true
		kinds[first[i].SkillTag]++
	}
	// Every skill of the pool is represented (round-robin over kinds).
	if kinds["listening:dictation"] < 4 || kinds["sound:iː-ɪ"] < 4 {
		t.Fatalf("draw not spread over skills: %v", kinds)
	}
	other := drawExercises(append([]models.Exercise{}, pool...), 15, 43)
	same := true
	for i := range other {
		if other[i].ID != first[i].ID {
			same = false
		}
	}
	if same {
		t.Fatal("a new attempt should draw differently")
	}
}
