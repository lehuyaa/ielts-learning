package contentimport

import (
	"encoding/json"
	"strings"
	"testing"
)

func exercise(typ string, payload string, answer string) ExerciseInput {
	return ExerciseInput{Type: typ, Prompt: "Prompt", Payload: json.RawMessage(payload), Answer: json.RawMessage(answer)}
}

func TestValidateExercise(t *testing.T) {
	tests := []struct {
		name    string
		input   ExerciseInput
		wantErr string
	}{
		{name: "multiple choice ok", input: exercise("MULTIPLE_CHOICE", `{"options":[{"id":"a","text":"x"},{"id":"b","text":"y"}]}`, `{"correct":"b"}`)},
		{name: "correct not in options", input: exercise("MULTIPLE_CHOICE", `{"options":[{"id":"a"},{"id":"b"}]}`, `{"correct":"c"}`), wantErr: "not one of the option ids"},
		{name: "duplicate option ids", input: exercise("MINIMAL_PAIR", `{"options":[{"id":"a"},{"id":"a"}]}`, `{"correct":"a"}`), wantErr: "duplicate option id"},
		{name: "ending sound plain options", input: exercise("ENDING_SOUND", `{"word":"works","options":["/s/","/z/","/ɪz/"]}`, `{"correct":"/s/"}`)},
		{name: "fill blank needs accepted", input: exercise("FILL_BLANK", `{"sentence":"I ___ ok."}`, `{"accepted":[]}`), wantErr: "answer.accepted"},
		{name: "reorder ok", input: exercise("REORDER", `{"tokens":["is","name","My","Huy","."]}`, `{"accepted":["My name is Huy."]}`)},
		{name: "reorder tokens mismatch", input: exercise("REORDER", `{"tokens":["is","name","My"]}`, `{"accepted":["My name is Huy."]}`), wantErr: "exactly the words"},
		{name: "reorder with contraction", input: exercise("REORDER", `{"tokens":["I'm","from","Vietnam","."]}`, `{"accepted":["I'm from Vietnam."]}`)},
		{name: "error correction index out of range", input: exercise("ERROR_CORRECTION", `{"tokens":["He","work"]}`, `{"wrongIndex":5,"accepted":["works"]}`), wantErr: "wrongIndex"},
		{name: "gap fill ok", input: exercise("LISTEN_GAP_FILL", `{"audioUrl":"/a.m4a","text":"At {{1}} by {{2}}."}`, `{"blanks":{"1":["nine"],"2":["bus"]}}`)},
		{name: "gap fill missing blank", input: exercise("LISTEN_GAP_FILL", `{"audioUrl":"/a.m4a","text":"At {{1}} by {{2}}."}`, `{"blanks":{"1":["nine"]}}`), wantErr: "answer.blanks.2"},
		{name: "dictation needs audio", input: exercise("DICTATION", `{}`, `{"accepted":["Hi."]}`), wantErr: "audioUrl"},
		{name: "word stress ok", input: exercise("WORD_STRESS", `{"word":"hello","syllables":["hel","lo"]}`, `{"stressIndex":1}`)},
		{name: "translate ok", input: exercise("TRANSLATE_VI_EN", `{"sourceVi":"Xin chào"}`, `{"reference":["Hello."]}`)},
		{name: "true false ok", input: exercise("TRUE_FALSE", `{"statement":"x"}`, `{"value":false}`)},
		{name: "speak ok", input: exercise("SPEAK", `{"text":"Hello"}`, `{"minAccuracy":70}`)},
		{name: "unknown type", input: exercise("ESSAY", `{}`, `{}`), wantErr: "unknown exercise type"},
		{name: "payload not object", input: exercise("TRUE_FALSE", `[1]`, `{}`), wantErr: "payload must be a JSON object"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validateExercise(tt.input)
			joined := strings.Join(errs, "; ")
			if tt.wantErr == "" && len(errs) > 0 {
				t.Fatalf("unexpected errors: %s", joined)
			}
			if tt.wantErr != "" && !strings.Contains(joined, tt.wantErr) {
				t.Fatalf("errors %q do not mention %q", joined, tt.wantErr)
			}
		})
	}
}

func TestValidateActivityVocabulary(t *testing.T) {
	known := map[string]bool{"a1-hello": true}
	errs := validateActivity(ActivityInput{
		Slug: "words", Type: "VOCABULARY", Title: "Words",
		Content: json.RawMessage(`{"vocabularySlugs":["a1-hello","a1-nope"]}`),
	}, known)
	if len(errs) != 1 || !strings.Contains(errs[0], "a1-nope") {
		t.Fatalf("errors = %v, want missing a1-nope", errs)
	}
	if errs := validateActivity(ActivityInput{Slug: "Bad Slug", Type: "DANCE", Title: ""}, known); len(errs) != 3 {
		t.Fatalf("errors = %v, want 3", errs)
	}
}

func TestParseBundleRejectsUnknownFields(t *testing.T) {
	if _, err := ParseBundle([]byte(`{"courseSlug":"x","topcs":[]}`)); err == nil {
		t.Fatal("expected unknown field error")
	}
}
