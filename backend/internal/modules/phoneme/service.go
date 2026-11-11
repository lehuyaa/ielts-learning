package phoneme

import (
	"encoding/json"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/lessonactivity"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) List(userID uint) (ListResponse, error) {
	phonemes, err := s.repository.FindAll()
	if err != nil {
		return ListResponse{}, err
	}
	progress, err := s.repository.FindProgress(userID)
	if err != nil {
		return ListResponse{}, err
	}

	items := make([]PhonemeItem, 0, len(phonemes))
	for _, phoneme := range phonemes {
		row, ok := progress[phoneme.ID]
		items = append(items, PhonemeItem{
			ID:          phoneme.ID,
			Symbol:      phoneme.Symbol,
			Kind:        string(phoneme.Kind),
			Voiced:      phoneme.Voiced,
			NameVI:      phoneme.NameVI,
			AudioURL:    phoneme.AudioURL,
			ExampleWord: firstExample(phoneme.ExampleWordsJSON),
			OrderIndex:  phoneme.OrderIndex,
			Progress:    toProgress(row, ok),
		})
	}
	topics, err := s.repository.FindPracticeTopics()
	if err != nil {
		return ListResponse{}, err
	}
	return ListResponse{Items: items, PracticeTopics: topics}, nil
}

func (s Service) Get(userID uint, symbol string) (DetailResponse, error) {
	phonemes, err := s.repository.FindAll()
	if err != nil {
		return DetailResponse{}, err
	}
	index := -1
	for i, phoneme := range phonemes {
		if phoneme.Symbol == symbol {
			index = i
			break
		}
	}
	if index < 0 {
		return DetailResponse{}, ErrPhonemeNotFound
	}
	phoneme := phonemes[index]

	progress, err := s.repository.FindProgress(userID)
	if err != nil {
		return DetailResponse{}, err
	}
	lessons, err := s.repository.FindLessons(symbol)
	if err != nil {
		return DetailResponse{}, err
	}

	row, ok := progress[phoneme.ID]
	detail := DetailResponse{
		Phoneme:  lessonactivity.ToPhonemeCard(phoneme),
		Progress: toProgress(row, ok),
		Lessons:  lessons,
	}
	if index > 0 {
		detail.Previous = &phonemes[index-1].Symbol
	}
	if index < len(phonemes)-1 {
		detail.Next = &phonemes[index+1].Symbol
	}
	return detail, nil
}

func firstExample(raw []byte) *ExampleWord {
	var words []ExampleWord
	if len(raw) == 0 || json.Unmarshal(raw, &words) != nil || len(words) == 0 {
		return nil
	}
	return &words[0]
}

func toProgress(row models.UserPhonemeProgress, ok bool) ProgressResponse {
	if !ok {
		return ProgressResponse{}
	}
	return ProgressResponse{
		Accuracy:        row.Accuracy,
		Attempts:        row.Attempts,
		Practised:       row.Attempts > 0,
		LastPracticedAt: row.LastPracticed,
	}
}
