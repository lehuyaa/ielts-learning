// Package pronunciation scores recordings with Azure pronunciation assessment
// (docs/foundation-spec.md §7, Tasks 71-72) and keeps per-sound progress.
package pronunciation

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
	"unicode"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/aiusage"
	"ielts-learning/backend/internal/modules/mistakes"
)

const (
	// MaxSeconds is the longest recording accepted (spec: recorder max 15 s).
	MaxSeconds = 15
	// MaxReferenceChars matches pronunciation_attempts.reference_text.
	MaxReferenceChars = 500
	// rollingWindow approximates "average of the last 30 occurrences" (spec §7).
	rollingWindow = 30
	// A sound becomes a learner mistake below this average after minMistakeAttempts.
	mistakeAccuracy    = 60
	minMistakeAttempts = 5
	defaultMinAccuracy = 70
)

var (
	ErrUnavailable      = errors.New("pronunciation scoring is not configured")
	ErrLimitReached     = errors.New("daily pronunciation limit reached")
	ErrInvalidAudio     = errors.New("audio must be a 16 kHz mono 16-bit WAV of at most 15 seconds")
	ErrInvalidReference = errors.New("invalid reference text")
	ErrExerciseNotFound = errors.New("speaking exercise not found")
	ErrScoringFailed    = errors.New("pronunciation scoring failed")
)

// assessor is the part of AzureClient the service needs (fake in tests).
type assessor interface {
	Configured() bool
	Assess(ctx context.Context, wav []byte, referenceText string) (Assessment, error)
}

type Service struct {
	db         *gorm.DB
	client     assessor
	dailyLimit int
	now        func() time.Time
}

func NewService(db *gorm.DB, client *AzureClient, dailyLimit int) Service {
	service := Service{db: db, dailyLimit: dailyLimit, now: time.Now}
	if client != nil {
		service.client = client
	}
	return service
}

func (s Service) enabled() bool {
	return s.client != nil && s.client.Configured()
}

func (s Service) Status() StatusResponse {
	return StatusResponse{Enabled: s.enabled(), MaxSeconds: MaxSeconds, DailyLimit: s.dailyLimit}
}

// AssessInput is one recording to score.
type AssessInput struct {
	Audio         []byte
	ReferenceText string
	ExerciseID    *uint
	Source        string
}

func (s Service) Assess(ctx context.Context, userID uint, input AssessInput) (AssessResponse, error) {
	if !s.enabled() {
		return AssessResponse{}, ErrUnavailable
	}

	reference := strings.Join(strings.Fields(input.ReferenceText), " ")
	source := normalizeSource(input.Source, input.ExerciseID != nil)
	minAccuracy := 0
	if input.ExerciseID != nil {
		text, min, err := s.speakExercise(*input.ExerciseID)
		if err != nil {
			return AssessResponse{}, err
		}
		reference, minAccuracy, source = text, min, models.PronunciationExercise
	}
	if err := validateReference(reference); err != nil {
		return AssessResponse{}, err
	}
	if err := validateWAV(input.Audio); err != nil {
		return AssessResponse{}, err
	}

	allowed, err := aiusage.Take(s.db, userID, models.AIUsagePronunciation, s.dailyLimit, s.now())
	if err != nil {
		return AssessResponse{}, err
	}
	if !allowed {
		return AssessResponse{}, ErrLimitReached
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	assessment, err := s.client.Assess(callCtx, input.Audio, reference)
	if err != nil {
		if errors.Is(err, ErrNoSpeech) {
			return AssessResponse{}, ErrNoSpeech
		}
		log.Printf("pronunciation: assess: %v", err)
		return AssessResponse{}, ErrScoringFailed
	}

	words := toWords(assessment.Words)
	wordsJSON, _ := json.Marshal(words)
	attempt := models.PronunciationAttempt{
		UserID:        userID,
		ReferenceText: truncate(reference, MaxReferenceChars),
		Source:        source,
		ExerciseID:    input.ExerciseID,
		Accuracy:      round(assessment.Accuracy),
		Fluency:       round(assessment.Fluency),
		Completeness:  round(assessment.Completeness),
		Pronunciation: round(assessment.Pronunciation),
		WordsJSON:     datatypes.JSON(wordsJSON),
	}

	now := s.now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&attempt).Error; err != nil {
			return fmt.Errorf("save pronunciation attempt: %w", err)
		}
		return recordSounds(tx, userID, words, now)
	})
	if err != nil {
		return AssessResponse{}, err
	}

	response := AssessResponse{
		AttemptID:      attempt.ID,
		ReferenceText:  reference,
		RecognizedText: assessment.RecognizedText,
		Accuracy:       attempt.Accuracy,
		Fluency:        attempt.Fluency,
		Completeness:   attempt.Completeness,
		Pronunciation:  attempt.Pronunciation,
		Words:          words,
	}
	if input.ExerciseID != nil {
		response.MinAccuracy = &minAccuracy
		passed := attempt.Accuracy >= minAccuracy
		response.Passed = &passed
	}
	return response, nil
}

// speakExercise returns a SPEAK exercise's sentence and pass mark.
func (s Service) speakExercise(id uint) (string, int, error) {
	var exercise models.Exercise
	err := s.db.Where("id = ? AND type = ?", id, models.ExerciseSpeak).Take(&exercise).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", 0, ErrExerciseNotFound
	}
	if err != nil {
		return "", 0, err
	}
	var payload struct {
		Text string `json:"text"`
	}
	_ = json.Unmarshal(exercise.PayloadJSON, &payload)
	return payload.Text, MinAccuracy(exercise.AnswerJSON), nil
}

// MinAccuracy is a SPEAK exercise's pass mark (answer.minAccuracy, default 70).
func MinAccuracy(answerJSON []byte) int {
	var key struct {
		MinAccuracy *int `json:"minAccuracy"`
	}
	_ = json.Unmarshal(answerJSON, &key)
	if key.MinAccuracy == nil || *key.MinAccuracy <= 0 || *key.MinAccuracy > 100 {
		return defaultMinAccuracy
	}
	return *key.MinAccuracy
}

func toWords(assessed []AssessedWord) []WordResult {
	words := make([]WordResult, 0, len(assessed))
	for _, word := range assessed {
		result := WordResult{Word: word.Word, Accuracy: round(word.Accuracy), ErrorType: word.ErrorType, Phonemes: []PhonemeResult{}}
		for _, phoneme := range word.Phonemes {
			mapping := mapSound(phoneme.Phoneme)
			result.Phonemes = append(result.Phonemes, PhonemeResult{
				Azure:    phoneme.Phoneme,
				Display:  mapping.display,
				Symbols:  append([]string{}, mapping.symbols...),
				Accuracy: round(phoneme.Accuracy),
			})
		}
		words = append(words, result)
	}
	return words
}

type soundSample struct {
	scores    []int
	worstWord string
	worst     int
}

// recordSounds updates user_phoneme_progresses with every scored phoneme
// and logs sounds that stay weak as learner mistakes (spec §7).
func recordSounds(tx *gorm.DB, userID uint, words []WordResult, now time.Time) error {
	samples := map[string]*soundSample{}
	var order []string
	for _, word := range words {
		if word.ErrorType == "Omission" || word.ErrorType == "Insertion" {
			continue
		}
		for _, phoneme := range word.Phonemes {
			for _, symbol := range phoneme.Symbols {
				sample, ok := samples[symbol]
				if !ok {
					sample = &soundSample{worst: 101}
					samples[symbol] = sample
					order = append(order, symbol)
				}
				sample.scores = append(sample.scores, phoneme.Accuracy)
				if phoneme.Accuracy < sample.worst {
					sample.worst, sample.worstWord = phoneme.Accuracy, word.Word
				}
			}
		}
	}
	if len(order) == 0 {
		return nil
	}

	var phonemes []models.Phoneme
	if err := tx.Where("symbol IN ?", order).Find(&phonemes).Error; err != nil {
		return fmt.Errorf("find phonemes: %w", err)
	}
	ids := make(map[string]uint, len(phonemes))
	for _, phoneme := range phonemes {
		ids[phoneme.Symbol] = phoneme.ID
	}

	for _, symbol := range order {
		phonemeID, ok := ids[symbol]
		if !ok {
			continue
		}
		sample := samples[symbol]
		seed := models.UserPhonemeProgress{UserID: userID, PhonemeID: phonemeID}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&seed).Error; err != nil {
			return fmt.Errorf("create phoneme progress: %w", err)
		}
		var progress models.UserPhonemeProgress
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND phoneme_id = ?", userID, phonemeID).Take(&progress).Error; err != nil {
			return fmt.Errorf("lock phoneme progress: %w", err)
		}
		progress.Accuracy, progress.Attempts = rollingAccuracy(progress.Accuracy, progress.Attempts, sample.scores)
		progress.LastPracticed = &now
		if err := tx.Model(&progress).Updates(map[string]any{
			"accuracy":          progress.Accuracy,
			"attempts":          progress.Attempts,
			"last_practiced_at": now,
		}).Error; err != nil {
			return fmt.Errorf("save phoneme progress: %w", err)
		}

		if progress.Attempts >= minMistakeAttempts && progress.Accuracy < mistakeAccuracy && sample.worst < mistakeAccuracy {
			err := mistakes.Upsert(tx, models.LearnerMistake{
				UserID:      userID,
				Source:      models.MistakeFromPronunciation,
				SkillTag:    "sound:" + symbol,
				WrongText:   sample.worstWord,
				CorrectText: "/" + symbol + "/",
				ExplanationVI: fmt.Sprintf("Âm /%s/ (ví dụ trong '%s') của bạn đang đạt trung bình %d/100. Mở thẻ âm /%s/ để xem cách đặt miệng và luyện thêm.",
					symbol, sample.worstWord, progress.Accuracy, symbol),
				LastSeenAt: now,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// rollingAccuracy folds new scores into an average that weighs roughly the
// last 30 occurrences.
func rollingAccuracy(accuracy, attempts int, scores []int) (int, int) {
	value := float64(accuracy)
	for _, score := range scores {
		weight := math.Min(float64(attempts), rollingWindow-1)
		value = (value*weight + float64(score)) / (weight + 1)
		attempts++
	}
	return int(math.Round(value)), attempts
}

func normalizeSource(source string, hasExercise bool) string {
	switch strings.ToUpper(strings.TrimSpace(source)) {
	case models.PronunciationWord:
		return models.PronunciationWord
	case models.PronunciationShadowing:
		return models.PronunciationShadowing
	case models.PronunciationReadAloud:
		return models.PronunciationReadAloud
	}
	if hasExercise {
		return models.PronunciationExercise
	}
	return models.PronunciationWord
}

func validateReference(text string) error {
	if text == "" || len([]rune(text)) > MaxReferenceChars {
		return ErrInvalidReference
	}
	for _, r := range text {
		if unicode.IsLetter(r) {
			return nil
		}
	}
	return ErrInvalidReference
}

// validateWAV accepts the PCM WAV the recorder produces: 16 kHz, mono,
// 16-bit, at most MaxSeconds (+1 s of slack).
func validateWAV(data []byte) error {
	if len(data) < 44 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return ErrInvalidAudio
	}
	var channels, bits uint16
	var rate uint32
	dataBytes := -1
	for offset := 12; offset+8 <= len(data); {
		id := string(data[offset : offset+4])
		size := int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
		body := offset + 8
		switch id {
		case "fmt ":
			if body+16 > len(data) {
				return ErrInvalidAudio
			}
			format := binary.LittleEndian.Uint16(data[body : body+2])
			channels = binary.LittleEndian.Uint16(data[body+2 : body+4])
			rate = binary.LittleEndian.Uint32(data[body+4 : body+8])
			bits = binary.LittleEndian.Uint16(data[body+14 : body+16])
			if format != 1 {
				return ErrInvalidAudio
			}
		case "data":
			dataBytes = min(size, len(data)-body)
		}
		if dataBytes >= 0 {
			break
		}
		offset = body + size + size%2
	}
	if channels != 1 || bits != 16 || rate != 16000 || dataBytes <= 0 {
		return ErrInvalidAudio
	}
	if dataBytes > (MaxSeconds+1)*16000*2 {
		return ErrInvalidAudio
	}
	return nil
}

func round(value float64) int {
	return int(math.Round(math.Max(0, math.Min(100, value))))
}

func truncate(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}
