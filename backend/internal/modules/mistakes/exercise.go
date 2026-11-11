package mistakes

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"ielts-learning/backend/internal/models"
)

var blankRun = regexp.MustCompile(`_{2,}(\s*\([^)]*\))?`)

type choiceOption struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Word  string `json:"word"`
	IPA   string `json:"ipa"`
	Label string `json:"label"`
}

type exercisePayload struct {
	Options   json.RawMessage `json:"options"`
	Tokens    []string        `json:"tokens"`
	Sentence  string          `json:"sentence"`
	Statement string          `json:"statement"`
	Syllables []string        `json:"syllables"`
}

type exerciseKey struct {
	Correct     json.RawMessage     `json:"correct"`
	Value       *bool               `json:"value"`
	Accepted    []string            `json:"accepted"`
	Reference   []string            `json:"reference"`
	WrongIndex  *int                `json:"wrongIndex"`
	Blanks      map[string][]string `json:"blanks"`
	StressIndex *int                `json:"stressIndex"`
}

type exerciseAnswer struct {
	OptionID    *string           `json:"optionId"`
	Value       *bool             `json:"value"`
	Text        *string           `json:"text"`
	Tokens      []string          `json:"tokens"`
	Index       *int              `json:"index"`
	Blanks      map[string]string `json:"blanks"`
	StressIndex *int              `json:"stressIndex"`
}

// DescribeExercise turns a wrong exercise answer into a (wrong → right) pair
// for the mistake log. ok is false when the type has nothing to show
// (SPEAK is scored by pronunciation instead).
func DescribeExercise(exerciseType models.ExerciseType, payloadJSON, answerJSON, submittedJSON []byte) (wrong string, correct string, ok bool) {
	var payload exercisePayload
	var key exerciseKey
	var submitted exerciseAnswer
	if json.Unmarshal(payloadJSON, &payload) != nil || json.Unmarshal(answerJSON, &key) != nil || json.Unmarshal(submittedJSON, &submitted) != nil {
		return "", "", false
	}

	switch exerciseType {
	case models.ExerciseMultipleChoice, models.ExerciseListenPickImage, models.ExerciseMinimalPair,
		models.ExerciseIPAToWord, models.ExerciseWordToIPA, models.ExerciseEndingSound:
		var correctID string
		if json.Unmarshal(key.Correct, &correctID) != nil || submitted.OptionID == nil {
			return "", "", false
		}
		labels := optionLabels(payload.Options)
		wrong, correct = labels[*submitted.OptionID], labels[correctID]
		if wrong == "" || correct == "" {
			return "", "", false
		}
		if blankRun.MatchString(payload.Sentence) {
			return fillBlank(payload.Sentence, wrong), fillBlank(payload.Sentence, correct), true
		}
		return wrong, correct, true
	case models.ExerciseTrueFalse:
		if submitted.Value == nil || key.Value == nil || payload.Statement == "" {
			return "", "", false
		}
		return payload.Statement + " → " + yesNo(*submitted.Value), payload.Statement + " → " + yesNo(*key.Value), true
	case models.ExerciseFillBlank:
		if submitted.Text == nil || len(key.Accepted) == 0 {
			return "", "", false
		}
		if blankRun.MatchString(payload.Sentence) {
			return fillBlank(payload.Sentence, *submitted.Text), fillBlank(payload.Sentence, key.Accepted[0]), true
		}
		return *submitted.Text, key.Accepted[0], true
	case models.ExerciseTransform, models.ExerciseTranslateViEn, models.ExerciseDictation, models.ExerciseListenNumbers:
		reference := append(append([]string{}, key.Accepted...), key.Reference...)
		if submitted.Text == nil || len(reference) == 0 {
			return "", "", false
		}
		return *submitted.Text, reference[0], true
	case models.ExerciseReorder:
		if len(submitted.Tokens) == 0 || len(key.Accepted) == 0 {
			return "", "", false
		}
		return joinTokens(submitted.Tokens), key.Accepted[0], true
	case models.ExerciseErrorCorrection:
		if key.WrongIndex == nil || len(key.Accepted) == 0 || *key.WrongIndex >= len(payload.Tokens) {
			return "", "", false
		}
		fixed := append([]string{}, payload.Tokens...)
		fixed[*key.WrongIndex] = key.Accepted[0]
		return joinTokens(payload.Tokens), joinTokens(fixed), true
	case models.ExerciseListenGapFill:
		if len(key.Blanks) == 0 {
			return "", "", false
		}
		names := make([]string, 0, len(key.Blanks))
		for name := range key.Blanks {
			names = append(names, name)
		}
		sort.Strings(names)
		var wrongParts, correctParts []string
		for _, name := range names {
			if len(key.Blanks[name]) == 0 {
				continue
			}
			wrongParts = append(wrongParts, "("+name+") "+submitted.Blanks[name])
			correctParts = append(correctParts, "("+name+") "+key.Blanks[name][0])
		}
		return strings.Join(wrongParts, ", "), strings.Join(correctParts, ", "), true
	case models.ExerciseWordStress:
		if submitted.StressIndex == nil || key.StressIndex == nil || len(payload.Syllables) == 0 {
			return "", "", false
		}
		return stressed(payload.Syllables, *submitted.StressIndex), stressed(payload.Syllables, *key.StressIndex), true
	default:
		return "", "", false
	}
}

func optionLabels(raw json.RawMessage) map[string]string {
	labels := map[string]string{}
	var plain []string
	if json.Unmarshal(raw, &plain) == nil {
		for _, option := range plain {
			labels[option] = option
		}
		return labels
	}
	var options []choiceOption
	_ = json.Unmarshal(raw, &options)
	for _, option := range options {
		for _, label := range []string{option.Text, option.Word, option.Label, option.IPA, option.ID} {
			if strings.TrimSpace(label) != "" {
				labels[option.ID] = label
				break
			}
		}
	}
	return labels
}

func fillBlank(sentence, value string) string {
	done := false
	return blankRun.ReplaceAllStringFunc(sentence, func(string) string {
		if done {
			return "___"
		}
		done = true
		return strings.TrimSpace(value)
	})
}

func yesNo(value bool) string {
	if value {
		return "Đúng"
	}
	return "Sai"
}

func stressed(syllables []string, index int) string {
	parts := make([]string, len(syllables))
	for i, syllable := range syllables {
		if i == index {
			parts[i] = strings.ToUpper(syllable)
		} else {
			parts[i] = syllable
		}
	}
	return strings.Join(parts, "-")
}

func joinTokens(tokens []string) string {
	var b strings.Builder
	for i, token := range tokens {
		if i > 0 && !strings.ContainsAny(token, ".,!?;:") {
			b.WriteString(" ")
		}
		b.WriteString(token)
	}
	return b.String()
}
