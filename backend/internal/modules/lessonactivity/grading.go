package lessonactivity

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"ielts-learning/backend/internal/models"
)

// ErrInvalidAnswerShape means the submitted answer does not match what the
// exercise type expects (e.g. text for a multiple-choice item).
var ErrInvalidAnswerShape = errors.New("invalid answer shape")

// ErrInvalidAnswerKey means the stored exercise data is malformed.
var ErrInvalidAnswerKey = errors.New("invalid exercise answer key")

// correctThreshold is the score ratio from which a partially graded item
// (dictation, gap fill) counts as correct.
const correctThreshold = 0.9

// GradeResult is the outcome of grading one answer.
type GradeResult struct {
	IsCorrect  bool
	ScoreRatio float64
	// CorrectAnswer is safe to show after the learner has answered.
	CorrectAnswer any
	// NeedsAI is true when matching failed but the key allows AI grading
	// (TRANSFORM / TRANSLATE_VI_EN, Task 55).
	NeedsAI bool
	// PronunciationAttempt is set for a SPEAK answer scored by Azure
	// ({"attemptId"}); the service grades it from the stored attempt (Task 71).
	PronunciationAttempt uint
}

// Submitted is the union of answer shapes from docs/foundation-spec.md §5.1.
type Submitted struct {
	OptionID    *string           `json:"optionId"`
	Value       *bool             `json:"value"`
	Text        *string           `json:"text"`
	Tokens      []string          `json:"tokens"`
	Index       *int              `json:"index"`
	Blanks      map[string]string `json:"blanks"`
	StressIndex *int              `json:"stressIndex"`
	SelfCheck   *bool             `json:"selfCheck"`
	AttemptID   *uint             `json:"attemptId"`
}

type answerKey struct {
	Correct     json.RawMessage     `json:"correct"`
	Value       *bool               `json:"value"`
	Accepted    []string            `json:"accepted"`
	Reference   []string            `json:"reference"`
	WrongIndex  *int                `json:"wrongIndex"`
	Blanks      map[string][]string `json:"blanks"`
	StressIndex *int                `json:"stressIndex"`
	AIFallback  bool                `json:"aiFallback"`
	AIGraded    bool                `json:"aiGraded"`
}

type numbersPayload struct {
	Format string `json:"format"`
}

// Grade checks a submitted answer against an exercise's private answer key.
func Grade(exerciseType models.ExerciseType, payloadJSON []byte, answerJSON []byte, submittedJSON []byte) (GradeResult, error) {
	var key answerKey
	if err := json.Unmarshal(answerJSON, &key); err != nil {
		return GradeResult{}, fmt.Errorf("%w: %v", ErrInvalidAnswerKey, err)
	}

	var submitted Submitted
	if err := json.Unmarshal(submittedJSON, &submitted); err != nil {
		return GradeResult{}, ErrInvalidAnswerShape
	}

	switch exerciseType {
	case models.ExerciseMultipleChoice, models.ExerciseListenPickImage, models.ExerciseMinimalPair,
		models.ExerciseIPAToWord, models.ExerciseWordToIPA, models.ExerciseEndingSound:
		return gradeChoice(key, submitted)
	case models.ExerciseTrueFalse:
		if submitted.Value == nil || key.Value == nil {
			return GradeResult{}, shapeOrKeyError(submitted.Value == nil)
		}
		return exact(*submitted.Value == *key.Value, *key.Value), nil
	case models.ExerciseFillBlank:
		return gradeText(key.Accepted, submitted.Text, false)
	case models.ExerciseTransform:
		return gradeText(key.Accepted, submitted.Text, key.AIFallback)
	case models.ExerciseTranslateViEn:
		return gradeText(append(append([]string{}, key.Accepted...), key.Reference...), submitted.Text, key.AIGraded)
	case models.ExerciseReorder:
		if len(submitted.Tokens) == 0 {
			return GradeResult{}, ErrInvalidAnswerShape
		}
		sentence := joinTokens(submitted.Tokens)
		return gradeText(key.Accepted, &sentence, false)
	case models.ExerciseErrorCorrection:
		return gradeErrorCorrection(key, submitted)
	case models.ExerciseDictation:
		return gradeDictation(key, submitted)
	case models.ExerciseListenGapFill:
		return gradeGapFill(key, submitted)
	case models.ExerciseListenNumbers:
		var payload numbersPayload
		_ = json.Unmarshal(payloadJSON, &payload)
		return gradeNumbers(key, submitted, payload.Format)
	case models.ExerciseWordStress:
		if submitted.StressIndex == nil || key.StressIndex == nil {
			return GradeResult{}, shapeOrKeyError(submitted.StressIndex == nil)
		}
		return exact(*submitted.StressIndex == *key.StressIndex, *key.StressIndex), nil
	case models.ExerciseSpeak:
		// A scored recording is graded from its attempt by the service; a
		// self-check (no microphone / scoring not configured) counts as done.
		if submitted.AttemptID != nil && *submitted.AttemptID > 0 {
			return GradeResult{PronunciationAttempt: *submitted.AttemptID}, nil
		}
		if submitted.SelfCheck == nil {
			return GradeResult{}, ErrInvalidAnswerShape
		}
		return GradeResult{IsCorrect: true, ScoreRatio: 1}, nil
	default:
		return GradeResult{}, fmt.Errorf("%w: unsupported exercise type %q", ErrInvalidAnswerKey, exerciseType)
	}
}

func shapeOrKeyError(isShape bool) error {
	if isShape {
		return ErrInvalidAnswerShape
	}
	return ErrInvalidAnswerKey
}

func exact(ok bool, correct any) GradeResult {
	if ok {
		return GradeResult{IsCorrect: true, ScoreRatio: 1, CorrectAnswer: correct}
	}
	return GradeResult{IsCorrect: false, ScoreRatio: 0, CorrectAnswer: correct}
}

func gradeChoice(key answerKey, submitted Submitted) (GradeResult, error) {
	if submitted.OptionID == nil {
		return GradeResult{}, ErrInvalidAnswerShape
	}
	var correct string
	if err := json.Unmarshal(key.Correct, &correct); err != nil || correct == "" {
		return GradeResult{}, ErrInvalidAnswerKey
	}
	return exact(*submitted.OptionID == correct, correct), nil
}

func gradeText(accepted []string, text *string, aiAllowed bool) (GradeResult, error) {
	if text == nil {
		return GradeResult{}, ErrInvalidAnswerShape
	}
	if len(accepted) == 0 {
		return GradeResult{}, ErrInvalidAnswerKey
	}

	if matchesAny(*text, accepted) {
		return GradeResult{IsCorrect: true, ScoreRatio: 1, CorrectAnswer: accepted[0]}, nil
	}
	return GradeResult{IsCorrect: false, ScoreRatio: 0, CorrectAnswer: accepted[0], NeedsAI: aiAllowed && strings.TrimSpace(*text) != ""}, nil
}

func gradeErrorCorrection(key answerKey, submitted Submitted) (GradeResult, error) {
	if submitted.Index == nil || submitted.Text == nil {
		return GradeResult{}, ErrInvalidAnswerShape
	}
	if key.WrongIndex == nil || len(key.Accepted) == 0 {
		return GradeResult{}, ErrInvalidAnswerKey
	}

	correct := map[string]any{"wrongIndex": *key.WrongIndex, "text": key.Accepted[0]}
	if *submitted.Index != *key.WrongIndex {
		return GradeResult{IsCorrect: false, ScoreRatio: 0, CorrectAnswer: correct}, nil
	}
	if matchesAny(*submitted.Text, key.Accepted) {
		return GradeResult{IsCorrect: true, ScoreRatio: 1, CorrectAnswer: correct}, nil
	}
	// Found the mistake but could not fix it.
	return GradeResult{IsCorrect: false, ScoreRatio: 0.5, CorrectAnswer: correct}, nil
}

func gradeDictation(key answerKey, submitted Submitted) (GradeResult, error) {
	if submitted.Text == nil {
		return GradeResult{}, ErrInvalidAnswerShape
	}
	if len(key.Accepted) == 0 {
		return GradeResult{}, ErrInvalidAnswerKey
	}

	best := 0.0
	answerWords := strings.Fields(normalizeText(*submitted.Text))
	for _, reference := range key.Accepted {
		referenceWords := strings.Fields(normalizeText(reference))
		if len(referenceWords) == 0 {
			continue
		}
		ratio := 1 - float64(wordEditDistance(answerWords, referenceWords))/float64(len(referenceWords))
		if ratio > best {
			best = ratio
		}
	}
	if best < 0 {
		best = 0
	}
	best = roundRatio(best)
	return GradeResult{IsCorrect: best >= correctThreshold, ScoreRatio: best, CorrectAnswer: key.Accepted[0]}, nil
}

func gradeGapFill(key answerKey, submitted Submitted) (GradeResult, error) {
	if submitted.Blanks == nil {
		return GradeResult{}, ErrInvalidAnswerShape
	}
	if len(key.Blanks) == 0 {
		return GradeResult{}, ErrInvalidAnswerKey
	}

	correct := make(map[string]string, len(key.Blanks))
	matched := 0
	for blank, accepted := range key.Blanks {
		if len(accepted) == 0 {
			return GradeResult{}, ErrInvalidAnswerKey
		}
		correct[blank] = accepted[0]
		if answer, ok := submitted.Blanks[blank]; ok && matchesAnyNumberAware(answer, accepted) {
			matched++
		}
	}

	ratio := roundRatio(float64(matched) / float64(len(key.Blanks)))
	return GradeResult{IsCorrect: matched == len(key.Blanks), ScoreRatio: ratio, CorrectAnswer: correct}, nil
}

func gradeNumbers(key answerKey, submitted Submitted, format string) (GradeResult, error) {
	if submitted.Text == nil {
		return GradeResult{}, ErrInvalidAnswerShape
	}
	if len(key.Accepted) == 0 {
		return GradeResult{}, ErrInvalidAnswerKey
	}

	answer := normalizeNumberAnswer(*submitted.Text, format)
	for _, accepted := range key.Accepted {
		if answer == normalizeNumberAnswer(accepted, format) {
			return GradeResult{IsCorrect: true, ScoreRatio: 1, CorrectAnswer: key.Accepted[0]}, nil
		}
	}
	return GradeResult{IsCorrect: false, ScoreRatio: 0, CorrectAnswer: key.Accepted[0]}, nil
}

// ---------------------------------------------------------------------------
// Normalisation (docs/foundation-spec.md §5.2)

var (
	whitespacePattern  = regexp.MustCompile(`\s+`)
	spaceBeforePunct   = regexp.MustCompile(`\s+([.,!?;:])`)
	trailingPunct      = regexp.MustCompile(`[.!?]+$`)
	innerPunctPattern  = regexp.MustCompile(`[,;:"]`)
	contractionPattern = regexp.MustCompile(`\b([a-z]+)'([a-z]+)\b`)
)

// Contractions expand to one canonical form so "I'm" and "I am" compare equal.
var contractions = map[string]string{
	"i'm": "i am", "you're": "you are", "we're": "we are", "they're": "they are",
	"he's": "he is", "she's": "she is", "it's": "it is", "that's": "that is",
	"there's": "there is", "what's": "what is", "where's": "where is", "who's": "who is",
	"how's": "how is", "here's": "here is",
	"i've": "i have", "you've": "you have", "we've": "we have", "they've": "they have",
	"i'll": "i will", "you'll": "you will", "he'll": "he will", "she'll": "she will",
	"it'll": "it will", "we'll": "we will", "they'll": "they will",
	"i'd": "i would", "you'd": "you would", "he'd": "he would", "she'd": "she would",
	"we'd": "we would", "they'd": "they would",
	"don't": "do not", "doesn't": "does not", "didn't": "did not",
	"isn't": "is not", "aren't": "are not", "wasn't": "was not", "weren't": "were not",
	"haven't": "have not", "hasn't": "has not", "hadn't": "had not",
	"won't": "will not", "wouldn't": "would not", "can't": "cannot", "couldn't": "could not",
	"shouldn't": "should not", "mustn't": "must not", "let's": "let us",
}

func normalizeText(value string) string {
	text := strings.ToLower(value)
	text = strings.NewReplacer("’", "'", "‘", "'", "“", `"`, "”", `"`).Replace(text)
	text = whitespacePattern.ReplaceAllString(strings.TrimSpace(text), " ")
	text = spaceBeforePunct.ReplaceAllString(text, "$1")
	text = trailingPunct.ReplaceAllString(text, "")
	text = innerPunctPattern.ReplaceAllString(text, "")
	text = contractionPattern.ReplaceAllStringFunc(text, func(word string) string {
		if expanded, ok := contractions[word]; ok {
			return expanded
		}
		return word
	})
	text = strings.ReplaceAll(text, "can not", "cannot")
	return whitespacePattern.ReplaceAllString(strings.TrimSpace(text), " ")
}

func matchesAny(answer string, accepted []string) bool {
	normalized := normalizeText(answer)
	if normalized == "" {
		return false
	}
	for _, candidate := range accepted {
		if normalized == normalizeText(candidate) {
			return true
		}
	}
	return false
}

func matchesAnyNumberAware(answer string, accepted []string) bool {
	if matchesAny(answer, accepted) {
		return true
	}
	normalized := normalizeNumberAnswer(answer, "NUMBER")
	for _, candidate := range accepted {
		if normalized != "" && normalized == normalizeNumberAnswer(candidate, "NUMBER") {
			return true
		}
	}
	return false
}

// joinTokens turns reordered tokens back into a sentence, keeping punctuation
// attached to the previous word.
func joinTokens(tokens []string) string {
	sentence := strings.Join(tokens, " ")
	return spaceBeforePunct.ReplaceAllString(sentence, "$1")
}

var numberWords = map[string]int{
	"zero": 0, "oh": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6,
	"seven": 7, "eight": 8, "nine": 9, "ten": 10, "eleven": 11, "twelve": 12,
	"thirteen": 13, "fourteen": 14, "fifteen": 15, "sixteen": 16, "seventeen": 17,
	"eighteen": 18, "nineteen": 19, "twenty": 20, "thirty": 30, "forty": 40,
	"fifty": 50, "sixty": 60, "seventy": 70, "eighty": 80, "ninety": 90,
}

var (
	timePattern  = regexp.MustCompile(`^(\d{1,2})[:.h ](\d{2})$`)
	digitsOnly   = regexp.MustCompile(`[^0-9]`)
	pricePattern = regexp.MustCompile(`^\$?\s*(\d+)(?:[.,](\d{1,2}))?\s*(?:dollars?|usd|\$)?$`)
)

// wordsToNumberText converts number words to digits ("seven thirty" -> "7 30",
// "twenty five" -> "25") and leaves other words unchanged.
func wordsToNumberText(value string) string {
	words := strings.Fields(strings.ReplaceAll(strings.ToLower(value), "-", " "))
	out := make([]string, 0, len(words))
	pending := -1
	flush := func() {
		if pending >= 0 {
			out = append(out, strconv.Itoa(pending))
			pending = -1
		}
	}
	for _, word := range words {
		number, ok := numberWords[word]
		switch {
		case !ok:
			flush()
			out = append(out, word)
		case pending >= 20 && pending%10 == 0 && number < 10:
			pending += number
		default:
			flush()
			pending = number
		}
	}
	flush()
	return strings.Join(out, " ")
}

func normalizeNumberAnswer(value string, format string) string {
	text := strings.TrimSpace(strings.ToLower(value))
	text = strings.TrimSuffix(text, ".")
	text = wordsToNumberText(text)
	text = strings.ReplaceAll(text, "o'clock", "")
	text = strings.ReplaceAll(text, "oclock", "")
	text = strings.TrimSpace(text)

	switch strings.ToUpper(format) {
	case "TIME":
		text = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(text, "am"), "pm"))
		text = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(text, "a.m"), "p.m"))
		if match := timePattern.FindStringSubmatch(text); match != nil {
			hour, _ := strconv.Atoi(match[1])
			return fmt.Sprintf("%d:%s", hour, match[2])
		}
		if hour, err := strconv.Atoi(text); err == nil {
			return fmt.Sprintf("%d:00", hour)
		}
	case "PRICE":
		if match := pricePattern.FindStringSubmatch(text); match != nil {
			cents := match[2]
			if cents == "" || strings.Trim(cents, "0") == "" {
				return match[1]
			}
			if len(cents) == 1 {
				cents += "0"
			}
			return match[1] + "." + cents
		}
	case "PHONE":
		return digitsOnly.ReplaceAllString(text, "")
	}

	return whitespacePattern.ReplaceAllString(text, " ")
}

// wordEditDistance is the Levenshtein distance counted in words.
func wordEditDistance(a []string, b []string) int {
	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(a); i++ {
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			current[j] = minOf(previous[j]+1, current[j-1]+1, previous[j-1]+cost)
		}
		previous, current = current, previous
	}
	return previous[len(b)]
}

func minOf(values ...int) int {
	result := values[0]
	for _, value := range values[1:] {
		if value < result {
			result = value
		}
	}
	return result
}

func roundRatio(value float64) float64 {
	return float64(int(value*1000+0.5)) / 1000
}

// gradeSpeak grades a scored recording: correct from minAccuracy, half credit
// within 15 points of it.
func gradeSpeak(accuracy, minAccuracy int) GradeResult {
	correct := map[string]int{"accuracy": accuracy, "minAccuracy": minAccuracy}
	switch {
	case accuracy >= minAccuracy:
		return GradeResult{IsCorrect: true, ScoreRatio: 1, CorrectAnswer: correct}
	case accuracy >= minAccuracy-15:
		return GradeResult{ScoreRatio: 0.5, CorrectAnswer: correct}
	default:
		return GradeResult{CorrectAnswer: correct}
	}
}
