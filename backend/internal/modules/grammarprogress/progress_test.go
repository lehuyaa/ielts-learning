package grammarprogress

import "testing"

func TestMasteryWindow(t *testing.T) {
	var recent []float64
	for i := 0; i < Window; i++ {
		recent = AddRecent(recent, 0)
	}
	if got := Mastery(recent); got != 0 {
		t.Fatalf("all wrong: mastery = %d, want 0", got)
	}
	// Ten correct answers replace the ten oldest wrong ones.
	for i := 0; i < 10; i++ {
		recent = AddRecent(recent, 1)
	}
	if len(recent) != Window {
		t.Fatalf("window length = %d, want %d", len(recent), Window)
	}
	if got := Mastery(recent); got != 50 {
		t.Fatalf("mastery = %d, want 50", got)
	}
	if got := Mastery(AddRecent(nil, 0.5)); got != 50 {
		t.Fatalf("partial credit: mastery = %d, want 50", got)
	}
	if got := Mastery(AddRecent(nil, 3)); got != 100 {
		t.Fatalf("ratio is clamped: mastery = %d, want 100", got)
	}
}

func TestSlugFromSkillTag(t *testing.T) {
	if slug, ok := SlugFromSkillTag("grammar:g1-to-be"); !ok || slug != "g1-to-be" {
		t.Fatalf("got %q %v", slug, ok)
	}
	for _, tag := range []string{"", "sound:iː", "grammar:"} {
		if _, ok := SlugFromSkillTag(tag); ok {
			t.Fatalf("tag %q should not be a grammar tag", tag)
		}
	}
}
