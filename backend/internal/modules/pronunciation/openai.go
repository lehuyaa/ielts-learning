package pronunciation

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"

	"ielts-learning/backend/internal/shared/ai"
)

// Providers that can score a recording.
const (
	ProviderAzure  = "AZURE"
	ProviderOpenAI = "OPENAI"
)

// transcriber is the part of ai.Client the OpenAI scorer needs (fake in tests).
type transcriber interface {
	Configured() bool
	Transcribe(ctx context.Context, filename string, data []byte) (ai.Transcription, error)
}

// OpenAIScorer is the fallback when Azure is not available: the recording is
// transcribed (Whisper) and the transcript is aligned with the reference text.
// It scores words only — a word counts as right when it is heard as that word
// — so phoneme lists stay empty and sound progress is not updated.
type OpenAIScorer struct {
	client transcriber
}

func NewOpenAIScorer(client *ai.Client) *OpenAIScorer {
	scorer := &OpenAIScorer{}
	if client != nil {
		scorer.client = client
	}
	return scorer
}

func (s *OpenAIScorer) Configured() bool {
	return s != nil && s.client != nil && s.client.Configured()
}

func (s *OpenAIScorer) Assess(ctx context.Context, wav []byte, referenceText string) (Assessment, error) {
	transcript, err := s.client.Transcribe(ctx, "recording.wav", wav)
	if err != nil {
		return Assessment{}, err
	}
	heard := strings.TrimSpace(transcript.Text)
	if heard == "" {
		return Assessment{}, ErrNoSpeech
	}
	return scoreTranscript(referenceText, heard, transcript.Duration), nil
}

var (
	wordToken  = regexp.MustCompile(`[A-Za-z0-9'’]+(?::[0-9]+)?`)
	clockToken = regexp.MustCompile(`^(\d{1,2}):(\d{2})$`)
)

// scoreTranscript aligns the reference words with the words heard
// (edit distance over words) and scores each reference word: heard as is
// 100, heard as another word up to 60 by spelling closeness, not heard 0.
func scoreTranscript(reference, heard string, duration float64) Assessment {
	refDisplay := wordToken.FindAllString(reference, -1)
	refWords := make([]string, len(refDisplay))
	for i, word := range refDisplay {
		refWords[i] = normalizeWord(word)
	}
	var heardWords []string
	for _, word := range wordToken.FindAllString(heard, -1) {
		heardWords = append(heardWords, strings.Fields(spellOut(word))...)
	}

	n, m := len(refWords), len(heardWords)
	cost := make([][]int, n+1)
	for i := range cost {
		cost[i] = make([]int, m+1)
		cost[i][0] = i
	}
	for j := 0; j <= m; j++ {
		cost[0][j] = j
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			substitution := cost[i-1][j-1]
			if refWords[i-1] != heardWords[j-1] {
				substitution++
			}
			cost[i][j] = min(substitution, cost[i-1][j]+1, cost[i][j-1]+1)
		}
	}

	words := make([]AssessedWord, n)
	for i, j := n, m; i > 0; {
		switch {
		case j > 0 && cost[i][j] == cost[i-1][j-1] && refWords[i-1] == heardWords[j-1]:
			words[i-1] = AssessedWord{Word: refDisplay[i-1], Accuracy: 100, ErrorType: "None"}
			i, j = i-1, j-1
		case j > 0 && cost[i][j] == cost[i-1][j-1]+1:
			score := math.Round(similarity(refWords[i-1], heardWords[j-1]) * 60)
			words[i-1] = AssessedWord{Word: refDisplay[i-1], Accuracy: score, ErrorType: "Mispronunciation"}
			i, j = i-1, j-1
		case cost[i][j] == cost[i-1][j]+1:
			words[i-1] = AssessedWord{Word: refDisplay[i-1], Accuracy: 0, ErrorType: "Omission"}
			i--
		default: // an extra word the learner said
			j--
		}
	}

	total, said := 0.0, 0
	for _, word := range words {
		total += word.Accuracy
		if word.ErrorType != "Omission" {
			said++
		}
	}
	assessment := Assessment{RecognizedText: heard, Words: words}
	if n > 0 {
		assessment.Accuracy = total / float64(n)
		assessment.Completeness = float64(said) * 100 / float64(n)
	}
	assessment.Pronunciation = assessment.Accuracy
	// Without timing per word, fluency is the speaking rate on sentences of
	// 4+ words (1.5 words a second is fine at A1); short ones use completeness.
	assessment.Fluency = assessment.Completeness
	if duration > 0 && said >= 4 {
		assessment.Fluency = math.Min(100, float64(said)/duration/1.5*100)
	}
	return assessment
}

func normalizeWord(word string) string {
	return strings.Trim(strings.ToLower(strings.ReplaceAll(word, "’", "'")), "'")
}

// spellOut writes numbers the way they were read: "8:30" → "eight thirty",
// "15" → "fifteen", so a transcript with digits matches the reference words.
func spellOut(word string) string {
	word = normalizeWord(word)
	if match := clockToken.FindStringSubmatch(word); match != nil {
		hour, _ := strconv.Atoi(match[1])
		minute, _ := strconv.Atoi(match[2])
		if minute == 0 {
			return numberWords(hour)
		}
		if minute < 10 {
			return numberWords(hour) + " oh " + numberWords(minute)
		}
		return numberWords(hour) + " " + numberWords(minute)
	}
	if value, err := strconv.Atoi(word); err == nil && value >= 0 && value < 10000 {
		return numberWords(value)
	}
	return word
}

var (
	ones = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten",
		"eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
)

func numberWords(value int) string {
	switch {
	case value < 20:
		return ones[value]
	case value < 100:
		if value%10 == 0 {
			return tens[value/10]
		}
		return tens[value/10] + " " + ones[value%10]
	case value < 1000:
		rest := ""
		if value%100 != 0 {
			rest = " and " + numberWords(value%100)
		}
		return ones[value/100] + " hundred" + rest
	default:
		rest := ""
		if value%1000 != 0 {
			rest = " " + numberWords(value%1000)
		}
		return numberWords(value/1000) + " thousand" + rest
	}
}

// similarity is 1 - edit distance / longer length, over letters.
func similarity(a, b string) float64 {
	ra, rb := []rune(a), []rune(b)
	if len(ra) == 0 && len(rb) == 0 {
		return 1
	}
	previous := make([]int, len(rb)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		current := make([]int, len(rb)+1)
		current[0] = i
		for j := 1; j <= len(rb); j++ {
			change := 1
			if ra[i-1] == rb[j-1] {
				change = 0
			}
			current[j] = min(previous[j-1]+change, previous[j]+1, current[j-1]+1)
		}
		previous = current
	}
	return 1 - float64(previous[len(rb)])/float64(max(len(ra), len(rb)))
}
