package passage

import (
	"encoding/json"
	"strings"
	"time"

	"ielts-learning/backend/internal/modules/lessonactivity"
)

type Service struct {
	repository Repository
	now        func() time.Time
}

func NewService(repository Repository) Service {
	return Service{repository: repository, now: time.Now}
}

func (s Service) Get(id uint) (lessonactivity.PassageCard, error) {
	passage, err := s.repository.FindPassage(id)
	if err != nil {
		return lessonactivity.PassageCard{}, err
	}
	return lessonactivity.ToPassageCard(passage), nil
}

type glossaryEntry struct {
	Word           string `json:"word"`
	VocabularySlug string `json:"vocabularySlug"`
	MeaningVI      string `json:"meaningVi"`
	IPA            string `json:"ipa"`
}

// Lookup finds a tapped word: the passage glossary first (its meaning and
// linked word), then the vocabulary by the word's base forms.
func (s Service) Lookup(userID uint, word string, passageID uint) (WordLookupResponse, error) {
	result := WordLookupResponse{Word: word}
	candidates := wordCandidates(word)

	var entry *glossaryEntry
	if passageID != 0 {
		passage, err := s.repository.FindPassage(passageID)
		if err != nil {
			return WordLookupResponse{}, err
		}
		var glossary []glossaryEntry
		_ = json.Unmarshal(passage.GlossaryJSON, &glossary)
		for i := range glossary {
			for _, candidate := range candidates {
				if strings.EqualFold(glossary[i].Word, candidate) {
					entry = &glossary[i]
					break
				}
			}
			if entry != nil {
				break
			}
		}
	}

	slug := ""
	if entry != nil {
		result.Found = true
		result.Note = entry.MeaningVI
		result.IPA = entry.IPA
		slug = entry.VocabularySlug
	}

	vocabulary, ok, err := s.repository.FindWord(candidates, slug)
	if err != nil {
		return WordLookupResponse{}, err
	}
	if !ok && slug != "" {
		vocabulary, ok, err = s.repository.FindWord(candidates, "")
		if err != nil {
			return WordLookupResponse{}, err
		}
	}
	if ok {
		inFlashcards, err := s.repository.HasProgress(userID, vocabulary.ID)
		if err != nil {
			return WordLookupResponse{}, err
		}
		result.Found = true
		if result.IPA == "" {
			result.IPA = vocabulary.IPA
		}
		result.Vocabulary = &WordResult{
			ID:           vocabulary.ID,
			Word:         vocabulary.Word,
			Slug:         vocabulary.Slug,
			IPA:          vocabulary.IPA,
			PartOfSpeech: vocabulary.PartOfSpeech,
			MeaningVI:    vocabulary.MeaningVI,
			AudioURL:     vocabulary.AudioURL,
			InFlashcards: inFlashcards,
		}
	}
	return result, nil
}

// AddFlashcard puts a word into the learner's reviews (status NEW, due now).
func (s Service) AddFlashcard(userID uint, vocabularyID uint) (AddFlashcardResponse, error) {
	exists, err := s.repository.VocabularyExists(vocabularyID)
	if err != nil {
		return AddFlashcardResponse{}, err
	}
	if !exists {
		return AddFlashcardResponse{}, ErrVocabularyNotFound
	}
	row, added, err := s.repository.AddToReviews(userID, vocabularyID, s.now())
	if err != nil {
		return AddFlashcardResponse{}, err
	}
	return AddFlashcardResponse{VocabularyID: vocabularyID, Added: added, Status: string(row.Status)}, nil
}
