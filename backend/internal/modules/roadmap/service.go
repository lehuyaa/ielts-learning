package roadmap

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

func (s Service) Get(userID uint) (Response, error) {
	user, err := s.repository.FindUser(userID)
	if err != nil {
		return Response{}, err
	}

	course, err := s.repository.FindPublishedCourse()
	if err != nil {
		return Response{}, err
	}

	bands, err := s.repository.FindBandLevels(course.ID)
	if err != nil {
		return Response{}, err
	}

	bandIDs := make([]uint, 0, len(bands))
	for _, band := range bands {
		bandIDs = append(bandIDs, band.ID)
	}

	topics, err := s.repository.FindTopics(bandIDs)
	if err != nil {
		return Response{}, err
	}

	topicIDs := make([]uint, 0, len(topics))
	for _, topic := range topics {
		topicIDs = append(topicIDs, topic.ID)
	}

	lessons, err := s.repository.FindLessons(topicIDs)
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

	wordsMastered, err := s.repository.CountMasteredWords(userID)
	if err != nil {
		return Response{}, err
	}

	return s.toResponse(course, user, bands, topics, lessons, progressRows, wordsMastered), nil
}

func (s Service) toResponse(
	course models.Course,
	user models.User,
	bands []models.BandLevel,
	topics []models.Topic,
	lessons []models.Lesson,
	progressRows []models.UserLessonProgress,
	wordsMastered int,
) Response {
	topicsByBand := groupTopicsByBand(topics)
	lessonsByTopic := groupLessonsByTopic(lessons)
	progressByLesson := mapProgressByLesson(progressRows)

	totalTopics := 0
	topicsCompleted := 0
	previousLessonCompleted := true

	bandResponses := make([]BandLevelResponse, 0, len(bands))
	for _, band := range bands {
		bandTopics := topicsByBand[band.ID]
		topicResponses := make([]TopicResponse, 0, len(bandTopics))
		bandCompletedTopics := 0

		for _, topic := range bandTopics {
			topicLessons := lessonsByTopic[topic.ID]
			lessonResponses := make([]LessonResponse, 0, len(topicLessons))
			topicCompletedLessons := 0
			topicStartedLessons := 0
			topicLockedLessons := 0

			for _, lesson := range topicLessons {
				progress, hasProgress := progressByLesson[lesson.ID]
				status := lessonStatus(hasProgress, progress, previousLessonCompleted)
				if status == models.LessonStatusCompleted {
					topicCompletedLessons++
					previousLessonCompleted = true
				} else {
					previousLessonCompleted = false
				}
				if status == models.LessonStatusInProgress {
					topicStartedLessons++
				}
				if status == models.LessonStatusLocked {
					topicLockedLessons++
				}

				lessonResponses = append(lessonResponses, LessonResponse{
					ID:               lesson.ID,
					Title:            lesson.Title,
					Slug:             lesson.Slug,
					Status:           status,
					RequiredScore:    lesson.RequiredScore,
					EstimatedMinutes: lesson.EstimatedMinutes,
					XPReward:         lesson.XPReward,
					Score:            progress.Score,
					BestScore:        progress.BestScore,
				})
			}

			topicStatus := aggregateTopicStatus(len(topicLessons), topicCompletedLessons, topicStartedLessons, topicLockedLessons)
			if topicStatus == models.LessonStatusCompleted {
				bandCompletedTopics++
				topicsCompleted++
			}
			totalTopics++

			topicResponses = append(topicResponses, TopicResponse{
				ID:                 topic.ID,
				Title:              topic.Title,
				Slug:               topic.Slug,
				Emoji:              topic.Emoji,
				Color:              topic.Color,
				Status:             topicStatus,
				LessonsCompleted:   topicCompletedLessons,
				TotalLessons:       len(topicLessons),
				ProgressPercentage: percentage(topicCompletedLessons, len(topicLessons)),
				Lessons:            lessonResponses,
			})
		}

		// Band progress uses completed topics over total topics because the roadmap
		// screenshots group progress at topic-card level.
		bandResponses = append(bandResponses, BandLevelResponse{
			ID:                 band.ID,
			BandScore:          band.BandScore,
			Title:              band.Title,
			Description:        band.Description,
			Status:             aggregateBandStatus(len(bandTopics), bandCompletedTopics, topicResponses),
			ProgressPercentage: percentage(bandCompletedTopics, len(bandTopics)),
			TopicsCompleted:    bandCompletedTopics,
			TotalTopics:        len(bandTopics),
			Topics:             topicResponses,
		})
	}

	return Response{
		Course: CourseResponse{
			ID:           course.ID,
			Title:        course.Title,
			Slug:         course.Slug,
			BandMin:      course.BandMin,
			BandMax:      course.BandMax,
			TotalWords:   course.TotalWords,
			TotalLessons: course.TotalLessons,
			TotalTopics:  course.TotalTopics,
		},
		Summary: SummaryResponse{
			TopicsCompleted: topicsCompleted,
			TotalTopics:     totalTopics,
			CurrentBand:     user.CurrentBand,
			WordsMastered:   wordsMastered,
			CurrentStreak:   user.CurrentStreak,
		},
		BandLevels: bandResponses,
	}
}

func groupTopicsByBand(topics []models.Topic) map[uint][]models.Topic {
	result := make(map[uint][]models.Topic)
	for _, topic := range topics {
		result[topic.BandLevelID] = append(result[topic.BandLevelID], topic)
	}
	return result
}

func groupLessonsByTopic(lessons []models.Lesson) map[uint][]models.Lesson {
	result := make(map[uint][]models.Lesson)
	for _, lesson := range lessons {
		result[lesson.TopicID] = append(result[lesson.TopicID], lesson)
	}
	return result
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

func aggregateTopicStatus(totalLessons int, completedLessons int, startedLessons int, lockedLessons int) models.LessonStatus {
	switch {
	case totalLessons == 0:
		return models.LessonStatusLocked
	case completedLessons == totalLessons:
		return models.LessonStatusCompleted
	case completedLessons > 0 || startedLessons > 0:
		return models.LessonStatusInProgress
	case lockedLessons == totalLessons:
		return models.LessonStatusLocked
	default:
		return models.LessonStatusUnlocked
	}
}

func aggregateBandStatus(totalTopics int, completedTopics int, topics []TopicResponse) models.LessonStatus {
	if totalTopics == 0 {
		return models.LessonStatusLocked
	}
	if completedTopics == totalTopics {
		return models.LessonStatusCompleted
	}

	hasProgress := false
	allLocked := true
	for _, topic := range topics {
		if topic.Status == models.LessonStatusInProgress || topic.Status == models.LessonStatusCompleted {
			hasProgress = true
		}
		if topic.Status != models.LessonStatusLocked {
			allLocked = false
		}
	}

	switch {
	case hasProgress:
		return models.LessonStatusInProgress
	case allLocked:
		return models.LessonStatusLocked
	default:
		return models.LessonStatusUnlocked
	}
}

func percentage(completed int, total int) int {
	if total == 0 {
		return 0
	}
	return (completed * 100) / total
}
