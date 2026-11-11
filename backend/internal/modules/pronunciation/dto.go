package pronunciation

// StatusResponse tells the app whether recordings can be scored.
type StatusResponse struct {
	Enabled    bool `json:"enabled"`
	MaxSeconds int  `json:"maxSeconds"`
	DailyLimit int  `json:"dailyLimit"`
	// Provider is AZURE or OPENAI (empty when disabled). Phonemes is false
	// for OPENAI, which scores whole words only.
	Provider string `json:"provider,omitempty"`
	Phonemes bool   `json:"phonemes"`
}

// AssessResponse is the score of one recording (spec §7).
type AssessResponse struct {
	AttemptID      uint         `json:"attemptId"`
	ReferenceText  string       `json:"referenceText"`
	RecognizedText string       `json:"recognizedText"`
	Accuracy       int          `json:"accuracy"`
	Fluency        int          `json:"fluency"`
	Completeness   int          `json:"completeness"`
	Pronunciation  int          `json:"pronunciation"`
	Words          []WordResult `json:"words"`
	// MinAccuracy and Passed are set for a SPEAK exercise.
	MinAccuracy *int  `json:"minAccuracy,omitempty"`
	Passed      *bool `json:"passed,omitempty"`
}

// WordResult is one word of the reference text. ErrorType is None,
// Mispronunciation, Omission or Insertion.
type WordResult struct {
	Word      string          `json:"word"`
	Accuracy  int             `json:"accuracy"`
	ErrorType string          `json:"errorType"`
	Phonemes  []PhonemeResult `json:"phonemes"`
}

// PhonemeResult is one sound of a word: Display is in our IPA notation,
// Symbols are the sound cards it counts for.
type PhonemeResult struct {
	Azure    string   `json:"azure"`
	Display  string   `json:"display"`
	Symbols  []string `json:"symbols"`
	Accuracy int      `json:"accuracy"`
}
