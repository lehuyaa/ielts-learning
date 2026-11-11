package contentimport

import "encoding/json"

// Bundle is one content file (docs/foundation-spec.md §8.2). Everything is
// upserted by slug; exercises are matched by their position in the activity.
type Bundle struct {
	Version    int          `json:"version"`
	CourseSlug string       `json:"courseSlug"`
	CEFRLevel  string       `json:"cefrLevel"`
	Topics     []TopicInput `json:"topics"`

	GrammarPoints []GrammarPointInput `json:"grammarPoints"`
	Phonemes      []PhonemeInput      `json:"phonemes"`

	Passages []PassageInput `json:"passages"`

	Scenarios []ScenarioInput `json:"scenarios"`
}

// GrammarPointInput is a grammar theory card (docs/foundation-spec.md §3.3).
type GrammarPointInput struct {
	Slug           string          `json:"slug"`
	CEFRLevel      string          `json:"cefrLevel"`
	Title          string          `json:"title"`
	TitleVI        string          `json:"titleVi"`
	SummaryVI      string          `json:"summaryVi"`
	Formula        string          `json:"formula"`
	Examples       json.RawMessage `json:"examples" swaggertype:"array,object"`
	CommonMistakes json.RawMessage `json:"commonMistakes" swaggertype:"array,object"`
	OrderIndex     int             `json:"orderIndex"`
}

// PhonemeInput is a sound card (docs/foundation-spec.md §3.3), upserted by symbol.
type PhonemeInput struct {
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

// PassageInput is a text/recording for LISTENING and READING (spec §3.3),
// upserted by slug.
type PassageInput struct {
	Slug      string          `json:"slug"`
	CEFRLevel string          `json:"cefrLevel"`
	Genre     string          `json:"genre"`
	Title     string          `json:"title"`
	BodyText  string          `json:"bodyText"`
	AudioURL  string          `json:"audioUrl"`
	Segments  json.RawMessage `json:"segments" swaggertype:"array,object"`
	Glossary  json.RawMessage `json:"glossary" swaggertype:"array,object"`
	Source    string          `json:"source"`
}

// ScenarioInput is a system role-play scenario (spec §3.1), upserted by slug.
type ScenarioInput struct {
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	CEFRLevel        string   `json:"cefrLevel"`
	Role             string   `json:"role"`
	Task             string   `json:"task"`
	SituationContext string   `json:"situationContext"`
	TargetPatterns   []string `json:"targetPatterns"`
	MaxTurns         int      `json:"maxTurns"`
	OpeningLine      string   `json:"openingLine"`
	// PassageSlug is the model dialogue (shared with the listening activity).
	PassageSlug string `json:"passageSlug"`
	Group       string `json:"group"`
	Duration    string `json:"duration"`
	OrderIndex  int    `json:"orderIndex"`
}

type TopicInput struct {
	Slug        string        `json:"slug"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	Emoji       string        `json:"emoji"`
	Color       string        `json:"color"`
	Track       string        `json:"track"`
	OrderIndex  int           `json:"orderIndex"`
	Lessons     []LessonInput `json:"lessons"`
}

type LessonInput struct {
	Slug             string          `json:"slug"`
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	EstimatedMinutes int             `json:"estimatedMinutes"`
	XPReward         int             `json:"xpReward"`
	RequiredScore    int             `json:"requiredScore"`
	OrderIndex       int             `json:"orderIndex"`
	Activities       []ActivityInput `json:"activities"`
}

type ActivityInput struct {
	Slug           string          `json:"slug"`
	Type           string          `json:"type"`
	Title          string          `json:"title"`
	InstructionsVI string          `json:"instructionsVi"`
	Content        json.RawMessage `json:"content,omitempty" swaggertype:"object"`
	// GrammarPointSlug links a GRAMMAR activity to its theory card; it may
	// also be given as content.grammarPointSlug.
	GrammarPointSlug string `json:"grammarPointSlug"`
	// PassageSlug links a LISTENING / READING activity to its passage; it
	// may also be given as content.passageSlug.
	PassageSlug string `json:"passageSlug"`
	// ScenarioSlug links a SPEAKING activity to its role-play; it may also
	// be given as content.scenarioSlug.
	ScenarioSlug string          `json:"scenarioSlug"`
	PassScore    *int            `json:"passScore"`
	XPReward     *int            `json:"xpReward"`
	IsRequired   *bool           `json:"isRequired"`
	Status       string          `json:"status"`
	Exercises    []ExerciseInput `json:"exercises"`
}

type ExerciseInput struct {
	Type          string          `json:"type"`
	Prompt        string          `json:"prompt"`
	Payload       json.RawMessage `json:"payload" swaggertype:"object"`
	Answer        json.RawMessage `json:"answer" swaggertype:"object"`
	ExplanationVI string          `json:"explanationVi"`
	SkillTag      string          `json:"skillTag"`
	Points        int             `json:"points"`
}
