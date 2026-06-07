package quiz

import (
	"errors"
	"math"
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

func (s Service) Get(userID uint, lessonID uint) (GetResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return GetResponse{}, err
	}

	context, err := s.loadContext(userID, lessonID)
	if err != nil {
		return GetResponse{}, err
	}
	if context.status == models.LessonStatusLocked {
		return GetResponse{}, ErrLessonLocked
	}

	return toGetResponse(context.lesson, context.questions), nil
}

func (s Service) Submit(userID uint, lessonID uint, request SubmitRequest) (SubmitResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return SubmitResponse{}, err
	}
	if len(request.Answers) == 0 {
		return SubmitResponse{}, ErrInvalidQuizAnswers
	}

	context, err := s.loadContext(userID, lessonID)
	if err != nil {
		return SubmitResponse{}, err
	}
	if context.status == models.LessonStatusLocked {
		return SubmitResponse{}, ErrLessonLocked
	}

	graded, err := grade(context.lesson, context.questions, request)
	if err != nil {
		return SubmitResponse{}, err
	}

	user, err := s.repository.FindUser(userID)
	if err != nil {
		return SubmitResponse{}, err
	}

	return s.repository.SubmitQuiz(user, context.lesson, context.questions, s.now().UTC(), graded)
}

func (s Service) CheckAnswer(userID uint, lessonID uint, request CheckAnswerRequest) (CheckAnswerResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return CheckAnswerResponse{}, err
	}
	if request.QuestionID == 0 || request.OptionID == 0 {
		return CheckAnswerResponse{}, ErrInvalidQuizAnswers
	}

	context, err := s.loadContext(userID, lessonID)
	if err != nil {
		return CheckAnswerResponse{}, err
	}
	if context.status == models.LessonStatusLocked {
		return CheckAnswerResponse{}, ErrLessonLocked
	}

	question, ok := findQuestion(context.questions, request.QuestionID)
	if !ok {
		return CheckAnswerResponse{}, ErrInvalidQuizAnswers
	}
	if !optionBelongsToQuestion(question, request.OptionID) {
		return CheckAnswerResponse{}, ErrInvalidQuizAnswers
	}

	correct, err := correctOption(question)
	if err != nil {
		return CheckAnswerResponse{}, err
	}

	isCorrect := request.OptionID == correct.ID
	earnedPoints := 0
	if isCorrect {
		earnedPoints = question.Points
	}

	return CheckAnswerResponse{
		QuestionID:       question.ID,
		SelectedOptionID: request.OptionID,
		CorrectOptionID:  correct.ID,
		IsCorrect:        isCorrect,
		Explanation:      question.Explanation,
		EarnedPoints:     earnedPoints,
	}, nil
}

type quizContext struct {
	lesson    models.Lesson
	questions []models.QuizQuestion
	status    models.LessonStatus
}

func (s Service) loadContext(userID uint, lessonID uint) (quizContext, error) {
	lesson, err := s.repository.FindLesson(lessonID)
	if err != nil {
		return quizContext{}, err
	}

	topicLessons, err := s.repository.FindTopicLessons(lesson.TopicID)
	if err != nil {
		return quizContext{}, err
	}
	lessonIDs := make([]uint, 0, len(topicLessons))
	for _, topicLesson := range topicLessons {
		lessonIDs = append(lessonIDs, topicLesson.ID)
	}

	progressRows, err := s.repository.FindLessonProgress(userID, lessonIDs)
	if err != nil {
		return quizContext{}, err
	}

	questions, err := s.repository.FindQuizQuestions(lesson.ID)
	if err != nil {
		return quizContext{}, err
	}

	return quizContext{
		lesson:    lesson,
		questions: questions,
		status:    lessonStatus(lesson.ID, topicLessons, mapLessonProgressByID(progressRows)),
	}, nil
}

type GradedQuiz struct {
	Score        int
	EarnedPoints int
	TotalPoints  int
	CorrectCount int
	Passed       bool
	Results      []SubmitResultResponse
}

func grade(lesson models.Lesson, questions []models.QuizQuestion, request SubmitRequest) (GradedQuiz, error) {
	selectedByQuestion := make(map[uint]uint, len(request.Answers))
	for _, answer := range request.Answers {
		if answer.QuestionID == 0 || answer.OptionID == 0 {
			return GradedQuiz{}, ErrInvalidQuizAnswers
		}
		if _, exists := selectedByQuestion[answer.QuestionID]; exists {
			return GradedQuiz{}, ErrInvalidQuizAnswers
		}
		selectedByQuestion[answer.QuestionID] = answer.OptionID
	}

	totalPoints := 0
	earnedPoints := 0
	correctCount := 0
	results := make([]SubmitResultResponse, 0, len(questions))

	for _, question := range questions {
		totalPoints += question.Points
		correctOption, err := correctOption(question)
		if err != nil {
			return GradedQuiz{}, err
		}

		selectedOptionID, answered := selectedByQuestion[question.ID]
		var selectedOptionIDPtr *uint
		if answered {
			if !optionBelongsToQuestion(question, selectedOptionID) {
				return GradedQuiz{}, ErrInvalidQuizAnswers
			}
			selectedOptionIDPtr = &selectedOptionID
		}

		isCorrect := answered && selectedOptionID == correctOption.ID
		pointsAwarded := 0
		if isCorrect {
			pointsAwarded = question.Points
			earnedPoints += question.Points
			correctCount++
		}

		results = append(results, SubmitResultResponse{
			QuestionID:       question.ID,
			SelectedOptionID: selectedOptionIDPtr,
			CorrectOptionID:  correctOption.ID,
			IsCorrect:        isCorrect,
			Explanation:      question.Explanation,
			PointsAwarded:    pointsAwarded,
		})
	}

	for questionID := range selectedByQuestion {
		if !questionExists(questions, questionID) {
			return GradedQuiz{}, ErrInvalidQuizAnswers
		}
	}

	score := 0
	if totalPoints > 0 {
		score = int(math.Round(float64(earnedPoints) / float64(totalPoints) * 100))
	} else if len(questions) > 0 {
		score = int(math.Round(float64(correctCount) / float64(len(questions)) * 100))
	}

	return GradedQuiz{
		Score:        clampScore(score),
		EarnedPoints: earnedPoints,
		TotalPoints:  totalPoints,
		CorrectCount: correctCount,
		Passed:       clampScore(score) >= lesson.RequiredScore,
		Results:      results,
	}, nil
}

func toGetResponse(lesson models.Lesson, questions []models.QuizQuestion) GetResponse {
	responses := make([]QuestionResponse, 0, len(questions))
	for _, question := range questions {
		options := make([]OptionResponse, 0, len(question.Options))
		for _, option := range question.Options {
			options = append(options, OptionResponse{
				ID:    option.ID,
				Label: option.Label,
				Text:  option.Content,
			})
		}

		responses = append(responses, QuestionResponse{
			ID:      question.ID,
			Type:    question.Type,
			Prompt:  question.Question,
			Points:  question.Points,
			Options: options,
		})
	}

	return GetResponse{
		Lesson: LessonResponse{
			ID:            lesson.ID,
			Title:         lesson.Title,
			RequiredScore: lesson.RequiredScore,
			TimeLimit:     lesson.QuizTimeLimitSeconds,
		},
		Questions: responses,
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

func mapLessonProgressByID(rows []models.UserLessonProgress) map[uint]models.UserLessonProgress {
	result := make(map[uint]models.UserLessonProgress, len(rows))
	for _, row := range rows {
		result[row.LessonID] = row
	}
	return result
}

func correctOption(question models.QuizQuestion) (models.QuizOption, error) {
	for _, option := range question.Options {
		if option.IsCorrect {
			return option, nil
		}
	}

	return models.QuizOption{}, ErrInvalidQuizAnswers
}

func findQuestion(questions []models.QuizQuestion, questionID uint) (models.QuizQuestion, bool) {
	for _, question := range questions {
		if question.ID == questionID {
			return question, true
		}
	}

	return models.QuizQuestion{}, false
}

func optionBelongsToQuestion(question models.QuizQuestion, optionID uint) bool {
	for _, option := range question.Options {
		if option.ID == optionID {
			return true
		}
	}
	return false
}

func questionExists(questions []models.QuizQuestion, questionID uint) bool {
	for _, question := range questions {
		if question.ID == questionID {
			return true
		}
	}
	return false
}

func clampScore(score int) int {
	switch {
	case score < 0:
		return 0
	case score > 100:
		return 100
	default:
		return score
	}
}
