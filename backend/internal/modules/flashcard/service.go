package flashcard

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"ielts-learning/backend/internal/models"
)

const flashcardReviewXP = 5

var (
	ErrLessonLocked  = errors.New("lesson locked")
	ErrInvalidRating = errors.New("invalid rating")
)

type DueReviewQuery struct {
	Limit   int
	TopicID *uint
}

type Service struct {
	repository Repository
	now        func() time.Time
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
		now:        time.Now,
	}
}

func (s Service) GetLessonFlashcards(userID uint, lessonID uint) (LessonFlashcardsResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return LessonFlashcardsResponse{}, err
	}

	context, err := s.loadLessonContext(userID, lessonID)
	if err != nil {
		return LessonFlashcardsResponse{}, err
	}
	if context.status == models.LessonStatusLocked {
		return LessonFlashcardsResponse{}, ErrLessonLocked
	}

	cards := make([]FlashcardResponse, 0, len(context.lessonVocabularies))
	for _, row := range context.lessonVocabularies {
		progress, hasProgress := context.progressByVocabulary[row.VocabularyID]
		cards = append(cards, flashcardResponse(row.ID, row.Vocabulary, progress, hasProgress, context.lesson.Topic.Title))
	}

	done := completedReviewCount(context.lessonVocabularies, context.progressByVocabulary)
	total := len(cards)

	return LessonFlashcardsResponse{
		Lesson: LessonResponse{
			ID:         context.lesson.ID,
			Title:      context.lesson.Title,
			TopicID:    context.lesson.Topic.ID,
			TopicTitle: context.lesson.Topic.Title,
			BandLabel:  bandLabel(context.lesson.Topic.BandLevel),
		},
		Progress: SessionProgress{
			Done:      done,
			Remaining: maxInt(total-done, 0),
			Total:     total,
		},
		Cards: cards,
	}, nil
}

func (s Service) GetDueReviews(userID uint, query DueReviewQuery) (DueReviewsResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return DueReviewsResponse{}, err
	}

	query.Limit = normalizeLimit(query.Limit)
	progressRows, err := s.repository.FindDueProgress(userID, query, s.now().UTC())
	if err != nil {
		return DueReviewsResponse{}, err
	}

	vocabularyIDs := make([]uint, 0, len(progressRows))
	for _, progress := range progressRows {
		vocabularyIDs = append(vocabularyIDs, progress.VocabularyID)
	}

	contexts, err := s.repository.FindTopicContext(vocabularyIDs)
	if err != nil {
		return DueReviewsResponse{}, err
	}

	cards := make([]FlashcardResponse, 0, len(progressRows))
	for _, progress := range progressRows {
		context := contexts[progress.VocabularyID]
		cards = append(cards, flashcardResponse(progress.Vocabulary.ID, progress.Vocabulary, progress, true, context.TopicTitle))
	}

	return DueReviewsResponse{
		Cards: cards,
		Count: len(cards),
	}, nil
}

func (s Service) Review(userID uint, request ReviewRequest) (ReviewResponse, error) {
	if request.VocabularyID == 0 {
		return ReviewResponse{}, ErrVocabularyNotFound
	}
	if !validRating(request.Rating) {
		return ReviewResponse{}, ErrInvalidRating
	}

	user, err := s.repository.FindUser(userID)
	if err != nil {
		return ReviewResponse{}, err
	}

	now := s.now().UTC()
	progress, updatedUser, xpAwarded, err := s.repository.SaveReview(user, request.VocabularyID, request.LessonID, now, func(progress models.UserVocabularyProgress, _ bool) ReviewUpdate {
		return ReviewUpdate{
			Progress:  applyRating(progress, request.Rating, now),
			XPAwarded: flashcardReviewXP,
		}
	})
	if err != nil {
		return ReviewResponse{}, err
	}

	lastRating := models.FlashcardRating("")
	if progress.LastRating != nil {
		lastRating = *progress.LastRating
	}

	return ReviewResponse{
		VocabularyID:   progress.VocabularyID,
		Status:         progress.Status,
		ReviewCount:    progress.ReviewCount,
		CorrectCount:   progress.CorrectCount,
		WrongCount:     progress.WrongCount,
		LastRating:     lastRating,
		LastReviewedAt: progress.LastReviewedAt,
		NextReviewAt:   progress.NextReviewAt,
		MasteryScore:   masteryScore(progress, true),
		XPAwarded:      xpAwarded,
		TotalXP:        updatedUser.TotalXP,
	}, nil
}

type lessonContext struct {
	lesson               models.Lesson
	lessonVocabularies   []models.LessonVocabulary
	progressByVocabulary map[uint]models.UserVocabularyProgress
	status               models.LessonStatus
}

func (s Service) loadLessonContext(userID uint, lessonID uint) (lessonContext, error) {
	lesson, err := s.repository.FindLesson(lessonID)
	if err != nil {
		return lessonContext{}, err
	}

	topicLessons, err := s.repository.FindTopicLessons(lesson.TopicID)
	if err != nil {
		return lessonContext{}, err
	}

	lessonIDs := make([]uint, 0, len(topicLessons))
	for _, topicLesson := range topicLessons {
		lessonIDs = append(lessonIDs, topicLesson.ID)
	}

	lessonProgressRows, err := s.repository.FindLessonProgress(userID, lessonIDs)
	if err != nil {
		return lessonContext{}, err
	}

	lessonVocabularies, err := s.repository.FindLessonVocabularies(lesson.ID)
	if err != nil {
		return lessonContext{}, err
	}

	vocabularyIDs := make([]uint, 0, len(lessonVocabularies))
	for _, row := range lessonVocabularies {
		vocabularyIDs = append(vocabularyIDs, row.VocabularyID)
	}

	progressRows, err := s.repository.FindVocabularyProgress(userID, vocabularyIDs)
	if err != nil {
		return lessonContext{}, err
	}

	return lessonContext{
		lesson:               lesson,
		lessonVocabularies:   lessonVocabularies,
		progressByVocabulary: mapProgressByVocabulary(progressRows),
		status:               lessonStatus(lesson.ID, topicLessons, mapLessonProgressByID(lessonProgressRows)),
	}, nil
}

func applyRating(progress models.UserVocabularyProgress, rating models.FlashcardRating, now time.Time) models.UserVocabularyProgress {
	progress.ReviewCount++
	progress.LastRating = &rating
	progress.LastReviewedAt = &now

	if progress.FirstLearnedAt == nil {
		progress.FirstLearnedAt = &now
	}

	switch rating {
	case models.FlashcardRatingAgain:
		progress.WrongCount++
		progress.IntervalDays = 0
		progress.Status = models.VocabularyStatusLearning
		progress.NextReviewAt = &now
	case models.FlashcardRatingHard:
		progress.CorrectCount++
		progress.IntervalDays = 1
		progress.Status = models.VocabularyStatusLearning
		next := now.AddDate(0, 0, 1)
		progress.NextReviewAt = &next
	case models.FlashcardRatingGood:
		progress.CorrectCount++
		progress.IntervalDays = 3
		progress.Status = models.VocabularyStatusReview
		if progress.LearnedAt == nil {
			progress.LearnedAt = &now
		}
		next := now.AddDate(0, 0, 3)
		progress.NextReviewAt = &next
	case models.FlashcardRatingEasy:
		progress.CorrectCount++
		progress.IntervalDays = 7
		if progress.LearnedAt == nil {
			progress.LearnedAt = &now
		}
		if progress.CorrectCount >= 4 {
			progress.Status = models.VocabularyStatusMastered
			if progress.MasteredAt == nil {
				progress.MasteredAt = &now
			}
		} else {
			progress.Status = models.VocabularyStatusReview
		}
		next := now.AddDate(0, 0, 7)
		progress.NextReviewAt = &next
	}

	return progress
}

func flashcardResponse(
	id uint,
	vocabulary models.Vocabulary,
	progress models.UserVocabularyProgress,
	hasProgress bool,
	topicTitle string,
) FlashcardResponse {
	status := models.VocabularyStatusNew
	if hasProgress {
		status = progress.Status
	}

	return FlashcardResponse{
		ID:              id,
		VocabularyID:    vocabulary.ID,
		Word:            vocabulary.Word,
		Slug:            vocabulary.Slug,
		IPA:             vocabulary.IPA,
		AudioURL:        vocabulary.AudioURL,
		PartOfSpeech:    vocabulary.PartOfSpeech,
		MeaningVI:       vocabulary.MeaningVI,
		MeaningEN:       firstNonEmpty(vocabulary.MeaningEN, vocabulary.ShortDefinition, vocabulary.MeaningVI),
		ExampleSentence: vocabulary.ExampleSentence,
		Synonyms:        stringSlice(vocabulary.SynonymsJSON),
		Collocations:    stringSlice(vocabulary.CollocationsJSON),
		IELTSUsage:      vocabulary.IELTSUsage,
		Difficulty:      vocabulary.Difficulty,
		Band:            vocabulary.TargetBand,
		TargetBand:      vocabulary.TargetBand,
		TopicTitle:      topicTitle,
		Status:          status,
		ReviewCount:     progress.ReviewCount,
		NextReviewAt:    progress.NextReviewAt,
	}
}

func lessonStatus(
	lessonID uint,
	topicLessons []models.Lesson,
	progressByLesson map[uint]models.UserLessonProgress,
) models.LessonStatus {
	if progress, ok := progressByLesson[lessonID]; ok {
		return progress.Status
	}

	previousLessonCompleted := true
	for _, lesson := range topicLessons {
		if lesson.ID == lessonID {
			if previousLessonCompleted {
				return models.LessonStatusUnlocked
			}
			return models.LessonStatusLocked
		}

		progress, ok := progressByLesson[lesson.ID]
		previousLessonCompleted = ok && progress.Status == models.LessonStatusCompleted
	}

	return models.LessonStatusLocked
}

func completedReviewCount(
	lessonVocabularies []models.LessonVocabulary,
	progressByVocabulary map[uint]models.UserVocabularyProgress,
) int {
	done := 0
	for _, row := range lessonVocabularies {
		progress, ok := progressByVocabulary[row.VocabularyID]
		if ok && progress.Status != models.VocabularyStatusNew {
			done++
		}
	}
	return done
}

func mapLessonProgressByID(rows []models.UserLessonProgress) map[uint]models.UserLessonProgress {
	result := make(map[uint]models.UserLessonProgress, len(rows))
	for _, row := range rows {
		result[row.LessonID] = row
	}
	return result
}

func mapProgressByVocabulary(rows []models.UserVocabularyProgress) map[uint]models.UserVocabularyProgress {
	result := make(map[uint]models.UserVocabularyProgress, len(rows))
	for _, row := range rows {
		result[row.VocabularyID] = row
	}
	return result
}

func validRating(rating models.FlashcardRating) bool {
	switch rating {
	case models.FlashcardRatingAgain, models.FlashcardRatingHard, models.FlashcardRatingGood, models.FlashcardRatingEasy:
		return true
	default:
		return false
	}
}

func stringSlice(data []byte) []string {
	if len(data) == 0 {
		return []string{}
	}

	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return []string{}
	}

	return values
}

func masteryScore(progress models.UserVocabularyProgress, hasProgress bool) int {
	if !hasProgress {
		return 0
	}
	if progress.Status == models.VocabularyStatusMastered {
		return 100
	}

	totalAnswers := progress.CorrectCount + progress.WrongCount
	if totalAnswers > 0 {
		return clampPercentage((progress.CorrectCount * 100) / totalAnswers)
	}

	switch progress.Status {
	case models.VocabularyStatusReview:
		return 70
	case models.VocabularyStatusLearning:
		return 40
	default:
		return 0
	}
}

func normalizeLimit(limit int) int {
	switch {
	case limit <= 0:
		return 20
	case limit > 100:
		return 100
	default:
		return limit
	}
}

func bandLabel(bandLevel models.BandLevel) string {
	if bandLevel.Title != "" {
		return bandLevel.Title
	}
	if bandLevel.BandScore > 0 {
		return fmt.Sprintf("Band %.1f", bandLevel.BandScore)
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func clampPercentage(value int) int {
	switch {
	case value < 0:
		return 0
	case value > 100:
		return 100
	default:
		return value
	}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
