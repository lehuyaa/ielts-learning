package lessonactivity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/aiusage"
	"ielts-learning/backend/internal/shared/ai"
)

// How an answer was graded.
const (
	GradedByMatch = "MATCH"
	GradedByAI    = "AI"
)

// Why AI grading was skipped although the exercise allows it. The answer then
// keeps the matching result (incorrect).
const (
	AIStatusLimitReached = "LIMIT_REACHED"
	AIStatusUnavailable  = "UNAVAILABLE"
)

// AIFeedback is the AI's explanation for a short answer (spec §6.4).
type AIFeedback struct {
	Suggested     string `json:"suggested"`
	ExplanationVI string `json:"explanationVi"`
}

// GradeOutcome is a GradeResult plus how it was reached.
type GradeOutcome struct {
	GradeResult
	GradedBy   string
	AIFeedback *AIFeedback
	AIStatus   string
}

// jsonCompleter is the part of ai.Client the grader needs (fake in tests).
type jsonCompleter interface {
	Configured() bool
	CompleteJSON(ctx context.Context, messages []ai.Message, schema ai.Schema, out any) error
}

// AIGrader grades TRANSFORM / TRANSLATE_VI_EN answers with AI when normalised
// matching fails (spec §6.4), with a per-user daily limit (§6.7) and a shared
// cache per exercise + normalised answer.
type AIGrader struct {
	db         *gorm.DB
	client     jsonCompleter
	model      string
	dailyLimit int
	// writingLimit caps writing gradings (incl. rewrite checks) per user per day.
	writingLimit int
	timeout      time.Duration
	now          func() time.Time
}

// WithWritingLimit sets the daily limit of AI writing gradings (spec §6.7).
func (g *AIGrader) WithWritingLimit(limit int) *AIGrader {
	g.writingLimit = limit
	return g
}

// NewAIGrader returns a grader. A nil or unconfigured client disables AI.
func NewAIGrader(db *gorm.DB, client *ai.Client, model string, dailyLimit int) *AIGrader {
	grader := &AIGrader{db: db, model: model, dailyLimit: dailyLimit, writingLimit: 30, timeout: 25 * time.Second, now: time.Now}
	if client != nil {
		grader.client = client
	}
	return grader
}

type shortAnswerVerdict struct {
	IsCorrect     bool    `json:"isCorrect"`
	ScoreRatio    float64 `json:"scoreRatio"`
	Suggested     string  `json:"suggested"`
	ExplanationVI string  `json:"explanationVi"`
}

var shortAnswerSchema = ai.Schema{
	Name: "short_answer_grade",
	Schema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"isCorrect":     map[string]any{"type": "boolean"},
			"scoreRatio":    map[string]any{"type": "number", "description": "1 correct, 0.5 meaning right with one small error, 0 wrong"},
			"suggested":     map[string]any{"type": "string", "description": "the learner's answer corrected with minimal changes"},
			"explanationVi": map[string]any{"type": "string", "description": "one or two short sentences in Vietnamese"},
		},
		"required":             []string{"isCorrect", "scoreRatio", "suggested", "explanationVi"},
		"additionalProperties": false,
	},
}

// Resolve completes a matching result: when result.NeedsAI it asks the AI
// (or the cache); otherwise it returns the matching result unchanged.
// It never fails — any AI problem keeps the matching result.
func (g *AIGrader) Resolve(ctx context.Context, userID uint, exercise models.Exercise, submittedJSON []byte, result GradeResult) GradeOutcome {
	outcome := GradeOutcome{GradeResult: result, GradedBy: GradedByMatch}
	if !result.NeedsAI {
		return outcome
	}
	if g == nil || g.client == nil || !g.client.Configured() {
		outcome.AIStatus = AIStatusUnavailable
		return outcome
	}

	var submitted Submitted
	if err := json.Unmarshal(submittedJSON, &submitted); err != nil || submitted.Text == nil {
		return outcome
	}
	normalized := normalizeText(*submitted.Text)
	hash := answerHash(normalized)

	if verdict, ok := g.cached(exercise.ID, hash); ok {
		return applyVerdict(outcome, verdict)
	}

	allowed, err := g.takeQuota(userID, models.AIUsageShortAnswer, g.dailyLimit)
	if err != nil {
		log.Printf("ai grading: usage counter: %v", err)
		outcome.AIStatus = AIStatusUnavailable
		return outcome
	}
	if !allowed {
		outcome.AIStatus = AIStatusLimitReached
		return outcome
	}

	callCtx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()
	var verdict shortAnswerVerdict
	if err := g.client.CompleteJSON(callCtx, g.prompt(exercise, *submitted.Text), shortAnswerSchema, &verdict); err != nil {
		log.Printf("ai grading: exercise %d: %v", exercise.ID, err)
		outcome.AIStatus = AIStatusUnavailable
		return outcome
	}
	verdict = sanitizeVerdict(verdict)
	g.store(exercise.ID, hash, normalized, verdict)
	return applyVerdict(outcome, verdict)
}

func applyVerdict(outcome GradeOutcome, verdict shortAnswerVerdict) GradeOutcome {
	outcome.GradedBy = GradedByAI
	outcome.IsCorrect = verdict.IsCorrect
	outcome.ScoreRatio = verdict.ScoreRatio
	outcome.NeedsAI = false
	outcome.AIFeedback = &AIFeedback{Suggested: verdict.Suggested, ExplanationVI: verdict.ExplanationVI}
	return outcome
}

// sanitizeVerdict keeps scores on the 0 / 0.5 / 1 scale and consistent with isCorrect.
func sanitizeVerdict(verdict shortAnswerVerdict) shortAnswerVerdict {
	verdict.Suggested = strings.TrimSpace(verdict.Suggested)
	verdict.ExplanationVI = strings.TrimSpace(verdict.ExplanationVI)
	if verdict.IsCorrect {
		verdict.ScoreRatio = 1
		return verdict
	}
	if math.IsNaN(verdict.ScoreRatio) || verdict.ScoreRatio < 0.5 {
		verdict.ScoreRatio = 0
	} else {
		verdict.ScoreRatio = 0.5
	}
	return verdict
}

func answerHash(normalized string) string {
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

func (g *AIGrader) cached(exerciseID uint, hash string) (shortAnswerVerdict, bool) {
	var row models.AIShortAnswerGrading
	err := g.db.Where("exercise_id = ? AND answer_hash = ?", exerciseID, hash).Take(&row).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("ai grading: cache lookup: %v", err)
		}
		return shortAnswerVerdict{}, false
	}
	var verdict shortAnswerVerdict
	if err := json.Unmarshal(row.ResultJSON, &verdict); err != nil {
		return shortAnswerVerdict{}, false
	}
	return verdict, true
}

func (g *AIGrader) store(exerciseID uint, hash string, normalized string, verdict shortAnswerVerdict) {
	payload, _ := json.Marshal(verdict)
	row := models.AIShortAnswerGrading{
		ExerciseID: exerciseID,
		AnswerHash: hash,
		AnswerText: normalized,
		ResultJSON: datatypes.JSON(payload),
		Model:      g.model,
	}
	if err := g.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
		log.Printf("ai grading: cache store: %v", err)
	}
}

// takeQuota counts one AI call of a kind for today and reports whether it is allowed.
func (g *AIGrader) takeQuota(userID uint, kind string, limit int) (bool, error) {
	return aiusage.Take(g.db, userID, kind, limit, g.now())
}

func (g *AIGrader) prompt(exercise models.Exercise, answer string) []ai.Message {
	var payload struct {
		Sentence    string `json:"sentence"`
		Instruction string `json:"instruction"`
		SourceVI    string `json:"sourceVi"`
	}
	_ = json.Unmarshal(exercise.PayloadJSON, &payload)
	var key answerKey
	_ = json.Unmarshal(exercise.AnswerJSON, &key)
	references := append(append([]string{}, key.Accepted...), key.Reference...)

	level := g.levelOf(exercise.ID)
	var task string
	if exercise.Type == models.ExerciseTransform {
		task = fmt.Sprintf("Rewrite the sentence %q following the instruction %s (%s).",
			payload.Sentence, payload.Instruction, transformInstructionEN(payload.Instruction))
	} else {
		task = fmt.Sprintf("Translate this Vietnamese sentence into English: %q", payload.SourceVI)
	}

	system := "You grade one short answer from a Vietnamese adult learning English at CEFR " + level + ".\n" +
		"Accept any answer that is natural, grammatically correct English and does what the task asks, " +
		"even when it differs from the reference answers. Ignore capitalisation, final punctuation and contractions.\n" +
		"isCorrect=false with scoreRatio 0.5 when the meaning is right but there is exactly one small grammar or spelling error; " +
		"scoreRatio 0 when the task is not done, the meaning is wrong, or there are several errors.\n" +
		"suggested: the learner's answer corrected with as few changes as possible (equal to the answer when it is correct).\n" +
		"explanationVi: one or two short sentences in Vietnamese saying what is wrong and why (when correct, a short confirmation). " +
		"Quote English words in single quotes."
	user := fmt.Sprintf("Task: %s\nReference answers: %s\nLearner answer: %q",
		task, strings.Join(quoteAll(references), " | "), strings.TrimSpace(answer))
	return []ai.Message{{Role: "system", Content: system}, {Role: "user", Content: user}}
}

func (g *AIGrader) levelOf(exerciseID uint) string {
	var level *string
	g.db.Raw(`SELECT bl.cefr_level FROM exercises e
		JOIN lesson_activities a ON a.id = e.activity_id
		JOIN lessons l ON l.id = a.lesson_id
		JOIN topics t ON t.id = l.topic_id
		JOIN band_levels bl ON bl.id = t.band_level_id
		WHERE e.id = ?`, exerciseID).Scan(&level)
	if level == nil || *level == "" {
		return "A1"
	}
	return *level
}

func transformInstructionEN(instruction string) string {
	switch instruction {
	case "NEGATIVE":
		return "make it negative"
	case "QUESTION":
		return "make it a yes/no question"
	case "PAST":
		return "put it in the past simple"
	case "THIRD_PERSON":
		return "change the subject to he/she and adjust the verb"
	default:
		return "follow the instruction"
	}
}

func quoteAll(values []string) []string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, fmt.Sprintf("%q", value))
	}
	return quoted
}
