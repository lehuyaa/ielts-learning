package lessonactivity

import "testing"

func TestCountSentences(t *testing.T) {
	cases := map[string]int{
		"My name is Huy. I am a developer.":  2,
		"Hi! How are you? I'm fine...":       3,
		"I live in Hanoi":                    1,
		"  ":                                 0,
		"I like 3.5 hours of sleep. Really.": 2, // "3.5" is not a sentence end
	}
	for text, want := range cases {
		if got := countSentences(text); got != want {
			t.Errorf("countSentences(%q) = %d, want %d", text, got, want)
		}
	}
}

func TestFinishReportLocatesMistakes(t *testing.T) {
	text := "My name is Huy. I works in Hà Nội. He work too."
	report := finishReport(WritingReport{
		Scores: WritingScores{TaskAchievement: 8, Grammar: 6, Vocabulary: 7, Coherence: 12},
		Mistakes: []WritingMistake{
			{Wrong: "I works", Correct: "I work", SkillTag: "grammar:g3-present-simple"},
			{Wrong: "He work", Correct: "He works", SkillTag: "grammar:g3-present-simple"},
			{Wrong: "not in text", Correct: "x"},
			{Wrong: "same", Correct: "same"},
		},
	}, text)

	if report.Scores.Coherence != 10 || report.Overall != 78 {
		t.Fatalf("scores %+v overall %d", report.Scores, report.Overall)
	}
	if len(report.Mistakes) != 3 || !report.RewriteRequired {
		t.Fatalf("mistakes %+v", report.Mistakes)
	}
	// UTF-16 offsets: "Hà Nội" has multi-byte letters before "He work".
	runes := []rune(text)
	first, second := report.Mistakes[0], report.Mistakes[1]
	if string(runes[first.Start:first.End]) != "I works" || string(runes[second.Start:second.End]) != "He work" {
		t.Fatalf("offsets %d-%d, %d-%d", first.Start, first.End, second.Start, second.End)
	}
	if report.Mistakes[2].Start != -1 {
		t.Fatalf("missing quote should have -1, got %d", report.Mistakes[2].Start)
	}
}
