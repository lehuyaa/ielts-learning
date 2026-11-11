package lessonactivity

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"unicode/utf16"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/shared/ai"
)

var (
	// ErrAILimitReached means the learner used today's AI gradings (spec §6.7).
	ErrAILimitReached = errors.New("ai limit reached")
	// ErrAIUnavailable means AI is not configured or the call failed.
	ErrAIUnavailable = errors.New("ai unavailable")
)

// WritingContent is a WRITING activity's content (docs/foundation-spec.md §4.1).
type WritingContent struct {
	PromptVI         string   `json:"promptVi"`
	PromptEN         string   `json:"promptEn"`
	Frame            []string `json:"frame"`
	MinSentences     int      `json:"minSentences"`
	MaxSentences     int      `json:"maxSentences"`
	TargetGrammar    []string `json:"targetGrammar"`
	TargetVocabulary []string `json:"targetVocabulary"`
	ModelAnswer      string   `json:"modelAnswer"`
}

// WritingScores are 1-10 per criterion.
type WritingScores struct {
	TaskAchievement int `json:"taskAchievement"`
	Grammar         int `json:"grammar"`
	Vocabulary      int `json:"vocabulary"`
	Coherence       int `json:"coherence"`
}

// WritingMistake marks one error. Start/End are UTF-16 offsets into the text
// that was graded (-1 when the AI's quote could not be found).
type WritingMistake struct {
	Start         int    `json:"start"`
	End           int    `json:"end"`
	Wrong         string `json:"wrong"`
	Correct       string `json:"correct"`
	SkillTag      string `json:"skillTag"`
	ExplanationVI string `json:"explanationVi"`
}

type WritingTargets struct {
	Grammar    []string `json:"grammar"`
	Vocabulary []string `json:"vocabulary"`
}

// WritingReport is the AI grading result (spec §6.3).
type WritingReport struct {
	Scores          WritingScores    `json:"scores"`
	Overall         int              `json:"overall"`
	CorrectedText   string           `json:"correctedText"`
	Mistakes        []WritingMistake `json:"mistakes"`
	UsedTargets     WritingTargets   `json:"usedTargets"`
	FeedbackVI      string           `json:"feedbackVi"`
	RewriteRequired bool             `json:"rewriteRequired"`
}

func intField() map[string]any {
	return map[string]any{"type": "integer", "description": "1-10"}
}

var writingSchema = ai.Schema{
	Name: "writing_grade",
	Schema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"scores": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"taskAchievement": intField(), "grammar": intField(), "vocabulary": intField(), "coherence": intField(),
				},
				"required":             []string{"taskAchievement", "grammar", "vocabulary", "coherence"},
				"additionalProperties": false,
			},
			"correctedText": map[string]any{"type": "string"},
			"mistakes": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"wrong":         map[string]any{"type": "string", "description": "exact words copied from the learner text"},
						"correct":       map[string]any{"type": "string"},
						"skillTag":      map[string]any{"type": "string"},
						"explanationVi": map[string]any{"type": "string"},
					},
					"required":             []string{"wrong", "correct", "skillTag", "explanationVi"},
					"additionalProperties": false,
				},
			},
			"usedTargets": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"grammar":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"vocabulary": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
				},
				"required":             []string{"grammar", "vocabulary"},
				"additionalProperties": false,
			},
			"feedbackVi": map[string]any{"type": "string"},
		},
		"required":             []string{"scores", "correctedText", "mistakes", "usedTargets", "feedbackVi"},
		"additionalProperties": false,
	},
}

const writingRules = `Rules:
- Only mark real errors: grammar, spelling, wrong word, missing article, capital letters for "I", names and sentence starts. Do not rewrite correct sentences into "better" ones; at this level simple correct English is good English.
- "wrong" must be copied exactly from the learner text (the smallest span that shows the error, usually 1-4 words) and appear in text order. "correct" is what should replace exactly that span.
- skillTag: "grammar:<slug>" using one of these slugs when it fits: g1-to-be, g2-nouns-articles, g3-present-simple, g4-present-continuous, g5-there-is-prepositions, g6-adjectives-comparison, g7-can-must-should, g8-past-simple, g9-future; otherwise "spelling", "vocab:word-choice", "capitalization" or "punctuation".
- explanationVi: one short sentence in Vietnamese, quoting English words in single quotes.
- correctedText: the learner text with only the mistakes fixed.
- scores 1-10: taskAchievement (answers the prompt, enough sentences), grammar, vocabulary (uses the target words), coherence (clear order, linking words).
- usedTargets: the target grammar slugs and vocabulary words the learner used correctly.
- feedbackVi: 1-2 encouraging sentences in Vietnamese: one strength and the most important thing to fix.`

// GradeWriting asks the AI to grade a text for a WRITING activity.
func (g *AIGrader) GradeWriting(ctx context.Context, userID uint, content WritingContent, level string, text string) (WritingReport, error) {
	user := fmt.Sprintf("Task (Vietnamese): %s\nTask (English): %s\nSentence frames: %s\nLength: %d-%d sentences\nTarget grammar: %s\nTarget words: %s\n\nLearner text:\n%s",
		content.PromptVI, content.PromptEN, strings.Join(content.Frame, " | "), content.MinSentences, content.MaxSentences,
		strings.Join(content.TargetGrammar, ", "), strings.Join(targetWords(content.TargetVocabulary), ", "), text)
	system := "You grade short texts written by a Vietnamese adult learning English at CEFR " + level + ".\n" + writingRules
	return g.callWriting(ctx, userID, system, user, text)
}

// CheckRewrite asks the AI which mistakes remain in a corrected text.
func (g *AIGrader) CheckRewrite(ctx context.Context, userID uint, content WritingContent, level string, original WritingReport, rewrite string) (WritingReport, error) {
	var listed []string
	for _, mistake := range original.Mistakes {
		listed = append(listed, fmt.Sprintf("%q -> %q", mistake.Wrong, mistake.Correct))
	}
	user := fmt.Sprintf("Task: %s\nMistakes found in the first version: %s\n\nRewritten text:\n%s",
		content.PromptEN, strings.Join(listed, "; "), rewrite)
	system := "A Vietnamese adult learning English at CEFR " + level + " rewrote their text after feedback. " +
		"Report only the mistakes that are still in the rewritten text (an empty list when it is correct); new correct sentences are fine.\n" + writingRules
	return g.callWriting(ctx, userID, system, user, rewrite)
}

func (g *AIGrader) callWriting(ctx context.Context, userID uint, system string, user string, text string) (WritingReport, error) {
	if g == nil || g.client == nil || !g.client.Configured() {
		return WritingReport{}, ErrAIUnavailable
	}
	allowed, err := g.takeQuota(userID, models.AIUsageWriting, g.writingLimit)
	if err != nil {
		return WritingReport{}, err
	}
	if !allowed {
		return WritingReport{}, ErrAILimitReached
	}

	callCtx, cancel := context.WithTimeout(ctx, 2*g.timeout)
	defer cancel()
	var report WritingReport
	messages := []ai.Message{{Role: "system", Content: system}, {Role: "user", Content: user}}
	if err := g.client.CompleteJSON(callCtx, messages, writingSchema, &report); err != nil {
		log.Printf("writing grading: %v", err)
		return WritingReport{}, fmt.Errorf("%w: %v", ErrAIUnavailable, err)
	}
	return finishReport(report, text), nil
}

// finishReport clamps scores, computes overall (10 × average) and locates
// each mistake in the text.
func finishReport(report WritingReport, text string) WritingReport {
	clamp := func(value int) int { return max(1, min(10, value)) }
	report.Scores = WritingScores{
		TaskAchievement: clamp(report.Scores.TaskAchievement),
		Grammar:         clamp(report.Scores.Grammar),
		Vocabulary:      clamp(report.Scores.Vocabulary),
		Coherence:       clamp(report.Scores.Coherence),
	}
	sum := report.Scores.TaskAchievement + report.Scores.Grammar + report.Scores.Vocabulary + report.Scores.Coherence
	report.Overall = int(math.Round(float64(sum) * 10 / 4))

	mistakes := make([]WritingMistake, 0, len(report.Mistakes))
	cursor := 0
	for _, mistake := range report.Mistakes {
		mistake.Wrong = strings.TrimSpace(mistake.Wrong)
		mistake.Correct = strings.TrimSpace(mistake.Correct)
		if mistake.Wrong == "" || strings.EqualFold(mistake.Wrong, mistake.Correct) {
			continue
		}
		start, end := locate(text, mistake.Wrong, cursor)
		if start >= 0 {
			cursor = end
			mistake.Start, mistake.End = utf16Offset(text, start), utf16Offset(text, end)
		} else {
			mistake.Start, mistake.End = -1, -1
		}
		mistakes = append(mistakes, mistake)
	}
	report.Mistakes = mistakes
	if report.UsedTargets.Grammar == nil {
		report.UsedTargets.Grammar = []string{}
	}
	if report.UsedTargets.Vocabulary == nil {
		report.UsedTargets.Vocabulary = []string{}
	}
	report.RewriteRequired = len(mistakes) > 0
	return report
}

// locate finds needle in text (byte offsets), preferring matches after from,
// then anywhere, then ignoring case.
func locate(text string, needle string, from int) (int, int) {
	if index := strings.Index(text[from:], needle); index >= 0 {
		return from + index, from + index + len(needle)
	}
	if index := strings.Index(text, needle); index >= 0 {
		return index, index + len(needle)
	}
	lower := strings.ToLower(text)
	if index := strings.Index(lower, strings.ToLower(needle)); index >= 0 && len(lower) == len(text) {
		return index, index + len(needle)
	}
	return -1, -1
}

// utf16Offset converts a byte offset into a JavaScript string index.
func utf16Offset(text string, byteOffset int) int {
	return len(utf16.Encode([]rune(text[:byteOffset])))
}

func targetWords(slugs []string) []string {
	words := make([]string, 0, len(slugs))
	for _, slug := range slugs {
		words = append(words, strings.ReplaceAll(strings.TrimPrefix(slug, "a1-"), "-", " "))
	}
	return words
}
