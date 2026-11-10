package contentimport

import "encoding/json"

// Bundle is one content file (docs/foundation-spec.md §8.2). Everything is
// upserted by slug; exercises are matched by their position in the activity.
type Bundle struct {
	Version    int          `json:"version"`
	CourseSlug string       `json:"courseSlug"`
	CEFRLevel  string       `json:"cefrLevel"`
	Topics     []TopicInput `json:"topics"`

	// Reserved for later tasks; a bundle that uses them gets a warning.
	GrammarPoints json.RawMessage `json:"grammarPoints,omitempty" swaggertype:"array,object"`
	Phonemes      json.RawMessage `json:"phonemes,omitempty" swaggertype:"array,object"`
	Passages      json.RawMessage `json:"passages,omitempty" swaggertype:"array,object"`
	Scenarios     json.RawMessage `json:"scenarios,omitempty" swaggertype:"array,object"`
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
	PassScore      *int            `json:"passScore"`
	XPReward       *int            `json:"xpReward"`
	IsRequired     *bool           `json:"isRequired"`
	Status         string          `json:"status"`
	Exercises      []ExerciseInput `json:"exercises"`
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
