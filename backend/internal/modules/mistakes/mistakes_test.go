package mistakes

import (
	"testing"
	"time"

	"ielts-learning/backend/internal/models"
)

func TestDescribeExercise(t *testing.T) {
	cases := []struct {
		name                       string
		typ                        models.ExerciseType
		payload, answer, submitted string
		wantWrong, wantCorrect     string
		wantOK                     bool
	}{
		{"choice in sentence", models.ExerciseMultipleChoice, `{"sentence":"She ___ a nurse.","options":[{"id":"a","text":"is"},{"id":"b","text":"are"}]}`, `{"correct":"a"}`, `{"optionId":"b"}`, "She are a nurse.", "She is a nurse.", true},
		{"plain options", models.ExerciseMinimalPair, `{"options":["ship","sheep"]}`, `{"correct":"sheep"}`, `{"optionId":"ship"}`, "ship", "sheep", true},
		{"fill blank", models.ExerciseFillBlank, `{"sentence":"He ___ (work) here."}`, `{"accepted":["works"]}`, `{"text":"work"}`, "He work here.", "He works here.", true},
		{"reorder", models.ExerciseReorder, `{"tokens":["is","My","name","Huy","."]}`, `{"accepted":["My name is Huy."]}`, `{"tokens":["Is","my","name","Huy","."]}`, "Is my name Huy.", "My name is Huy.", true},
		{"error correction", models.ExerciseErrorCorrection, `{"tokens":["He","are","here","."]}`, `{"wrongIndex":1,"accepted":["is"]}`, `{"index":2,"text":"x"}`, "He are here.", "He is here.", true},
		{"true false", models.ExerciseTrueFalse, `{"statement":"Lan is a nurse."}`, `{"value":true}`, `{"value":false}`, "Lan is a nurse. → Sai", "Lan is a nurse. → Đúng", true},
		{"stress", models.ExerciseWordStress, `{"syllables":["tea","cher"]}`, `{"stressIndex":0}`, `{"stressIndex":1}`, "tea-CHER", "TEA-cher", true},
		{"speak skipped", models.ExerciseSpeak, `{"text":"Hi"}`, `{}`, `{"selfCheck":true}`, "", "", false},
	}
	for _, c := range cases {
		wrong, correct, ok := DescribeExercise(c.typ, []byte(c.payload), []byte(c.answer), []byte(c.submitted))
		if ok != c.wantOK || wrong != c.wantWrong || correct != c.wantCorrect {
			t.Errorf("%s: got (%q, %q, %v)", c.name, wrong, correct, ok)
		}
	}
}

func TestSkillLabelAndPractice(t *testing.T) {
	titles := map[string]string{"g1-to-be": "Động từ to be"}
	cases := map[string]string{
		"grammar:g1-to-be":    "Ngữ pháp: Động từ to be",
		"sound:θ-s":           "Phân biệt /θ/ và /s/",
		"sound:θ":             "Âm /θ/",
		"sound:ending-ed":     "Phát âm đuôi -ed",
		"listening:dictation": "Nghe chép chính tả",
		"reading":             "Đọc hiểu",
	}
	for tag, want := range cases {
		if got := SkillLabel(tag, titles); got != want {
			t.Errorf("SkillLabel(%q) = %q, want %q", tag, got, want)
		}
	}
	for tag, want := range map[string]bool{"grammar:g1-to-be": true, "sound:θ": true, "stress": true, "listening:dictation": true, "listening": false, "reading": false, "speaking": false} {
		if Practiceable(tag) != want {
			t.Errorf("Practiceable(%q) != %v", tag, want)
		}
	}
}

func TestSchedule(t *testing.T) {
	now := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	streak, interval, next, _ := schedule(0, 0, models.FlashcardRatingGood, now)
	if streak != 1 || interval != 1 || !next.Equal(now.AddDate(0, 0, 1)) {
		t.Fatalf("good: %d %d %v", streak, interval, next)
	}
	streak, interval, _, _ = schedule(streak, interval, models.FlashcardRatingGood, now)
	if streak != 2 || interval != 3 {
		t.Fatalf("good 2: %d %d", streak, interval)
	}
	streak, interval, next, _ = schedule(streak, interval, models.FlashcardRatingAgain, now)
	if streak != 0 || interval != 0 || !next.Equal(now.Add(10*time.Minute)) {
		t.Fatalf("again: %d %d %v", streak, interval, next)
	}
	if _, _, _, ok := schedule(0, 0, "MAYBE", now); ok {
		t.Fatal("invalid rating accepted")
	}
}
