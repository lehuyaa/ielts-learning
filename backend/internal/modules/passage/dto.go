package passage

// WordLookupResponse is what tap-to-translate shows for one tapped word.
type WordLookupResponse struct {
	Word  string `json:"word"`
	Found bool   `json:"found"`
	// Note is the passage glossary's meaning for this word, when given.
	Note       string      `json:"note"`
	IPA        string      `json:"ipa"`
	Vocabulary *WordResult `json:"vocabulary"`
}

type WordResult struct {
	ID           uint   `json:"id"`
	Word         string `json:"word"`
	Slug         string `json:"slug"`
	IPA          string `json:"ipa"`
	PartOfSpeech string `json:"partOfSpeech"`
	MeaningVI    string `json:"meaningVi"`
	AudioURL     string `json:"audioUrl"`
	// InFlashcards is true when the word is already in the learner's reviews.
	InFlashcards bool `json:"inFlashcards"`
}

type AddFlashcardResponse struct {
	VocabularyID uint   `json:"vocabularyId"`
	Added        bool   `json:"added"`
	Status       string `json:"status"`
}
