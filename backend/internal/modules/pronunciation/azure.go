package pronunciation

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

// ErrNoSpeech means Azure heard nothing it could score.
var ErrNoSpeech = errors.New("no speech recognised")

// AzureClient calls Azure Speech pronunciation assessment through the REST
// API for short audio (docs/foundation-spec.md §7): en-US, scripted mode with
// the reference text, phoneme granularity, IPA phonemes.
type AzureClient struct {
	key      string
	endpoint string
	http     *http.Client
}

// NewAzureClient returns a client for a region (e.g. "southeastasia").
// endpoint overrides https://{region}.stt.speech.microsoft.com (tests).
func NewAzureClient(key, region, endpoint string) *AzureClient {
	key = strings.TrimSpace(key)
	region = strings.TrimSpace(region)
	endpoint = strings.TrimRight(strings.TrimSpace(endpoint), "/")
	if endpoint == "" && region != "" {
		endpoint = fmt.Sprintf("https://%s.stt.speech.microsoft.com", region)
	}
	return &AzureClient{key: key, endpoint: endpoint, http: &http.Client{Timeout: 30 * time.Second}}
}

// Configured reports whether a key and region (or endpoint) are set.
func (c *AzureClient) Configured() bool {
	return c != nil && c.key != "" && c.endpoint != ""
}

// Assessment is Azure's result, phonemes still in Azure's notation.
type Assessment struct {
	RecognizedText string
	Accuracy       float64
	Fluency        float64
	Completeness   float64
	Pronunciation  float64
	Words          []AssessedWord
}

type AssessedWord struct {
	Word      string
	Accuracy  float64
	ErrorType string
	Phonemes  []AssessedPhoneme
}

type AssessedPhoneme struct {
	Phoneme  string
	Accuracy float64
}

// Assess scores a 16 kHz mono 16-bit WAV recording against referenceText.
func (c *AzureClient) Assess(ctx context.Context, wav []byte, referenceText string) (Assessment, error) {
	config, _ := json.Marshal(map[string]any{
		"ReferenceText":   referenceText,
		"GradingSystem":   "HundredMark",
		"Granularity":     "Phoneme",
		"Dimension":       "Comprehensive",
		"EnableMiscue":    true,
		"PhonemeAlphabet": "IPA",
	})
	url := c.endpoint + "/speech/recognition/conversation/cognitiveservices/v1?language=en-US&format=detailed"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(wav))
	if err != nil {
		return Assessment{}, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.key)
	req.Header.Set("Content-Type", "audio/wav; codecs=audio/pcm; samplerate=16000")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Pronunciation-Assessment", base64.StdEncoding.EncodeToString(config))

	resp, err := c.http.Do(req)
	if err != nil {
		return Assessment{}, fmt.Errorf("azure speech: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if resp.StatusCode != http.StatusOK {
		return Assessment{}, fmt.Errorf("azure speech: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return parseAzure(body)
}

// azureScores appear flat on the REST result and nested under
// "PronunciationAssessment" in the SDK/newer format; both are read.
type azureScores struct {
	AccuracyScore     *float64 `json:"AccuracyScore"`
	FluencyScore      *float64 `json:"FluencyScore"`
	CompletenessScore *float64 `json:"CompletenessScore"`
	PronScore         *float64 `json:"PronScore"`
	ErrorType         string   `json:"ErrorType"`
}

type azureResult struct {
	RecognitionStatus string `json:"RecognitionStatus"`
	DisplayText       string `json:"DisplayText"`
	NBest             []struct {
		Display string `json:"Display"`
		azureScores
		PronunciationAssessment *azureScores `json:"PronunciationAssessment"`
		Words                   []struct {
			Word string `json:"Word"`
			azureScores
			PronunciationAssessment *azureScores `json:"PronunciationAssessment"`
			Phonemes                []struct {
				Phoneme string `json:"Phoneme"`
				azureScores
				PronunciationAssessment *azureScores `json:"PronunciationAssessment"`
			} `json:"Phonemes"`
		} `json:"Words"`
	} `json:"NBest"`
}

func merge(flat azureScores, nested *azureScores) azureScores {
	if nested == nil {
		return flat
	}
	if nested.AccuracyScore != nil {
		flat.AccuracyScore = nested.AccuracyScore
	}
	if nested.FluencyScore != nil {
		flat.FluencyScore = nested.FluencyScore
	}
	if nested.CompletenessScore != nil {
		flat.CompletenessScore = nested.CompletenessScore
	}
	if nested.PronScore != nil {
		flat.PronScore = nested.PronScore
	}
	if nested.ErrorType != "" {
		flat.ErrorType = nested.ErrorType
	}
	return flat
}

func value(score *float64) float64 {
	if score == nil || math.IsNaN(*score) {
		return 0
	}
	return math.Max(0, math.Min(100, *score))
}

func parseAzure(body []byte) (Assessment, error) {
	var result azureResult
	if err := json.Unmarshal(body, &result); err != nil {
		return Assessment{}, fmt.Errorf("azure speech: decode: %w", err)
	}
	if result.RecognitionStatus != "" && result.RecognitionStatus != "Success" {
		if result.RecognitionStatus == "NoMatch" || result.RecognitionStatus == "InitialSilenceTimeout" || result.RecognitionStatus == "BabbleTimeout" {
			return Assessment{}, ErrNoSpeech
		}
		return Assessment{}, fmt.Errorf("azure speech: recognition status %s", result.RecognitionStatus)
	}
	if len(result.NBest) == 0 {
		return Assessment{}, ErrNoSpeech
	}
	best := result.NBest[0]
	scores := merge(best.azureScores, best.PronunciationAssessment)
	assessment := Assessment{
		RecognizedText: firstNonEmpty(result.DisplayText, best.Display),
		Accuracy:       value(scores.AccuracyScore),
		Fluency:        value(scores.FluencyScore),
		Completeness:   value(scores.CompletenessScore),
		Pronunciation:  value(scores.PronScore),
	}
	for _, word := range best.Words {
		ws := merge(word.azureScores, word.PronunciationAssessment)
		assessed := AssessedWord{Word: word.Word, Accuracy: value(ws.AccuracyScore), ErrorType: firstNonEmpty(ws.ErrorType, "None")}
		for _, phoneme := range word.Phonemes {
			ps := merge(phoneme.azureScores, phoneme.PronunciationAssessment)
			assessed.Phonemes = append(assessed.Phonemes, AssessedPhoneme{Phoneme: phoneme.Phoneme, Accuracy: value(ps.AccuracyScore)})
		}
		assessment.Words = append(assessment.Words, assessed)
	}
	if assessment.Pronunciation == 0 && assessment.Accuracy > 0 {
		assessment.Pronunciation = assessment.Accuracy
	}
	return assessment, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
