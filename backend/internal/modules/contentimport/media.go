package contentimport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/lessonactivity"
	"ielts-learning/backend/internal/shared/ai"
)

var (
	ErrAIUnavailable = errors.New("ai unavailable")
	ErrAudioInvalid  = errors.New("audio invalid")
	ErrPassageInput  = errors.New("invalid passage")
	ErrExtractInput  = errors.New("invalid extract input")
)

// MediaStore keeps uploaded files (e.g. imported recordings) on disk; the API
// serves them under /media.
type MediaStore struct {
	Dir string
}

func (m MediaStore) Save(folder string, name string, data []byte) (string, error) {
	dir := filepath.Join(m.Dir, folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create media folder: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
		return "", fmt.Errorf("save media: %w", err)
	}
	return "/media/" + folder + "/" + name, nil
}

// PassageDraft is a transcribed recording for the admin to edit and save.
type PassageDraft struct {
	PassageInput
	Duration float64 `json:"duration"`
}

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	return strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(value), "-"), "-")
}

var audioExtensions = map[string]bool{".mp3": true, ".m4a": true, ".wav": true, ".webm": true, ".mp4": true, ".mpeg": true, ".ogg": true}

// Transcribe stores the recording and asks Whisper for the text with segment
// times; the result is a draft, nothing is written to the database.
func (s Service) Transcribe(ctx context.Context, filename string, data []byte, title string, slug string, genre string, source string) (PassageDraft, error) {
	if s.ai == nil || !s.ai.Configured() {
		return PassageDraft{}, ErrAIUnavailable
	}
	ext := strings.ToLower(filepath.Ext(filename))
	if !audioExtensions[ext] {
		return PassageDraft{}, fmt.Errorf("%w: use mp3, m4a, wav, webm or ogg", ErrAudioInvalid)
	}

	slug = slugify(defaultString(slug, title))
	if slug == "" {
		slug = "passage-" + time.Now().Format("20060102-150405")
	}
	transcription, err := s.ai.Transcribe(ctx, filename, data)
	if err != nil {
		return PassageDraft{}, fmt.Errorf("%w: %v", ErrAIUnavailable, err)
	}
	// Stored only once transcribed, so failed attempts leave no files behind.
	stored := fmt.Sprintf("%s-%d%s", slug, time.Now().Unix(), ext)
	audioURL, err := s.media.Save("passages", stored, data)
	if err != nil {
		return PassageDraft{}, err
	}

	segments := make([]map[string]any, 0, len(transcription.Segments))
	var lines []string
	for _, segment := range transcription.Segments {
		text := strings.TrimSpace(segment.Text)
		if text == "" {
			continue
		}
		segments = append(segments, map[string]any{
			"index": len(segments), "text": text,
			"start": round2(segment.Start), "end": round2(math.Max(segment.End, segment.Start+0.1)),
		})
		lines = append(lines, text)
	}
	if len(lines) == 0 && strings.TrimSpace(transcription.Text) != "" {
		lines = append(lines, strings.TrimSpace(transcription.Text))
	}
	segmentsJSON, _ := jsonMarshal(segments)

	return PassageDraft{
		PassageInput: PassageInput{
			Slug:      slug,
			CEFRLevel: "A1",
			Genre:     defaultString(genre, string(models.GenreDialogue)),
			Title:     defaultString(title, slug),
			BodyText:  strings.Join(lines, "\n"),
			AudioURL:  audioURL,
			Segments:  segmentsJSON,
			Source:    defaultString(source, "audio import: "+filename),
		},
		Duration: transcription.Duration,
	}, nil
}

// SavePassage validates a passage (e.g. an edited Whisper draft) and upserts it by slug.
func (s Service) SavePassage(input PassageInput, cefrLevel string) (lessonactivity.PassageCard, bool, error) {
	if errs := validatePassage(input); len(errs) > 0 {
		return lessonactivity.PassageCard{}, false, fmt.Errorf("%w: %s", ErrPassageInput, strings.Join(errs, "; "))
	}
	result := ImportResult{BundleErrors: []string{}, Warnings: []string{}, Items: []ItemResult{}}
	var card lessonactivity.PassageCard
	created := false
	err := s.repository.Transaction(func(tx Repository) error {
		bundle := Bundle{CEFRLevel: defaultString(cefrLevel, "A1"), Passages: []PassageInput{input}}
		if _, err := s.applyPassages(tx, bundle, true, &result); err != nil {
			return err
		}
		passage, _, err := tx.FindPassage(input.Slug)
		if err != nil {
			return err
		}
		card = lessonactivity.ToPassageCard(passage)
		created = len(result.Items) > 0 && result.Items[0].Action == ActionCreate
		return nil
	})
	return card, created, err
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

// aiClient is the part of ai.Client used here.
type aiClient interface {
	Configured() bool
	Transcribe(ctx context.Context, filename string, data []byte) (ai.Transcription, error)
	CompleteJSONWithParts(ctx context.Context, system string, parts []ai.Part, schema ai.Schema, out any) error
}

func jsonMarshal(value any) (json.RawMessage, error) {
	data, err := json.Marshal(value)
	return json.RawMessage(data), err
}
