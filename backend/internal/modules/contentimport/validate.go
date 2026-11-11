package contentimport

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/lessonactivity"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

var activityTypes = map[string]bool{
	string(models.ActivityVocabulary): true, string(models.ActivityPronunciation): true,
	string(models.ActivityGrammar): true, string(models.ActivityListening): true,
	string(models.ActivityReading): true, string(models.ActivityWriting): true,
	string(models.ActivitySpeaking): true, string(models.ActivityCheckpoint): true,
}

func validateSlug(slug string) []string {
	if slug == "" {
		return []string{"slug is required"}
	}
	if !slugPattern.MatchString(slug) {
		return []string{fmt.Sprintf("slug %q must use lowercase letters, digits and dashes", slug)}
	}
	return nil
}

func validateTopic(topic TopicInput) []string {
	errs := validateSlug(topic.Slug)
	if strings.TrimSpace(topic.Title) == "" {
		errs = append(errs, "title is required")
	}
	if topic.Track != "" && topic.Track != string(models.TopicTrackCore) && topic.Track != string(models.TopicTrackVocabulary) {
		errs = append(errs, fmt.Sprintf("track %q must be CORE or VOCABULARY", topic.Track))
	}
	return errs
}

func validateLesson(lesson LessonInput) []string {
	errs := validateSlug(lesson.Slug)
	if strings.TrimSpace(lesson.Title) == "" {
		errs = append(errs, "title is required")
	}
	if lesson.RequiredScore < 0 || lesson.RequiredScore > 100 {
		errs = append(errs, "requiredScore must be between 0 and 100")
	}
	return errs
}

// validateActivity checks one activity; knownVocabulary holds the vocabulary
// slugs that exist in the database.
func validateActivity(activity ActivityInput, knownVocabulary map[string]bool) []string {
	errs := validateSlug(activity.Slug)
	if !activityTypes[activity.Type] {
		errs = append(errs, fmt.Sprintf("unknown activity type %q", activity.Type))
	}
	if strings.TrimSpace(activity.Title) == "" {
		errs = append(errs, "title is required")
	}
	if activity.Status != "" && activity.Status != string(models.ContentDraft) && activity.Status != string(models.ContentPublished) {
		errs = append(errs, fmt.Sprintf("status %q must be DRAFT or PUBLISHED", activity.Status))
	}
	if activity.PassScore != nil && (*activity.PassScore < 0 || *activity.PassScore > 100) {
		errs = append(errs, "passScore must be between 0 and 100")
	}
	if len(activity.Content) > 0 && !isJSONObject(activity.Content) {
		errs = append(errs, "content must be a JSON object")
	}

	if activity.Type == string(models.ActivityWriting) {
		errs = append(errs, validateWritingContent(activity.Content)...)
	}

	if activity.Type == string(models.ActivityVocabulary) {
		slugs := vocabularySlugs(activity.Content)
		if len(slugs) == 0 {
			errs = append(errs, "content.vocabularySlugs must list at least one word")
		}
		var missing []string
		for _, slug := range slugs {
			if !knownVocabulary[slug] {
				missing = append(missing, slug)
			}
		}
		if len(missing) > 0 {
			errs = append(errs, "vocabulary not found: "+strings.Join(missing, ", "))
		}
	}
	return errs
}

func vocabularySlugs(content json.RawMessage) []string {
	var parsed struct {
		VocabularySlugs []string `json:"vocabularySlugs"`
	}
	if len(content) == 0 {
		return nil
	}
	_ = json.Unmarshal(content, &parsed)
	return parsed.VocabularySlugs
}

// exercisePayload / exerciseAnswer list every field any exercise type uses.
type exercisePayload struct {
	Options   json.RawMessage `json:"options"`
	Tokens    []string        `json:"tokens"`
	Text      string          `json:"text"`
	Sentence  string          `json:"sentence"`
	SourceVI  string          `json:"sourceVi"`
	Statement string          `json:"statement"`
	Syllables []string        `json:"syllables"`
	Word      string          `json:"word"`
	IPA       string          `json:"ipa"`
	AudioURL  string          `json:"audioUrl"`
	// A DICTATION item may play one segment of a passage instead of audioUrl.
	PassageSlug  string `json:"passageSlug"`
	PassageID    uint   `json:"passageId"`
	SegmentIndex *int   `json:"segmentIndex"`
}

type exerciseAnswer struct {
	Correct     *string             `json:"correct"`
	Value       *bool               `json:"value"`
	Accepted    []string            `json:"accepted"`
	Reference   []string            `json:"reference"`
	WrongIndex  *int                `json:"wrongIndex"`
	Blanks      map[string][]string `json:"blanks"`
	StressIndex *int                `json:"stressIndex"`
	MinAccuracy *int                `json:"minAccuracy"`
}

var blankPattern = regexp.MustCompile(`\{\{(\w+)\}\}`)

// validateExercise checks an exercise's payload and answer for its type, then
// grades the expected correct answer to prove the item is answerable.
func validateExercise(exercise ExerciseInput) []string {
	var errs []string
	exerciseType := models.ExerciseType(exercise.Type)
	if strings.TrimSpace(exercise.Prompt) == "" {
		errs = append(errs, "prompt is required")
	}
	if !isJSONObject(exercise.Payload) {
		return append(errs, "payload must be a JSON object")
	}
	if !isJSONObject(exercise.Answer) {
		return append(errs, "answer must be a JSON object")
	}

	var payload exercisePayload
	var answer exerciseAnswer
	if err := json.Unmarshal(exercise.Payload, &payload); err != nil {
		return append(errs, "payload: "+err.Error())
	}
	if err := json.Unmarshal(exercise.Answer, &answer); err != nil {
		return append(errs, "answer: "+err.Error())
	}

	var canonical any
	switch exerciseType {
	case models.ExerciseMultipleChoice, models.ExerciseListenPickImage, models.ExerciseMinimalPair,
		models.ExerciseIPAToWord, models.ExerciseWordToIPA, models.ExerciseEndingSound:
		ids, err := optionIDs(payload.Options)
		if err != nil {
			errs = append(errs, err.Error())
		} else if len(ids) < 2 {
			errs = append(errs, "payload.options needs at least 2 options")
		}
		if answer.Correct == nil {
			errs = append(errs, "answer.correct is required")
		} else if !contains(ids, *answer.Correct) {
			errs = append(errs, fmt.Sprintf("answer.correct %q is not one of the option ids", *answer.Correct))
		}
		if exerciseType == models.ExerciseIPAToWord && payload.IPA == "" {
			errs = append(errs, "payload.ipa is required")
		}
		if answer.Correct != nil {
			canonical = map[string]any{"optionId": *answer.Correct}
		}
	case models.ExerciseTrueFalse:
		if payload.Statement == "" {
			errs = append(errs, "payload.statement is required")
		}
		if answer.Value == nil {
			errs = append(errs, "answer.value (true/false) is required")
		} else {
			canonical = map[string]any{"value": *answer.Value}
		}
	case models.ExerciseFillBlank, models.ExerciseTransform, models.ExerciseDictation, models.ExerciseListenNumbers:
		if exerciseType == models.ExerciseFillBlank && payload.Sentence == "" {
			errs = append(errs, "payload.sentence is required")
		}
		if exerciseType == models.ExerciseTransform && payload.Sentence == "" {
			errs = append(errs, "payload.sentence is required")
		}
		usesSegment := (payload.PassageSlug != "" || payload.PassageID != 0) && payload.SegmentIndex != nil
		if exerciseType == models.ExerciseDictation && payload.AudioURL == "" && !usesSegment {
			errs = append(errs, "payload.audioUrl (or passageSlug + segmentIndex) is required")
		}
		if exerciseType == models.ExerciseListenNumbers && payload.AudioURL == "" {
			errs = append(errs, "payload.audioUrl is required")
		}
		if len(answer.Accepted) == 0 {
			errs = append(errs, "answer.accepted needs at least one answer")
		} else {
			canonical = map[string]any{"text": answer.Accepted[0]}
		}
	case models.ExerciseTranslateViEn:
		if payload.SourceVI == "" {
			errs = append(errs, "payload.sourceVi is required")
		}
		all := append(append([]string{}, answer.Accepted...), answer.Reference...)
		if len(all) == 0 {
			errs = append(errs, "answer.reference needs at least one translation")
		} else {
			canonical = map[string]any{"text": all[0]}
		}
	case models.ExerciseReorder:
		if len(payload.Tokens) < 2 {
			errs = append(errs, "payload.tokens needs at least 2 tokens")
		}
		if len(answer.Accepted) == 0 {
			errs = append(errs, "answer.accepted needs at least one sentence")
		} else if len(payload.Tokens) > 0 && !sameTokens(payload.Tokens, answer.Accepted[0]) {
			errs = append(errs, "payload.tokens must be exactly the words of answer.accepted[0]")
		}
		// Tokens are shuffled on purpose, so the self-check grades the sentence order.
		if len(answer.Accepted) > 0 {
			canonical = map[string]any{"tokens": sentenceTokens(answer.Accepted[0])}
		}
	case models.ExerciseErrorCorrection:
		if len(payload.Tokens) < 2 {
			errs = append(errs, "payload.tokens needs at least 2 tokens")
		}
		if answer.WrongIndex == nil || *answer.WrongIndex < 0 || *answer.WrongIndex >= len(payload.Tokens) {
			errs = append(errs, "answer.wrongIndex must point to a token")
		}
		if len(answer.Accepted) == 0 {
			errs = append(errs, "answer.accepted needs the corrected word")
		}
		if answer.WrongIndex != nil && len(answer.Accepted) > 0 {
			canonical = map[string]any{"index": *answer.WrongIndex, "text": answer.Accepted[0]}
		}
	case models.ExerciseListenGapFill:
		blanks := blankPattern.FindAllStringSubmatch(payload.Text, -1)
		if len(blanks) == 0 {
			errs = append(errs, "payload.text needs blanks like {{1}}")
		}
		if payload.AudioURL == "" {
			errs = append(errs, "payload.audioUrl is required")
		}
		first := map[string]string{}
		for _, blank := range blanks {
			accepted := answer.Blanks[blank[1]]
			if len(accepted) == 0 {
				errs = append(errs, fmt.Sprintf("answer.blanks.%s needs at least one answer", blank[1]))
				continue
			}
			first[blank[1]] = accepted[0]
		}
		if len(answer.Blanks) != len(blanks) {
			errs = append(errs, "answer.blanks must have one entry per blank in payload.text")
		}
		canonical = map[string]any{"blanks": first}
	case models.ExerciseWordStress:
		if len(payload.Syllables) < 2 {
			errs = append(errs, "payload.syllables needs at least 2 syllables")
		}
		if answer.StressIndex == nil || *answer.StressIndex < 0 || *answer.StressIndex >= len(payload.Syllables) {
			errs = append(errs, "answer.stressIndex must point to a syllable")
		} else {
			canonical = map[string]any{"stressIndex": *answer.StressIndex}
		}
	case models.ExerciseSpeak:
		if payload.Text == "" {
			errs = append(errs, "payload.text is required")
		}
		if answer.MinAccuracy != nil && (*answer.MinAccuracy < 1 || *answer.MinAccuracy > 100) {
			errs = append(errs, "answer.minAccuracy must be 1-100")
		}
		canonical = map[string]any{"selfCheck": true}
	default:
		return append(errs, fmt.Sprintf("unknown exercise type %q", exercise.Type))
	}

	if len(errs) > 0 || canonical == nil {
		return errs
	}

	submitted, _ := json.Marshal(canonical)
	result, err := lessonactivity.Grade(exerciseType, exercise.Payload, exercise.Answer, submitted)
	if err != nil {
		return append(errs, "self-check failed: "+err.Error())
	}
	if !result.IsCorrect {
		return append(errs, "self-check failed: the expected answer is graded as wrong")
	}
	return nil
}

func optionIDs(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("payload.options is required")
	}
	var plain []string
	if err := json.Unmarshal(raw, &plain); err == nil {
		return plain, nil
	}
	var objects []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &objects); err != nil {
		return nil, fmt.Errorf("payload.options must be a list of strings or objects with an id")
	}
	ids := make([]string, 0, len(objects))
	seen := map[string]bool{}
	for _, option := range objects {
		if option.ID == "" {
			return nil, fmt.Errorf("every option needs an id")
		}
		if seen[option.ID] {
			return nil, fmt.Errorf("duplicate option id %q", option.ID)
		}
		seen[option.ID] = true
		ids = append(ids, option.ID)
	}
	return ids, nil
}

var tokenPattern = regexp.MustCompile(`[\p{L}\p{N}'’-]+|[.,!?;:]`)

// sentenceTokens splits a sentence into words and punctuation marks.
func sentenceTokens(sentence string) []string {
	return tokenPattern.FindAllString(sentence, -1)
}

func sameTokens(tokens []string, sentence string) bool {
	expected := sentenceTokens(sentence)
	if len(expected) != len(tokens) {
		return false
	}
	a := lowerSorted(tokens)
	b := lowerSorted(expected)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func lowerSorted(values []string) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = strings.ToLower(strings.TrimSpace(value))
	}
	sort.Strings(result)
	return result
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func isJSONObject(raw json.RawMessage) bool {
	var value map[string]json.RawMessage
	return len(raw) > 0 && json.Unmarshal(raw, &value) == nil && value != nil
}

type grammarExample struct {
	EN string `json:"en"`
	VI string `json:"vi"`
}

type grammarMistake struct {
	Wrong string `json:"wrong"`
	Right string `json:"right"`
}

func validateGrammarPoint(point GrammarPointInput) []string {
	errs := validateSlug(point.Slug)
	if strings.TrimSpace(point.Title) == "" {
		errs = append(errs, "title is required")
	}
	if strings.TrimSpace(point.SummaryVI) == "" {
		errs = append(errs, "summaryVi is required")
	}

	var examples []grammarExample
	if err := json.Unmarshal(point.Examples, &examples); err != nil || len(examples) == 0 {
		errs = append(errs, "examples must be a non-empty list of {en, vi}")
	} else {
		for i, example := range examples {
			if strings.TrimSpace(example.EN) == "" {
				errs = append(errs, fmt.Sprintf("examples[%d].en is required", i))
			}
		}
	}

	if len(point.CommonMistakes) > 0 && string(point.CommonMistakes) != "null" {
		var mistakes []grammarMistake
		if err := json.Unmarshal(point.CommonMistakes, &mistakes); err != nil {
			errs = append(errs, "commonMistakes must be a list of {wrong, right, noteVi}")
		} else {
			for i, mistake := range mistakes {
				if mistake.Wrong == "" || mistake.Right == "" {
					errs = append(errs, fmt.Sprintf("commonMistakes[%d] needs wrong and right", i))
				}
			}
		}
	}
	return errs
}

// grammarPointSlug returns the activity's linked grammar point slug.
func grammarPointSlug(activity ActivityInput) string {
	if activity.GrammarPointSlug != "" {
		return activity.GrammarPointSlug
	}
	var content struct {
		GrammarPointSlug string `json:"grammarPointSlug"`
	}
	if len(activity.Content) > 0 {
		_ = json.Unmarshal(activity.Content, &content)
	}
	return content.GrammarPointSlug
}

type phonemeWord struct {
	Word string `json:"word"`
	IPA  string `json:"ipa"`
}

type phonemePair struct {
	A string `json:"a"`
	B string `json:"b"`
}

func validatePhoneme(phoneme PhonemeInput) []string {
	var errs []string
	if strings.TrimSpace(phoneme.Symbol) == "" || len(phoneme.Symbol) > 10 {
		errs = append(errs, "symbol is required (max 10 bytes)")
	}
	switch models.PhonemeKind(phoneme.Kind) {
	case models.PhonemeVowel, models.PhonemeDiphthong:
		if phoneme.Voiced != nil {
			errs = append(errs, "voiced is only for consonants")
		}
	case models.PhonemeConsonant:
		if phoneme.Voiced == nil {
			errs = append(errs, "voiced is required for consonants")
		}
	default:
		errs = append(errs, fmt.Sprintf("kind %q must be VOWEL, DIPHTHONG or CONSONANT", phoneme.Kind))
	}
	if strings.TrimSpace(phoneme.MouthVI) == "" {
		errs = append(errs, "mouthVi is required")
	}

	var words []phonemeWord
	if err := json.Unmarshal(phoneme.ExampleWords, &words); err != nil || len(words) == 0 {
		errs = append(errs, "exampleWords must be a non-empty list of {word, ipa}")
	} else {
		for i, word := range words {
			if word.Word == "" || word.IPA == "" {
				errs = append(errs, fmt.Sprintf("exampleWords[%d] needs word and ipa", i))
			}
		}
	}

	if len(phoneme.MinimalPairs) > 0 && string(phoneme.MinimalPairs) != "null" {
		var pairs []phonemePair
		if err := json.Unmarshal(phoneme.MinimalPairs, &pairs); err != nil {
			errs = append(errs, "minimalPairs must be a list of {a, b}")
		} else {
			for i, pair := range pairs {
				if pair.A == "" || pair.B == "" {
					errs = append(errs, fmt.Sprintf("minimalPairs[%d] needs a and b", i))
				}
			}
		}
	}
	return errs
}

func passageSlug(activity ActivityInput) string {
	if activity.PassageSlug != "" {
		return activity.PassageSlug
	}
	var content struct {
		PassageSlug string `json:"passageSlug"`
	}
	if len(activity.Content) > 0 {
		_ = json.Unmarshal(activity.Content, &content)
	}
	return content.PassageSlug
}

var passageGenres = map[string]bool{
	string(models.GenreDialogue): true, string(models.GenreMessage): true, string(models.GenreEmail): true,
	string(models.GenreNotice): true, string(models.GenreMenu): true, string(models.GenreSchedule): true,
	string(models.GenreArticle): true, string(models.GenreProfile): true, string(models.GenreMonologue): true,
}

type passageSegment struct {
	Index    *int     `json:"index"`
	Speaker  string   `json:"speaker"`
	Text     string   `json:"text"`
	Start    *float64 `json:"start"`
	End      *float64 `json:"end"`
	AudioURL string   `json:"audioUrl"`
}

type glossaryEntry struct {
	Word string `json:"word"`
}

func validatePassage(passage PassageInput) []string {
	errs := validateSlug(passage.Slug)
	if !passageGenres[passage.Genre] {
		errs = append(errs, fmt.Sprintf("genre %q must be one of DIALOGUE, MESSAGE, EMAIL, NOTICE, MENU, SCHEDULE, ARTICLE, PROFILE, MONOLOGUE", passage.Genre))
	}
	if strings.TrimSpace(passage.Title) == "" {
		errs = append(errs, "title is required")
	}
	if strings.TrimSpace(passage.BodyText) == "" {
		errs = append(errs, "bodyText is required")
	}

	if len(passage.Segments) > 0 && string(passage.Segments) != "null" {
		var segments []passageSegment
		if err := json.Unmarshal(passage.Segments, &segments); err != nil {
			errs = append(errs, "segments must be a list of {index, speaker?, text, start?, end?, audioUrl?}")
		} else {
			for i, segment := range segments {
				if strings.TrimSpace(segment.Text) == "" {
					errs = append(errs, fmt.Sprintf("segments[%d].text is required", i))
				}
				if segment.Index != nil && *segment.Index != i {
					errs = append(errs, fmt.Sprintf("segments[%d].index must be %d", i, i))
				}
				if (segment.Start == nil) != (segment.End == nil) {
					errs = append(errs, fmt.Sprintf("segments[%d] needs both start and end", i))
				} else if segment.Start != nil && *segment.End <= *segment.Start {
					errs = append(errs, fmt.Sprintf("segments[%d].end must be after start", i))
				}
				if segment.Start != nil && passage.AudioURL == "" {
					errs = append(errs, fmt.Sprintf("segments[%d] has times but the passage has no audioUrl", i))
				}
			}
		}
	}

	if len(passage.Glossary) > 0 && string(passage.Glossary) != "null" {
		var glossary []glossaryEntry
		if err := json.Unmarshal(passage.Glossary, &glossary); err != nil {
			errs = append(errs, "glossary must be a list of {word, vocabularySlug?, meaningVi?, ipa?}")
		} else {
			for i, entry := range glossary {
				if strings.TrimSpace(entry.Word) == "" {
					errs = append(errs, fmt.Sprintf("glossary[%d].word is required", i))
				}
			}
		}
	}
	return errs
}

type writingPrompt struct {
	PromptVI     string   `json:"promptVi"`
	PromptEN     string   `json:"promptEn"`
	Frame        []string `json:"frame"`
	MinSentences int      `json:"minSentences"`
	MaxSentences int      `json:"maxSentences"`
	ModelAnswer  string   `json:"modelAnswer"`
}

// validateWritingContent checks a WRITING prompt (docs/foundation-spec.md §4.1).
func validateWritingContent(content json.RawMessage) []string {
	var prompt writingPrompt
	if len(content) == 0 || json.Unmarshal(content, &prompt) != nil {
		return []string{"content must be a writing prompt {promptVi, promptEn, frame, minSentences, maxSentences, modelAnswer}"}
	}
	var errs []string
	if strings.TrimSpace(prompt.PromptVI) == "" {
		errs = append(errs, "content.promptVi is required")
	}
	if prompt.MinSentences < 1 || prompt.MaxSentences < prompt.MinSentences {
		errs = append(errs, "content needs 1 <= minSentences <= maxSentences")
	}
	if strings.TrimSpace(prompt.ModelAnswer) == "" {
		errs = append(errs, "content.modelAnswer is required")
	}
	return errs
}

func scenarioSlug(activity ActivityInput) string {
	if activity.ScenarioSlug != "" {
		return activity.ScenarioSlug
	}
	var content struct {
		ScenarioSlug string `json:"scenarioSlug"`
	}
	if len(activity.Content) > 0 {
		_ = json.Unmarshal(activity.Content, &content)
	}
	return content.ScenarioSlug
}

func validateScenario(scenario ScenarioInput) []string {
	errs := validateSlug(scenario.Slug)
	if strings.TrimSpace(scenario.Title) == "" {
		errs = append(errs, "title is required")
	}
	if strings.TrimSpace(scenario.Role) == "" || strings.TrimSpace(scenario.Task) == "" {
		errs = append(errs, "role and task are required")
	}
	if strings.TrimSpace(scenario.OpeningLine) == "" {
		errs = append(errs, "openingLine is required")
	}
	if scenario.MaxTurns < 0 || scenario.MaxTurns > 30 {
		errs = append(errs, "maxTurns must be between 1 and 30")
	}
	return errs
}
