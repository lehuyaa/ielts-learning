package lessonactivity

import (
	"encoding/json"
	"time"

	"ielts-learning/backend/internal/models"
)

// ActivityStatusNotStarted is returned for activities without progress.
const ActivityStatusNotStarted = "NOT_STARTED"

type LessonActivitiesResponse struct {
	Lesson     LessonSummary      `json:"lesson"`
	Progress   ActivitiesProgress `json:"progress"`
	Activities []ActivityItem     `json:"activities"`
}

type LessonSummary struct {
	ID     uint                `json:"id"`
	Title  string              `json:"title"`
	Slug   string              `json:"slug"`
	Status models.LessonStatus `json:"status"`
}

type ActivitiesProgress struct {
	PassedRequired int `json:"passedRequired"`
	TotalRequired  int `json:"totalRequired"`
	// NextActivityID is the first unfinished activity, for resuming.
	NextActivityID *uint `json:"nextActivityId"`
}

type ActivityItem struct {
	ID             uint                `json:"id"`
	Slug           string              `json:"slug"`
	Type           models.ActivityType `json:"type"`
	Title          string              `json:"title"`
	InstructionsVI string              `json:"instructionsVi"`
	OrderIndex     int                 `json:"orderIndex"`
	IsRequired     bool                `json:"isRequired"`
	PassScore      int                 `json:"passScore"`
	XPReward       int                 `json:"xpReward"`
	ExerciseCount  int                 `json:"exerciseCount"`
	Status         string              `json:"status"`
	Score          *int                `json:"score"`
	BestScore      *int                `json:"bestScore"`
	Attempts       int                 `json:"attempts"`
	IsLocked       bool                `json:"isLocked"`
}

type ActivityDetailResponse struct {
	Activity  ActivityDetail     `json:"activity"`
	Exercises []ExerciseResponse `json:"exercises"`
	// Vocabularies are the word cards of a VOCABULARY activity, in content order.
	Vocabularies []VocabularyCard `json:"vocabularies"`
	// GrammarPoint is the theory card of a GRAMMAR activity, when linked.
	GrammarPoint *GrammarPointCard `json:"grammarPoint"`
	// Phonemes are the sound cards of a PRONUNCIATION activity (content.phonemes).
	Phonemes []PhonemeCard `json:"phonemes"`
	// Passage is the text/recording of a LISTENING or READING activity.
	Passage *PassageCard `json:"passage"`
	// Writing is the learner's latest submission of a WRITING activity.
	Writing *WritingSubmissionResponse `json:"writing"`
	// Scenario is the role-play of a SPEAKING activity; Speaking its latest finished session.
	Scenario *ScenarioCard    `json:"scenario"`
	Speaking *SpeakingSummary `json:"speaking"`
	// Answers already given in the current (unfinished) attempt, for resuming.
	Answers []AnswerResult `json:"answers"`
}

type ActivityDetail struct {
	ActivityItem
	LessonID       uint            `json:"lessonId"`
	Content        json.RawMessage `json:"content" swaggertype:"object"`
	GrammarPointID *uint           `json:"grammarPointId"`
	PassageID      *uint           `json:"passageId"`
	ScenarioID     *uint           `json:"scenarioId"`
	CurrentAttempt int             `json:"currentAttempt"`
}

// ExerciseResponse never contains the answer key.
type ExerciseResponse struct {
	ID         uint                `json:"id"`
	Type       models.ExerciseType `json:"type"`
	Prompt     string              `json:"prompt"`
	Payload    json.RawMessage     `json:"payload" swaggertype:"object"`
	Points     int                 `json:"points"`
	OrderIndex int                 `json:"orderIndex"`
}

type SubmitAnswerRequest struct {
	ExerciseID uint            `json:"exerciseId"`
	Answer     json.RawMessage `json:"answer" swaggertype:"object"`
}

type AnswerResult struct {
	ExerciseID    uint            `json:"exerciseId"`
	Attempt       int             `json:"attempt"`
	IsCorrect     bool            `json:"isCorrect"`
	ScoreRatio    float64         `json:"scoreRatio"`
	Answer        json.RawMessage `json:"answer" swaggertype:"object"`
	CorrectAnswer any             `json:"correctAnswer"`
	ExplanationVI string          `json:"explanationVi"`
	AnsweredAt    time.Time       `json:"answeredAt"`
	// GradedBy is MATCH or AI; empty for answers loaded from history.
	GradedBy   string      `json:"gradedBy,omitempty"`
	AIFeedback *AIFeedback `json:"aiFeedback,omitempty"`
	// AIStatus is LIMIT_REACHED or UNAVAILABLE when AI grading was skipped.
	AIStatus string `json:"aiStatus,omitempty"`
}

type VocabularyCard struct {
	ID               uint   `json:"id"`
	Word             string `json:"word"`
	Slug             string `json:"slug"`
	IPA              string `json:"ipa"`
	AudioURL         string `json:"audioUrl"`
	PartOfSpeech     string `json:"partOfSpeech"`
	MeaningVI        string `json:"meaningVi"`
	MeaningEN        string `json:"meaningEn"`
	ExampleSentence  string `json:"exampleSentence"`
	ExampleMeaningVI string `json:"exampleMeaningVi"`
}

type CompleteActivityResponse struct {
	Activity       ActivityItem `json:"activity"`
	Score          int          `json:"score"`
	Passed         bool         `json:"passed"`
	CorrectCount   int          `json:"correctCount"`
	AnsweredCount  int          `json:"answeredCount"`
	TotalCount     int          `json:"totalCount"`
	XPAwarded      int          `json:"xpAwarded"`
	NextActivityID *uint        `json:"nextActivityId"`
	Lesson         LessonResult `json:"lesson"`
}

type LessonResult struct {
	ID             uint                `json:"id"`
	Status         models.LessonStatus `json:"status"`
	PassedRequired int                 `json:"passedRequired"`
	TotalRequired  int                 `json:"totalRequired"`
	// JustCompleted is true only for the activity that finished the lesson.
	JustCompleted bool `json:"justCompleted"`
	Score         *int `json:"score"`
	XPAwarded     int  `json:"xpAwarded"`
}

// GrammarPointCard is a grammar theory card (also used by the grammar module).
type GrammarPointCard struct {
	ID             uint            `json:"id"`
	Slug           string          `json:"slug"`
	CEFRLevel      string          `json:"cefrLevel"`
	Title          string          `json:"title"`
	TitleVI        string          `json:"titleVi"`
	SummaryVI      string          `json:"summaryVi"`
	Formula        string          `json:"formula"`
	Examples       json.RawMessage `json:"examples" swaggertype:"array,object"`
	CommonMistakes json.RawMessage `json:"commonMistakes" swaggertype:"array,object"`
}

func ToGrammarPointCard(point models.GrammarPoint) GrammarPointCard {
	return GrammarPointCard{
		ID:             point.ID,
		Slug:           point.Slug,
		CEFRLevel:      point.CEFRLevel,
		Title:          point.Title,
		TitleVI:        point.TitleVI,
		SummaryVI:      point.SummaryVI,
		Formula:        point.Formula,
		Examples:       jsonList(point.ExamplesJSON),
		CommonMistakes: jsonList(point.CommonMistakesJSON),
	}
}

func jsonList(raw []byte) json.RawMessage {
	if len(raw) == 0 || string(raw) == "null" {
		return json.RawMessage("[]")
	}
	return json.RawMessage(raw)
}

// PhonemeCard is a sound card (also used by the phoneme module).
type PhonemeCard struct {
	ID           uint            `json:"id"`
	Symbol       string          `json:"symbol"`
	Kind         string          `json:"kind"`
	Voiced       *bool           `json:"voiced"`
	NameVI       string          `json:"nameVi"`
	ExampleWords json.RawMessage `json:"exampleWords" swaggertype:"array,object"`
	MinimalPairs json.RawMessage `json:"minimalPairs" swaggertype:"array,object"`
	MouthVI      string          `json:"mouthVi"`
	MistakesVI   string          `json:"mistakesVi"`
	AudioURL     string          `json:"audioUrl"`
	ImageURL     *string         `json:"imageUrl"`
	OrderIndex   int             `json:"orderIndex"`
}

func ToPhonemeCard(phoneme models.Phoneme) PhonemeCard {
	return PhonemeCard{
		ID:           phoneme.ID,
		Symbol:       phoneme.Symbol,
		Kind:         string(phoneme.Kind),
		Voiced:       phoneme.Voiced,
		NameVI:       phoneme.NameVI,
		ExampleWords: jsonList(phoneme.ExampleWordsJSON),
		MinimalPairs: jsonList(phoneme.MinimalPairsJSON),
		MouthVI:      phoneme.MouthVI,
		MistakesVI:   phoneme.MistakesVI,
		AudioURL:     phoneme.AudioURL,
		ImageURL:     phoneme.ImageURL,
		OrderIndex:   phoneme.OrderIndex,
	}
}

// PassageCard is a listening/reading text (also used by the passage module).
type PassageCard struct {
	ID        uint            `json:"id"`
	Slug      string          `json:"slug"`
	CEFRLevel string          `json:"cefrLevel"`
	Genre     string          `json:"genre"`
	Title     string          `json:"title"`
	BodyText  string          `json:"bodyText"`
	AudioURL  *string         `json:"audioUrl"`
	Segments  json.RawMessage `json:"segments" swaggertype:"array,object"`
	Glossary  json.RawMessage `json:"glossary" swaggertype:"array,object"`
	Source    string          `json:"source"`
	WordCount int             `json:"wordCount"`
}

func ToPassageCard(passage models.Passage) PassageCard {
	return PassageCard{
		ID:        passage.ID,
		Slug:      passage.Slug,
		CEFRLevel: passage.CEFRLevel,
		Genre:     string(passage.Genre),
		Title:     passage.Title,
		BodyText:  passage.BodyText,
		AudioURL:  passage.AudioURL,
		Segments:  jsonList(passage.SegmentsJSON),
		Glossary:  jsonList(passage.GlossaryJSON),
		Source:    passage.Source,
		WordCount: passage.WordCount,
	}
}

// ScenarioCard is the role-play brief shown in a SPEAKING activity.
type ScenarioCard struct {
	ID             uint     `json:"id"`
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Role           string   `json:"role"`
	Task           string   `json:"task"`
	TargetPatterns []string `json:"targetPatterns"`
	MaxTurns       int      `json:"maxTurns"`
	PassageID      *uint    `json:"passageId"`
}

// SpeakingSummary is the learner's latest finished role-play of an activity.
type SpeakingSummary struct {
	SessionID     uint       `json:"sessionId"`
	Score         *int       `json:"score"`
	TaskCompleted bool       `json:"taskCompleted"`
	FinishedAt    *time.Time `json:"finishedAt"`
}
