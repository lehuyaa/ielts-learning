package contentimport

const (
	KindTopic    = "TOPIC"
	KindLesson   = "LESSON"
	KindActivity = "ACTIVITY"
	KindExercise = "EXERCISE"

	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionSkip   = "SKIP"
)

type ImportResult struct {
	Committed bool          `json:"committed"`
	Summary   ImportSummary `json:"summary"`
	// Errors that stop the whole bundle (bad course or level).
	BundleErrors []string     `json:"bundleErrors"`
	Warnings     []string     `json:"warnings"`
	Items        []ItemResult `json:"items"`
}

type ImportSummary struct {
	Topics     int `json:"topics"`
	Lessons    int `json:"lessons"`
	Activities int `json:"activities"`
	Exercises  int `json:"exercises"`
	Created    int `json:"created"`
	Updated    int `json:"updated"`
	Invalid    int `json:"invalid"`
}

type ItemResult struct {
	Kind   string   `json:"kind"`
	Path   string   `json:"path"`
	Title  string   `json:"title"`
	Action string   `json:"action"`
	Status string   `json:"status,omitempty"`
	Errors []string `json:"errors"`
}
