package pronunciation

import (
	"context"
	"errors"
	"testing"

	"ielts-learning/backend/internal/shared/ai"
)

func TestScoreTranscriptExactMatch(t *testing.T) {
	got := scoreTranscript("Hello, I’m Anna.", "hello i'm anna", 1.5)
	if got.Accuracy != 100 || got.Completeness != 100 || len(got.Words) != 3 {
		t.Fatalf("got %+v", got)
	}
	if got.Words[1].Word != "I’m" || got.Words[1].ErrorType != "None" {
		t.Fatalf("word keeps its spelling: %+v", got.Words[1])
	}
}

func TestScoreTranscriptMistakes(t *testing.T) {
	// "ship" heard as "sheep", "to" not heard, extra "um" ignored.
	got := scoreTranscript("I want to buy a ship", "um I want buy a sheep", 0)
	want := []string{"None", "None", "Omission", "None", "None", "Mispronunciation"}
	for i, word := range got.Words {
		if word.ErrorType != want[i] {
			t.Fatalf("word %d %q: got %s, want %s", i, word.Word, word.ErrorType, want[i])
		}
	}
	if got.Words[2].Accuracy != 0 || got.Words[5].Accuracy <= 0 || got.Words[5].Accuracy >= 60 {
		t.Fatalf("scores: %+v", got.Words)
	}
	if got.Completeness < 83 || got.Completeness > 84 {
		t.Fatalf("completeness %v", got.Completeness)
	}
}

func TestScoreTranscriptNumbers(t *testing.T) {
	got := scoreTranscript("It's eight thirty. I'm twenty five.", "It's 8:30. I'm 25.", 0)
	if got.Accuracy != 100 {
		t.Fatalf("numbers should match words: %+v", got.Words)
	}
	if spellOut("9:05") != "nine oh five" || spellOut("7:00") != "seven" || spellOut("120") != "one hundred and twenty" {
		t.Fatal("spellOut")
	}
}

type fakeTranscriber struct{ text string }

func (f fakeTranscriber) Configured() bool { return true }
func (f fakeTranscriber) Transcribe(context.Context, string, []byte) (ai.Transcription, error) {
	return ai.Transcription{Text: f.text, Duration: 2}, nil
}

func TestOpenAIScorerNoSpeech(t *testing.T) {
	scorer := &OpenAIScorer{client: fakeTranscriber{text: "  "}}
	if _, err := scorer.Assess(context.Background(), nil, "Hello"); !errors.Is(err, ErrNoSpeech) {
		t.Fatalf("err = %v", err)
	}
	scorer = &OpenAIScorer{client: fakeTranscriber{text: "Hello."}}
	if got, err := scorer.Assess(context.Background(), nil, "Hello"); err != nil || got.Accuracy != 100 {
		t.Fatalf("got %+v %v", got, err)
	}
}

func TestChooseScorer(t *testing.T) {
	azure := NewAzureClient("key", "southeastasia", "")
	none := NewAzureClient("", "", "")
	openai := ai.New("key", "model", "")
	cases := []struct {
		setting string
		azure   *AzureClient
		want    string
	}{
		{"auto", azure, ProviderAzure},
		{"auto", none, ProviderOpenAI},
		{"", nil, ProviderOpenAI},
		{"openai", azure, ProviderOpenAI},
		{"AZURE", none, ProviderAzure},
	}
	for _, c := range cases {
		if got := ChooseScorer(c.setting, c.azure, openai); got.Provider != c.want {
			t.Fatalf("%q: got %s want %s", c.setting, got.Provider, c.want)
		}
	}
	if ChooseScorer("openai", nil, ai.New("", "model", "")).client.Configured() {
		t.Fatal("OpenAI without a key is not configured")
	}
}
