package topic

import (
	"errors"

	"ielts-learning/backend/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) Get(userID uint, topicID uint) (Response, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return Response{}, err
	}

	topic, err := s.repository.FindTopic(topicID)
	if err != nil {
		return Response{}, err
	}

	lessons, err := s.repository.FindLessons(topicID)
	if err != nil {
		return Response{}, err
	}

	lessonIDs := make([]uint, 0, len(lessons))
	for _, lesson := range lessons {
		lessonIDs = append(lessonIDs, lesson.ID)
	}

	progressRows, err := s.repository.FindLessonProgress(userID, lessonIDs)
	if err != nil {
		return Response{}, err
	}

	wordCounts, err := s.repository.CountVocabularyByLesson(lessonIDs)
	if err != nil {
		return Response{}, err
	}

	return s.toResponse(topic, lessons, progressRows, wordCounts), nil
}

func (s Service) toResponse(
	topic models.Topic,
	lessons []models.Lesson,
	progressRows []models.UserLessonProgress,
	wordCounts map[uint]int,
) Response {
	progressByLesson := mapProgressByLesson(progressRows)

	totalXP := 0
	completedLessons := 0
	previousLessonCompleted := true
	lessonResponses := make([]LessonResponse, 0, len(lessons))

	for _, lesson := range lessons {
		progress, hasProgress := progressByLesson[lesson.ID]
		status := lessonStatus(hasProgress, progress, previousLessonCompleted)
		if status == models.LessonStatusCompleted {
			completedLessons++
		}

		totalXP += lesson.XPReward

		lessonResponses = append(lessonResponses, LessonResponse{
			ID:                 lesson.ID,
			Title:              lesson.Title,
			Slug:               lesson.Slug,
			Description:        lesson.Description,
			WordCount:          wordCounts[lesson.ID],
			EstimatedMinutes:   lesson.EstimatedMinutes,
			XPReward:           lesson.XPReward,
			Status:             status,
			ProgressPercentage: lessonProgressPercentage(status, progress),
			LockedReason:       lockedReason(status),
		})

		previousLessonCompleted = status == models.LessonStatusCompleted
	}

	return Response{
		Topic: TopicResponse{
			ID:          topic.ID,
			Title:       topic.Title,
			Slug:        topic.Slug,
			Icon:        topic.Icon,
			Emoji:       topic.Emoji,
			Color:       topic.Color,
			Description: topic.Description,
		},
		BandLevel: BandLevelResponse{
			ID:        topic.BandLevel.ID,
			BandScore: topic.BandLevel.BandScore,
			Title:     topic.BandLevel.Title,
		},
		Summary: SummaryResponse{
			ProgressPercentage: percentage(completedLessons, len(lessons)),
			CompletedLessons:   completedLessons,
			TotalLessons:       len(lessons),
			TotalXP:            totalXP,
		},
		Lessons: lessonResponses,
	}
}

func mapProgressByLesson(progressRows []models.UserLessonProgress) map[uint]models.UserLessonProgress {
	result := make(map[uint]models.UserLessonProgress)
	for _, progress := range progressRows {
		result[progress.LessonID] = progress
	}
	return result
}

func lessonStatus(hasProgress bool, progress models.UserLessonProgress, previousLessonCompleted bool) models.LessonStatus {
	if hasProgress {
		return progress.Status
	}
	if previousLessonCompleted {
		return models.LessonStatusUnlocked
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

func lockedReason(status models.LessonStatus) *string {
	if status != models.LessonStatusLocked {
		return nil
	}

	reason := "Complete previous lessons to unlock this lesson."
	return &reason
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
