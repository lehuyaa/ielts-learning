// Package practice builds extra practice for one weakness (Task 74): AI
// generates 10 validated exercises for grammar / vocabulary tags; sound,
// stress, dictation and number tags reuse lesson exercises (they carry
// audio). A set is cached per user, skill tag and day.
package practice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/aiusage"
	"ielts-learning/backend/internal/modules/grammarprogress"
	"ielts-learning/backend/internal/modules/lessonactivity"
	"ielts-learning/backend/internal/modules/mistakes"
	"ielts-learning/backend/internal/shared/ai"
)

const (
	setSize = 10
	// requestCount asks for a few extra items: some may fail validation.
	requestCount = 12
	// minItems is the smallest usable set.
	minItems = 5
)

var (
	ErrInvalidTag     = errors.New("invalid skill tag")
	ErrNotPracticable = errors.New("this skill is practised in lessons")
	ErrNotEnough      = errors.New("not enough exercises for this skill")
	ErrSetNotFound    = errors.New("practice set not found")
	ErrItemNotFound   = errors.New("exercise not in this set")
	ErrAlreadyAnswer  = errors.New("exercise already answered")
	ErrLimitReached   = errors.New("daily practice limit reached")
	ErrAIUnavailable  = errors.New("AI practice is unavailable")
)

// jsonCompleter is the part of ai.Client the service needs (fake in tests).
type jsonCompleter interface {
	Configured() bool
	CompleteJSON(ctx context.Context, messages []ai.Message, schema ai.Schema, out any) error
}

type Service struct {
	db         *gorm.DB
	ai         jsonCompleter
	dailyLimit int
	timeout    time.Duration
	now        func() time.Time
}

func NewService(db *gorm.DB, client *ai.Client, dailyLimit int) Service {
	service := Service{db: db, dailyLimit: dailyLimit, timeout: 60 * time.Second, now: time.Now}
	if client != nil {
		service.ai = client
	}
	return service
}

// Item is one stored exercise of a set (the answer key stays on the server).
type Item struct {
	Type             models.ExerciseType `json:"type"`
	Prompt           string              `json:"prompt"`
	Payload          json.RawMessage     `json:"payload"`
	Answer           json.RawMessage     `json:"answer"`
	ExplanationVI    string              `json:"explanationVi"`
	SkillTag         string              `json:"skillTag"`
	SourceExerciseID *uint               `json:"sourceExerciseId,omitempty"`
}

// Generate returns today's set for a skill tag, building it when needed.
func (s Service) Generate(ctx context.Context, userID uint, tag string) (SetResponse, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" || len(tag) > 80 {
		return SetResponse{}, ErrInvalidTag
	}
	if !mistakes.Practiceable(tag) {
		return SetResponse{}, ErrNotPracticable
	}
	day := s.now().Format("2006-01-02")
	var existing models.PracticeSet
	err := s.db.Where("user_id = ? AND skill_tag = ? AND day = ?", userID, tag, day).Take(&existing).Error
	if err == nil {
		return s.response(existing)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return SetResponse{}, fmt.Errorf("find practice set: %w", err)
	}

	source := models.PracticeFromBank
	var items []Item
	if mistakes.PracticeFromBank(tag) {
		items, err = s.fromBank(userID, tag)
	} else {
		source = models.PracticeFromAI
		items, err = s.fromAI(ctx, userID, tag)
	}
	if err != nil {
		return SetResponse{}, err
	}
	if len(items) < minItems {
		return SetResponse{}, ErrNotEnough
	}

	itemsJSON, _ := json.Marshal(items)
	set := models.PracticeSet{UserID: userID, SkillTag: tag, Day: day, Source: source, ItemsJSON: datatypes.JSON(itemsJSON), ResultsJSON: datatypes.JSON("{}")}
	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&set).Error; err != nil {
		return SetResponse{}, fmt.Errorf("save practice set: %w", err)
	}
	// A concurrent request may have saved first: return the stored set.
	if err := s.db.Where("user_id = ? AND skill_tag = ? AND day = ?", userID, tag, day).Take(&set).Error; err != nil {
		return SetResponse{}, fmt.Errorf("load practice set: %w", err)
	}
	return s.response(set)
}

func (s Service) fromAI(ctx context.Context, userID uint, tag string) ([]Item, error) {
	if s.ai == nil || !s.ai.Configured() {
		return nil, ErrAIUnavailable
	}
	f, err := s.focus(userID, tag)
	if err != nil {
		return nil, err
	}
	allowed, err := aiusage.Take(s.db, userID, models.AIUsagePractice, s.dailyLimit, s.now())
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrLimitReached
	}
	items, err := s.generateWithAI(ctx, f)
	if err != nil {
		log.Printf("practice: generate %s: %v", tag, err)
		return nil, ErrAIUnavailable
	}
	return items, nil
}

// focus collects what the AI needs to know about the weakness.
func (s Service) focus(userID uint, tag string) (focus, error) {
	f := focus{Tag: tag, Level: "A1", Title: mistakes.SkillLabel(tag, nil)}
	if slug, ok := strings.CutPrefix(tag, "grammar:"); ok {
		var point models.GrammarPoint
		err := s.db.Where("slug = ?", slug).Take(&point).Error
		if err == nil {
			f.Title = point.Title
			if point.CEFRLevel != "" {
				f.Level = point.CEFRLevel
			}
			var common []struct {
				Wrong string `json:"wrong"`
				Right string `json:"right"`
			}
			_ = json.Unmarshal(point.CommonMistakesJSON, &common)
			var b strings.Builder
			if point.Formula != "" {
				fmt.Fprintf(&b, "Form: %s\n", point.Formula)
			}
			for _, mistake := range common {
				fmt.Fprintf(&b, "Typical mistake: %q → %q\n", mistake.Wrong, mistake.Right)
			}
			f.Details = strings.TrimSpace(b.String())
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return f, fmt.Errorf("find grammar point: %w", err)
		}
	}
	if err := s.db.Where("user_id = ? AND skill_tag = ? AND status <> ?", userID, tag, models.MistakeFixed).
		Order("occurrences DESC, last_seen_at DESC").Limit(8).Find(&f.Mistakes).Error; err != nil {
		return f, fmt.Errorf("find mistakes: %w", err)
	}
	return f, nil
}

// fromBank picks published lesson exercises for the tag: answered wrong
// last time first, then never answered, then the rest.
func (s Service) fromBank(userID uint, tag string) ([]Item, error) {
	query := s.db.Model(&models.Exercise{}).
		Joins("JOIN lesson_activities a ON a.id = exercises.activity_id").
		Where("a.status = ? AND exercises.type <> ?", models.ContentPublished, models.ExerciseSpeak)
	switch {
	case tag == "stress" || strings.HasPrefix(tag, "stress:"):
		query = query.Where("exercises.skill_tag = 'stress' OR exercises.skill_tag LIKE 'stress:%'")
	case strings.HasPrefix(tag, "sound:") && !strings.Contains(strings.TrimPrefix(tag, "sound:"), "-"):
		// A single sound (from pronunciation scoring) matches every pair it is in.
		symbol := strings.TrimPrefix(tag, "sound:")
		query = query.Where("exercises.skill_tag = ? OR exercises.skill_tag LIKE ? OR exercises.skill_tag LIKE ?",
			tag, "sound:"+symbol+"-%", "sound:%-"+symbol)
	default:
		query = query.Where("exercises.skill_tag = ?", tag)
	}
	var exercises []models.Exercise
	if err := query.Select("exercises.*").Order("exercises.id").Limit(200).Find(&exercises).Error; err != nil {
		return nil, fmt.Errorf("find exercises: %w", err)
	}

	ids := make([]uint, 0, len(exercises))
	for _, exercise := range exercises {
		ids = append(ids, exercise.ID)
	}
	var answers []models.UserExerciseAnswer
	if len(ids) > 0 {
		if err := s.db.Where("user_id = ? AND exercise_id IN ?", userID, ids).Order("answered_at ASC").Find(&answers).Error; err != nil {
			return nil, fmt.Errorf("find answers: %w", err)
		}
	}
	latest := map[uint]bool{}
	answered := map[uint]bool{}
	for _, answer := range answers {
		latest[answer.ExerciseID] = answer.IsCorrect
		answered[answer.ExerciseID] = true
	}
	rank := func(exercise models.Exercise) int {
		switch {
		case answered[exercise.ID] && !latest[exercise.ID]:
			return 0
		case !answered[exercise.ID]:
			return 1
		default:
			return 2
		}
	}
	// Stable order by rank, rotating within a rank by day so repeats vary.
	offset := int(s.now().Unix()/86400) % max(1, len(exercises))
	rotated := append(append([]models.Exercise{}, exercises[offset:]...), exercises[:offset]...)
	items := make([]Item, 0, setSize)
	for r := 0; r <= 2 && len(items) < setSize; r++ {
		for _, exercise := range rotated {
			if rank(exercise) != r || len(items) == setSize {
				continue
			}
			id := exercise.ID
			items = append(items, Item{
				Type:             exercise.Type,
				Prompt:           exercise.Prompt,
				Payload:          json.RawMessage(exercise.PayloadJSON),
				Answer:           json.RawMessage(exercise.AnswerJSON),
				ExplanationVI:    exercise.ExplanationVI,
				SkillTag:         exercise.SkillTag,
				SourceExerciseID: &id,
			})
		}
	}
	return items, nil
}

// Get returns a learner's set (to resume it).
func (s Service) Get(userID, setID uint) (SetResponse, error) {
	set, err := s.find(s.db, userID, setID, false)
	if err != nil {
		return SetResponse{}, err
	}
	return s.response(set)
}

// Answer grades one item. Grammar mastery and the mistake log are updated
// like for lesson exercises.
func (s Service) Answer(userID, setID uint, req AnswerRequest) (AnswerResult, error) {
	now := s.now()
	var result AnswerResult
	err := s.db.Transaction(func(tx *gorm.DB) error {
		set, err := s.find(tx, userID, setID, true)
		if err != nil {
			return err
		}
		var items []Item
		_ = json.Unmarshal(set.ItemsJSON, &items)
		index := int(req.ExerciseID)
		if index < 1 || index > len(items) {
			return ErrItemNotFound
		}
		results := map[string]AnswerResult{}
		_ = json.Unmarshal(set.ResultsJSON, &results)
		key := fmt.Sprint(index)
		if _, done := results[key]; done {
			return ErrAlreadyAnswer
		}
		item := items[index-1]

		graded, err := lessonactivity.Grade(item.Type, item.Payload, item.Answer, req.Answer)
		if err != nil {
			return err
		}
		result = AnswerResult{
			ExerciseID:    req.ExerciseID,
			Attempt:       1,
			IsCorrect:     graded.IsCorrect,
			ScoreRatio:    graded.ScoreRatio,
			Answer:        req.Answer,
			CorrectAnswer: graded.CorrectAnswer,
			ExplanationVI: item.ExplanationVI,
			AnsweredAt:    now,
		}
		results[key] = result
		resultsJSON, _ := json.Marshal(results)
		if err := tx.Model(&set).Updates(map[string]any{"results_json": datatypes.JSON(resultsJSON), "updated_at": now}).Error; err != nil {
			return fmt.Errorf("save practice answer: %w", err)
		}
		if err := grammarprogress.Record(tx, userID, item.SkillTag, graded.ScoreRatio, graded.IsCorrect, now); err != nil {
			return err
		}
		exercise := models.Exercise{Type: item.Type, PayloadJSON: datatypes.JSON(item.Payload), AnswerJSON: datatypes.JSON(item.Answer), ExplanationVI: item.ExplanationVI, SkillTag: item.SkillTag}
		if mistake, ok := lessonactivity.ExerciseMistake(userID, exercise, "practice", req.Answer, lessonactivity.GradeOutcome{GradeResult: graded}, now); ok {
			return mistakes.Upsert(tx, mistake)
		}
		return nil
	})
	if err != nil {
		return AnswerResult{}, err
	}
	return result, nil
}

func (s Service) find(db *gorm.DB, userID, setID uint, lock bool) (models.PracticeSet, error) {
	var set models.PracticeSet
	query := db.Where("id = ? AND user_id = ?", setID, userID)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := query.Take(&set).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return set, ErrSetNotFound
	}
	if err != nil {
		return set, fmt.Errorf("find practice set: %w", err)
	}
	return set, nil
}

func (s Service) response(set models.PracticeSet) (SetResponse, error) {
	var items []Item
	_ = json.Unmarshal(set.ItemsJSON, &items)
	results := map[string]AnswerResult{}
	_ = json.Unmarshal(set.ResultsJSON, &results)

	titles := map[string]string{}
	if slug, ok := strings.CutPrefix(set.SkillTag, "grammar:"); ok {
		var point models.GrammarPoint
		if s.db.Where("slug = ?", slug).Take(&point).Error == nil {
			titles[slug] = point.TitleVI
			if titles[slug] == "" {
				titles[slug] = point.Title
			}
		}
	}

	response := SetResponse{
		ID:        set.ID,
		SkillTag:  set.SkillTag,
		Label:     mistakes.SkillLabel(set.SkillTag, titles),
		Source:    set.Source,
		Day:       set.Day,
		Exercises: make([]Exercise, 0, len(items)),
		Answers:   []AnswerResult{},
	}
	for index, item := range items {
		response.Exercises = append(response.Exercises, Exercise{
			ID:         uint(index + 1),
			Type:       item.Type,
			Prompt:     item.Prompt,
			Payload:    item.Payload,
			Points:     10,
			OrderIndex: index + 1,
		})
		if result, ok := results[fmt.Sprint(index+1)]; ok {
			response.Answers = append(response.Answers, result)
		}
	}
	return response, nil
}
