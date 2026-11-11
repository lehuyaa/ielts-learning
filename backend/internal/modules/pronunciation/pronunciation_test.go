package pronunciation

import (
	"encoding/binary"
	"errors"
	"testing"
)

func wav(rate uint32, channels, bits uint16, seconds float64) []byte {
	dataBytes := int(float64(rate)*seconds) * int(channels) * int(bits/8)
	buf := make([]byte, 44+dataBytes)
	copy(buf[0:], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:], uint32(36+dataBytes))
	copy(buf[8:], "WAVE")
	copy(buf[12:], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:], 16)
	binary.LittleEndian.PutUint16(buf[20:], 1)
	binary.LittleEndian.PutUint16(buf[22:], channels)
	binary.LittleEndian.PutUint32(buf[24:], rate)
	binary.LittleEndian.PutUint32(buf[28:], rate*uint32(channels)*uint32(bits/8))
	binary.LittleEndian.PutUint16(buf[32:], channels*bits/8)
	binary.LittleEndian.PutUint16(buf[34:], bits)
	copy(buf[36:], "data")
	binary.LittleEndian.PutUint32(buf[40:], uint32(dataBytes))
	return buf
}

func TestValidateWAV(t *testing.T) {
	if err := validateWAV(wav(16000, 1, 16, 2)); err != nil {
		t.Fatalf("valid wav rejected: %v", err)
	}
	cases := map[string][]byte{
		"44.1 kHz":  wav(44100, 1, 16, 1),
		"stereo":    wav(16000, 2, 16, 1),
		"too long":  wav(16000, 1, 16, 17),
		"not a wav": []byte("OggS this is not a wav file at all, padding padding"),
	}
	for name, data := range cases {
		if err := validateWAV(data); !errors.Is(err, ErrInvalidAudio) {
			t.Errorf("%s: want ErrInvalidAudio, got %v", name, err)
		}
	}
}

func TestMapSound(t *testing.T) {
	cases := map[string]string{"i": "iː", "ɛ": "e", "ɹ": "r", "g": "ɡ", "ɚ": "ər", "ɑɹ": "ɑːr", "th": "θ", "ˈæ": "æ"}
	for azure, want := range cases {
		if got := mapSound(azure).display; got != want {
			t.Errorf("mapSound(%q) = %q, want %q", azure, got, want)
		}
	}
	if symbols := mapSound("ɚ").symbols; len(symbols) != 2 || symbols[0] != "ə" || symbols[1] != "r" {
		t.Errorf("ɚ symbols = %v", symbols)
	}
	if mapping := mapSound("x"); mapping.display != "x" || len(mapping.symbols) != 0 {
		t.Errorf("unknown sound mapped to %+v", mapping)
	}
}

func TestParseAzureFlatAndNested(t *testing.T) {
	flat := `{"RecognitionStatus":"Success","DisplayText":"Hello.","NBest":[{"AccuracyScore":90,"FluencyScore":80,"CompletenessScore":100,"PronScore":88,
		"Words":[{"Word":"hello","AccuracyScore":90,"ErrorType":"None","Phonemes":[{"Phoneme":"h","AccuracyScore":100},{"Phoneme":"ə","AccuracyScore":70}]}]}]}`
	nested := `{"RecognitionStatus":"Success","NBest":[{"Display":"Hello.","PronunciationAssessment":{"AccuracyScore":90,"FluencyScore":80,"CompletenessScore":100,"PronScore":88},
		"Words":[{"Word":"hello","PronunciationAssessment":{"AccuracyScore":90,"ErrorType":"None"},"Phonemes":[{"Phoneme":"h","PronunciationAssessment":{"AccuracyScore":100}},{"Phoneme":"ə","PronunciationAssessment":{"AccuracyScore":70}}]}]}]}`
	for name, body := range map[string]string{"flat": flat, "nested": nested} {
		got, err := parseAzure([]byte(body))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got.Accuracy != 90 || got.Fluency != 80 || got.Completeness != 100 || got.Pronunciation != 88 || got.RecognizedText != "Hello." {
			t.Errorf("%s: scores %+v", name, got)
		}
		if len(got.Words) != 1 || len(got.Words[0].Phonemes) != 2 || got.Words[0].Phonemes[1].Accuracy != 70 || got.Words[0].ErrorType != "None" {
			t.Errorf("%s: words %+v", name, got.Words)
		}
	}
	if _, err := parseAzure([]byte(`{"RecognitionStatus":"InitialSilenceTimeout"}`)); !errors.Is(err, ErrNoSpeech) {
		t.Errorf("silence: want ErrNoSpeech, got %v", err)
	}
}

func TestRollingAccuracy(t *testing.T) {
	accuracy, attempts := rollingAccuracy(0, 0, []int{80})
	if accuracy != 80 || attempts != 1 {
		t.Fatalf("first score: %d/%d", accuracy, attempts)
	}
	accuracy, attempts = rollingAccuracy(80, 1, []int{40})
	if accuracy != 60 || attempts != 2 {
		t.Fatalf("second score: %d/%d", accuracy, attempts)
	}
	// After many attempts a new score moves the average by about 1/30.
	accuracy, _ = rollingAccuracy(90, 100, []int{0})
	if accuracy != 87 {
		t.Fatalf("windowed: %d", accuracy)
	}
}

func TestMinAccuracy(t *testing.T) {
	if MinAccuracy([]byte(`{"selfCheck":true}`)) != 70 || MinAccuracy([]byte(`{"minAccuracy":80}`)) != 80 || MinAccuracy([]byte(`{"minAccuracy":0}`)) != 70 {
		t.Fatal("MinAccuracy defaults")
	}
}
