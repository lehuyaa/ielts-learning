package passage

import (
	"slices"
	"testing"
)

func TestWordCandidates(t *testing.T) {
	cases := map[string]string{
		"Works.":   "work",
		"studies":  "study",
		"studied":  "study",
		"stopped":  "stop",
		"running":  "run",
		"making":   "make",
		"liked":    "like",
		"boxes":    "box",
		"Lan's":    "lan",
		"“Hello,”": "hello",
	}
	for input, want := range cases {
		got := wordCandidates(input)
		if !slices.Contains(got, want) {
			t.Errorf("wordCandidates(%q) = %v, want it to contain %q", input, got, want)
		}
	}
	if got := wordCandidates("class"); slices.Contains(got, "clas") {
		t.Errorf("class should not lose its s: %v", got)
	}
}
