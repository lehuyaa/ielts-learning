package contentimport

import (
	"encoding/json"
	"testing"
)

func TestValidatePhoneme(t *testing.T) {
	voiced := true
	valid := PhonemeInput{
		Symbol: "θ", Kind: "CONSONANT", Voiced: &voiced, MouthVI: "Đầu lưỡi giữa hai răng",
		ExampleWords: json.RawMessage(`[{"word":"think","ipa":"/θɪŋk/"}]`),
		MinimalPairs: json.RawMessage(`[{"a":"think","b":"sink"}]`),
	}
	if errs := validatePhoneme(valid); len(errs) != 0 {
		t.Fatalf("valid phoneme rejected: %v", errs)
	}

	vowel := valid
	vowel.Kind = "VOWEL"
	if errs := validatePhoneme(vowel); len(errs) == 0 {
		t.Fatal("vowel with voiced flag accepted")
	}

	consonant := valid
	consonant.Voiced = nil
	consonant.ExampleWords = json.RawMessage(`[]`)
	consonant.MinimalPairs = json.RawMessage(`[{"a":"think"}]`)
	if errs := validatePhoneme(consonant); len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %v", errs)
	}
}

func TestCollectPhonemeSymbols(t *testing.T) {
	bundle := Bundle{Topics: []TopicInput{{Lessons: []LessonInput{{Activities: []ActivityInput{
		{Content: json.RawMessage(`{"phonemes":["iː","ɪ"]}`)},
		{Content: json.RawMessage(`{"phonemes":["ɪ"],"rule":"x"}`)},
		{Content: json.RawMessage(`{"text":"hi"}`)},
	}}}}}}
	got := collectPhonemeSymbols(bundle)
	if len(got) != 2 || got[0] != "iː" || got[1] != "ɪ" {
		t.Fatalf("collectPhonemeSymbols = %v", got)
	}
}

func TestResolvePayloadPassage(t *testing.T) {
	ids := map[string]uint{"a1-dialogue": 7}
	got, errText := resolvePayloadPassage(json.RawMessage(`{"passageSlug":"a1-dialogue","segmentIndex":2}`), ids)
	if errText != "" || string(got) != `{"passageId":7,"segmentIndex":2}` {
		t.Fatalf("got %s %q", got, errText)
	}
	if _, errText := resolvePayloadPassage(json.RawMessage(`{"passageSlug":"missing"}`), ids); errText == "" {
		t.Fatal("missing passage accepted")
	}
	plain := json.RawMessage(`{"audioUrl":"/a.m4a"}`)
	if got, _ := resolvePayloadPassage(plain, ids); string(got) != string(plain) {
		t.Fatalf("plain payload changed: %s", got)
	}
}
