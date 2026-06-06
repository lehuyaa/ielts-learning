package lesson

import (
	"errors"
	"time"

	"ielts-learning/backend/internal/models"
)

var ErrLessonLocked = errors.New("lesson locked")

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

func (s Service) Get(userID uint, lessonID uint) (DetailResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return DetailResponse{}, err
	}

	context, err := s.loadContext(userID, lessonID)
	if err != nil {
		return DetailResponse{}, err
	}

	return s.toDetailResponse(context), nil
}

func (s Service) Start(userID uint, lessonID uint) (StartResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return StartResponse{}, err
	}

	context, err := s.loadContext(userID, lessonID)
	if err != nil {
		return StartResponse{}, err
	}

	if context.status == models.LessonStatusLocked {
		return StartResponse{}, ErrLessonLocked
	}

	progress, err := s.repository.StartLesson(userID, context.lesson, len(context.vocabularies), s.now().UTC())
	if err != nil {
		return StartResponse{}, err
	}

	return StartResponse{
		LessonID:      context.lesson.ID,
		Status:        progress.Status,
		StartedAt:     progress.StartedAt,
		LastStudiedAt: progress.LastStudiedAt,
	}, nil
}

type detailContext struct {
	lesson             models.Lesson
	topicLessons       []models.Lesson
	vocabularies       []models.Vocabulary
	lessonProgressByID map[uint]models.UserLessonProgress
	vocabProgressByID  map[uint]models.UserVocabularyProgress
	status             models.LessonStatus
	progress           models.UserLessonProgress
	hasProgress        bool
}

func (s Service) loadContext(userID uint, lessonID uint) (detailContext, error) {
	lesson, err := s.repository.FindLesson(lessonID)
	if err != nil {
		return detailContext{}, err
	}

	topicLessons, err := s.repository.FindTopicLessons(lesson.TopicID)
	if err != nil {
		return detailContext{}, err
	}

	lessonIDs := make([]uint, 0, len(topicLessons))
	for _, topicLesson := range topicLessons {
		lessonIDs = append(lessonIDs, topicLesson.ID)
	}

	lessonProgressRows, err := s.repository.FindLessonProgress(userID, lessonIDs)
	if err != nil {
		return detailContext{}, err
	}
	lessonProgressByID := mapLessonProgressByID(lessonProgressRows)

	vocabularies, err := s.repository.FindLessonVocabularies(lesson.ID)
	if err != nil {
		return detailContext{}, err
	}

	vocabularyIDs := make([]uint, 0, len(vocabularies))
	for _, vocabulary := range vocabularies {
		vocabularyIDs = append(vocabularyIDs, vocabulary.ID)
	}

	vocabProgressRows, err := s.repository.FindVocabularyProgress(userID, vocabularyIDs)
	if err != nil {
		return detailContext{}, err
	}

	progress, hasProgress := lessonProgressByID[lesson.ID]
	status := lessonStatus(lesson.ID, topicLessons, lessonProgressByID)

	return detailContext{
		lesson:             lesson,
		topicLessons:       topicLessons,
		vocabularies:       vocabularies,
		lessonProgressByID: lessonProgressByID,
		vocabProgressByID:  mapVocabularyProgressByID(vocabProgressRows),
		status:             status,
		progress:           progress,
		hasProgress:        hasProgress,
	}, nil
}

func (s Service) toDetailResponse(context detailContext) DetailResponse {
	progressPercentage := lessonProgressPercentage(context.status, context.progress)
	wordCount := len(context.vocabularies)

	return DetailResponse{
		Lesson: LessonResponse{
			ID:                 context.lesson.ID,
			Title:              context.lesson.Title,
			Slug:               context.lesson.Slug,
			Description:        context.lesson.Description,
			RequiredScore:      context.lesson.RequiredScore,
			EstimatedMinutes:   context.lesson.EstimatedMinutes,
			XPReward:           context.lesson.XPReward,
			WordCount:          wordCount,
			OrderIndex:         context.lesson.OrderIndex,
			Status:             context.status,
			ProgressPercentage: progressPercentage,
			LockedReason:       lockedReason(context.status),
		},
		Topic: TopicResponse{
			ID:    context.lesson.Topic.ID,
			Title: context.lesson.Topic.Title,
			Slug:  context.lesson.Topic.Slug,
			Icon:  context.lesson.Topic.Icon,
			Emoji: context.lesson.Topic.Emoji,
			Color: context.lesson.Topic.Color,
		},
		BandLevel: BandLevelResponse{
			ID:        context.lesson.Topic.BandLevel.ID,
			BandScore: context.lesson.Topic.BandLevel.BandScore,
			Title:     context.lesson.Topic.BandLevel.Title,
		},
		Progress: ProgressResponse{
			Status:             context.status,
			Score:              context.progress.Score,
			BestScore:          context.progress.BestScore,
			BestXP:             context.progress.BestXP,
			WordsLearned:       context.progress.WordsLearned,
			TotalWords:         totalWords(context.progress.TotalWords, wordCount),
			ProgressPercentage: progressPercentage,
			StartedAt:          context.progress.StartedAt,
			CompletedAt:        context.progress.CompletedAt,
			LastStudiedAt:      context.progress.LastStudiedAt,
		},
		Vocabularies: vocabularyResponses(context.vocabularies, context.vocabProgressByID),
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

func lessonProgressPercentage(status models.LessonStatus, progress models.UserLessonProgress) int {
	switch status {
	case models.LessonStatusCompleted:
		return 100
	case models.LessonStatusInProgress:
		if progress.TotalWords > 0 {
			return percentage(progress.WordsLearned, progress.TotalWords)
		}
		if progress.BestScore != nil {
			return clampPercentage(*progress.BestScore)
		}
		if progress.Score != nil {
			return clampPercentage(*progress.Score)
		}
	}

	return 0
}

func vocabularyResponses(
	vocabularies []models.Vocabulary,
	progressByVocabulary map[uint]models.UserVocabularyProgress,
) []VocabularyResponse {
	responses := make([]VocabularyResponse, 0, len(vocabularies))
	for _, vocabulary := range vocabularies {
		progress, ok := progressByVocabulary[vocabulary.ID]
		status := models.VocabularyStatusNew
		if ok {
			status = progress.Status
		}

		responses = append(responses, VocabularyResponse{
			ID:              vocabulary.ID,
			Word:            vocabulary.Word,
			Slug:            vocabulary.Slug,
			IPA:             vocabulary.IPA,
			AudioURL:        vocabulary.AudioURL,
			PartOfSpeech:    vocabulary.PartOfSpeech,
			MeaningVI:       vocabulary.MeaningVI,
			MeaningEN:       vocabulary.MeaningEN,
			ShortDefinition: vocabulary.ShortDefinition,
			ExampleSentence: vocabulary.ExampleSentence,
			Difficulty:      vocabulary.Difficulty,
			TargetBand:      vocabulary.TargetBand,
			Status:          status,
			ReviewCount:     progress.ReviewCount,
			CorrectCount:    progress.CorrectCount,
			WrongCount:      progress.WrongCount,
			Learned:         status == models.VocabularyStatusLearning || status == models.VocabularyStatusReview || status == models.VocabularyStatusMastered,
			LearnedAt:       progress.LearnedAt,
			LastReviewedAt:  progress.LastReviewedAt,
			NextReviewAt:    progress.NextReviewAt,
		})
	}

	return responses
}

func mapLessonProgressByID(rows []models.UserLessonProgress) map[uint]models.UserLessonProgress {
	result := make(map[uint]models.UserLessonProgress, len(rows))
	for _, row := range rows {
		result[row.LessonID] = row
	}
	return result
}

func mapVocabularyProgressByID(rows []models.UserVocabularyProgress) map[uint]models.UserVocabularyProgress {
	result := make(map[uint]models.UserVocabularyProgress, len(rows))
	for _, row := range rows {
		result[row.VocabularyID] = row
	}
	return result
}

func lockedReason(status models.LessonStatus) *string {
	if status != models.LessonStatusLocked {
		return nil
	}

	reason := "Complete previous lessons to unlock this lesson."
	return &reason
}

func totalWords(progressTotal int, vocabularyCount int) int {
	if progressTotal > 0 {
		return progressTotal
	}
	return vocabularyCount
}

func percentage(completed int, total int) int {
	if total == 0 {
		return 0
	}
	return (completed * 100) / total
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
