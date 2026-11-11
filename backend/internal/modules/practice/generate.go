package practice

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/contentimport"
	"ielts-learning/backend/internal/shared/ai"
)

// generatedTypes are the item types the AI may write: all graded by
// matching, so practice never needs a second AI call to grade.
var generatedTypes = []string{
	string(models.ExerciseMultipleChoice),
	string(models.ExerciseFillBlank),
	string(models.ExerciseReorder),
	string(models.ExerciseErrorCorrection),
}

var (
	str   = map[string]any{"type": "string"}
	words = regexp.MustCompile(`[A-Za-z0-9']+`)
)

func strictObject(properties map[string]any) map[string]any {
	required := make([]string, 0, len(properties))
	for key := range properties {
		required = append(required, key)
	}
	sort.Strings(required)
	return map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}
}

var generateSchema = ai.Schema{
	Name: "practice_generation",
	Schema: strictObject(map[string]any{
		"exercises": map[string]any{"type": "array", "items": strictObject(map[string]any{
			"type":          map[string]any{"type": "string", "enum": generatedTypes},
			"prompt":        str,
			"payloadJson":   str,
			"answerJson":    str,
			"explanationVi": str,
		})},
	}),
}

const generateSystemPrompt = `You write short practice exercises for a Vietnamese adult learning English at CEFR %s.
The practice targets ONE weakness of this learner (below). Write %d different items that practise exactly that point, mixing the types.

Rules:
- Use only everyday words a %s learner knows; every English sentence at most %d words.
- Each item has exactly one correct answer. Wrong options must be clearly wrong for this point (often the learner's own mistake).
- prompt: a short Vietnamese instruction. explanationVi: one or two short Vietnamese sentences on why the answer is right; quote English in single quotes.
- Vary the sentences; do not reuse the learner's mistakes word for word.

payloadJson / answerJson are JSON strings:
- MULTIPLE_CHOICE: {"sentence":"She ___ a teacher.","options":[{"id":"a","text":"is"},{"id":"b","text":"are"},{"id":"c","text":"am"}]} / {"correct":"a"}  (3-4 options; "sentence" with one ___ is optional)
- FILL_BLANK: {"sentence":"They ___ (be) from Hanoi."} / {"accepted":["are"]}  (exactly one ___, list every correct form incl. contractions)
- REORDER: {"tokens":["from","I","am","Vietnam","."]} (the words of the answer, shuffled) / {"accepted":["I am from Vietnam."]}
- ERROR_CORRECTION: {"tokens":["He","are","my","friend","."]} / {"wrongIndex":1,"accepted":["is"]}  (exactly one wrong token)`

type generatedItem struct {
	Type          string `json:"type"`
	Prompt        string `json:"prompt"`
	PayloadJSON   string `json:"payloadJson"`
	AnswerJSON    string `json:"answerJson"`
	ExplanationVI string `json:"explanationVi"`
}

// focus describes the weakness for the prompt.
type focus struct {
	Tag      string
	Level    string
	Title    string
	Details  string
	Mistakes []models.LearnerMistake
}

func (s Service) generateWithAI(ctx context.Context, f focus) ([]Item, error) {
	maxWords := 12
	if f.Level != "A1" {
		maxWords = 15
	}
	var user strings.Builder
	fmt.Fprintf(&user, "Weakness: %s (skill tag %s)\n", f.Title, f.Tag)
	if f.Details != "" {
		fmt.Fprintf(&user, "%s\n", f.Details)
	}
	if len(f.Mistakes) > 0 {
		user.WriteString("The learner's recent mistakes (wrong → right):\n")
		for _, mistake := range f.Mistakes {
			fmt.Fprintf(&user, "- %q → %q\n", mistake.WrongText, mistake.CorrectText)
		}
	}
	messages := []ai.Message{
		{Role: "system", Content: fmt.Sprintf(generateSystemPrompt, f.Level, requestCount, f.Level, maxWords)},
		{Role: "user", Content: user.String()},
	}

	callCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var out struct {
		Exercises []generatedItem `json:"exercises"`
	}
	if err := s.ai.CompleteJSON(callCtx, messages, generateSchema, &out); err != nil {
		return nil, err
	}
	return validItems(out.Exercises, f.Tag, maxWords), nil
}

// validItems keeps the items that pass the importer's checks (answer key
// grades as correct), the word limit and the type-specific checks, without
// duplicates, up to setSize.
func validItems(generated []generatedItem, tag string, maxWords int) []Item {
	items := make([]Item, 0, setSize)
	seen := map[string]bool{}
	for _, raw := range generated {
		input := contentimport.ExerciseInput{
			Type:          raw.Type,
			Prompt:        strings.TrimSpace(raw.Prompt),
			Payload:       json.RawMessage(raw.PayloadJSON),
			Answer:        json.RawMessage(raw.AnswerJSON),
			ExplanationVI: strings.TrimSpace(raw.ExplanationVI),
			SkillTag:      tag,
		}
		if !json.Valid(input.Payload) || !json.Valid(input.Answer) {
			continue
		}
		if errs := contentimport.ValidateExercise(input); len(errs) > 0 {
			continue
		}
		if !checkItem(raw.Type, input.Payload, input.Answer, maxWords) {
			continue
		}
		key := raw.Type + string(compact(input.Payload))
		if seen[key] {
			continue
		}
		seen[key] = true
		items = append(items, Item{
			Type:          models.ExerciseType(raw.Type),
			Prompt:        input.Prompt,
			Payload:       compact(input.Payload),
			Answer:        compact(input.Answer),
			ExplanationVI: input.ExplanationVI,
			SkillTag:      tag,
		})
		if len(items) == setSize {
			break
		}
	}
	return items
}

type itemPayload struct {
	Sentence string   `json:"sentence"`
	Tokens   []string `json:"tokens"`
	Options  []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"options"`
}

type itemAnswer struct {
	Correct    string   `json:"correct"`
	Accepted   []string `json:"accepted"`
	WrongIndex *int     `json:"wrongIndex"`
}

func checkItem(itemType string, payloadJSON, answerJSON []byte, maxWords int) bool {
	var payload itemPayload
	var answer itemAnswer
	if json.Unmarshal(payloadJSON, &payload) != nil || json.Unmarshal(answerJSON, &answer) != nil {
		return false
	}
	if countWords(payload.Sentence) > maxWords || countWords(strings.Join(payload.Tokens, " ")) > maxWords {
		return false
	}
	switch models.ExerciseType(itemType) {
	case models.ExerciseMultipleChoice:
		if len(payload.Options) < 2 || len(payload.Options) > 5 {
			return false
		}
		texts := map[string]bool{}
		for _, option := range payload.Options {
			text := strings.ToLower(strings.TrimSpace(option.Text))
			if text == "" || texts[text] || countWords(text) > maxWords {
				return false
			}
			texts[text] = true
		}
	case models.ExerciseFillBlank:
		if strings.Count(payload.Sentence, "___") != 1 {
			return false
		}
	case models.ExerciseReorder:
		if len(answer.Accepted) == 0 || len(payload.Tokens) < 3 {
			return false
		}
		return sameWords(strings.Join(payload.Tokens, " "), answer.Accepted[0]) && !sameOrder(payload.Tokens, answer.Accepted[0])
	case models.ExerciseErrorCorrection:
		if answer.WrongIndex == nil || len(answer.Accepted) == 0 || *answer.WrongIndex >= len(payload.Tokens) {
			return false
		}
		return !strings.EqualFold(payload.Tokens[*answer.WrongIndex], answer.Accepted[0])
	}
	return true
}

func countWords(text string) int {
	return len(words.FindAllString(text, -1))
}

// sameWords: the shuffled tokens are exactly the words of the answer.
func sameWords(a, b string) bool {
	left, right := words.FindAllString(strings.ToLower(a), -1), words.FindAllString(strings.ToLower(b), -1)
	sort.Strings(left)
	sort.Strings(right)
	return strings.Join(left, " ") == strings.Join(right, " ")
}

// sameOrder: tokens already in the answer's order would make the item trivial.
func sameOrder(tokens []string, answer string) bool {
	return strings.Join(words.FindAllString(strings.ToLower(strings.Join(tokens, " ")), -1), " ") ==
		strings.Join(words.FindAllString(strings.ToLower(answer), -1), " ")
}

func compact(raw json.RawMessage) json.RawMessage {
	var value any
	if json.Unmarshal(raw, &value) != nil {
		return raw
	}
	out, _ := json.Marshal(value)
	return out
}
