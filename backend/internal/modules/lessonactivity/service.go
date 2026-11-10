package lessonactivity

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/datatypes"

	"ielts-learning/backend/internal/models"
)

var (
	ErrLessonLocked          = errors.New("lesson locked")
	ErrActivityLocked        = errors.New("activity locked")
	ErrExerciseNotInActivity = errors.New("exercise not in activity")
)

type Service struct {
	repository Repository
	aiGrader   *AIGrader
	now        func() time.Time
}

// NewService builds the service. aiGrader may be nil (no AI fallback grading).
func NewService(repository Repository, aiGrader *AIGrader) Service {
	return Service{repository: repository, aiGrader: aiGrader, now: time.Now}
}

// lessonContext holds a lesson with its activities and the user's progress.
type lessonContext struct {
	lesson     models.Lesson
	status     models.LessonStatus
	activities []models.LessonActivity
	progress   map[uint]models.UserActivityProgress
}

func (s Service) ListActivities(userID uint, lessonID uint) (LessonActivitiesResponse, error) {
	ctx, err := s.loadLessonContext(userID, lessonID)
	if err != nil {
		return LessonActivitiesResponse{}, err
	}

	activityIDs := make([]uint, 0, len(ctx.activities))
	for _, activity := range ctx.activities {
		activityIDs = append(activityIDs, activity.ID)
	}
	counts, err := s.repository.CountExercises(activityIDs)
	if err != nil {
		return LessonActivitiesResponse{}, err
	}

	locks := activityLocks(ctx.activities, ctx.progress)
	items := make([]ActivityItem, 0, len(ctx.activities))
	summary := ActivitiesProgress{}
	for _, activity := range ctx.activities {
		progress, hasProgress := ctx.progress[activity.ID]
		item := toActivityItem(activity, progress, hasProgress, counts[activity.ID], locks[activity.ID])
		items = append(items, item)

		passed := hasProgress && hasPassed(progress)
		if activity.IsRequired {
			summary.TotalRequired++
			if passed {
				summary.PassedRequired++
			}
		}
		if summary.NextActivityID == nil && !passed && !item.IsLocked {
			id := activity.ID
			summary.NextActivityID = &id
		}
	}

	return LessonActivitiesResponse{
		Lesson: LessonSummary{
			ID:     ctx.lesson.ID,
			Title:  ctx.lesson.Title,
			Slug:   ctx.lesson.Slug,
			Status: ctx.status,
		},
		Progress:   summary,
		Activities: items,
	}, nil
}

func (s Service) GetActivity(userID uint, activityID uint) (ActivityDetailResponse, error) {
	activity, ctx, err := s.loadAccessibleActivity(userID, activityID)
	if err != nil {
		return ActivityDetailResponse{}, err
	}

	exercises, err := s.repository.FindExercises(activity.ID)
	if err != nil {
		return ActivityDetailResponse{}, err
	}

	progress, hasProgress := ctx.progress[activity.ID]
	attempt := currentAttempt(progress, hasProgress)

	exerciseIDs := make([]uint, 0, len(exercises))
	exerciseResponses := make([]ExerciseResponse, 0, len(exercises))
	exerciseByID := make(map[uint]models.Exercise, len(exercises))
	for _, exercise := range exercises {
		exerciseIDs = append(exerciseIDs, exercise.ID)
		exerciseByID[exercise.ID] = exercise
		exerciseResponses = append(exerciseResponses, ExerciseResponse{
			ID:         exercise.ID,
			Type:       exercise.Type,
			Prompt:     exercise.Prompt,
			Payload:    rawJSON(exercise.PayloadJSON),
			Points:     exercise.Points,
			OrderIndex: exercise.OrderIndex,
		})
	}

	answers, err := s.repository.FindAnswers(userID, exerciseIDs, attempt)
	if err != nil {
		return ActivityDetailResponse{}, err
	}
	answerResults := make([]AnswerResult, 0, len(answers))
	for _, answer := range answers {
		exercise := exerciseByID[answer.ExerciseID]
		answerResults = append(answerResults, toAnswerResult(exercise, answer))
	}

	vocabularies, err := s.vocabularyCards(activity)
	if err != nil {
		return ActivityDetailResponse{}, err
	}
	var grammarPoint *GrammarPointCard
	if activity.GrammarPointID != nil {
		point, found, err := s.repository.FindGrammarPoint(*activity.GrammarPointID)
		if err != nil {
			return ActivityDetailResponse{}, err
		}
		if found {
			card := ToGrammarPointCard(point)
			grammarPoint = &card
		}
	}

	phonemes, err := s.phonemeCards(activity)
	if err != nil {
		return ActivityDetailResponse{}, err
	}
	var passage *PassageCard
	if activity.PassageID != nil {
		found, ok, err := s.repository.FindPassage(*activity.PassageID)
		if err != nil {
			return ActivityDetailResponse{}, err
		}
		if ok {
			card := ToPassageCard(found)
			passage = &card
		}
	}

	counts := map[uint]int{activity.ID: len(exercises)}
	locks := activityLocks(ctx.activities, ctx.progress)
	return ActivityDetailResponse{
		Activity: ActivityDetail{
			ActivityItem:   toActivityItem(activity, progress, hasProgress, counts[activity.ID], locks[activity.ID]),
			LessonID:       activity.LessonID,
			Content:        rawJSON(activity.ContentJSON),
			GrammarPointID: activity.GrammarPointID,
			PassageID:      activity.PassageID,
			ScenarioID:     activity.ScenarioID,
			CurrentAttempt: attempt,
		},
		Exercises:    exerciseResponses,
		Vocabularies: vocabularies,
		GrammarPoint: grammarPoint,
		Phonemes:     phonemes,
		Passage:      passage,
		Answers:      answerResults,
	}, nil
}

func (s Service) SubmitAnswer(ctx context.Context, userID uint, activityID uint, req SubmitAnswerRequest) (AnswerResult, error) {
	activity, _, err := s.loadAccessibleActivity(userID, activityID)
	if err != nil {
		return AnswerResult{}, err
	}

	exercises, err := s.repository.FindExercises(activity.ID)
	if err != nil {
		return AnswerResult{}, err
	}
	var exercise models.Exercise
	for _, candidate := range exercises {
		if candidate.ID == req.ExerciseID {
			exercise = candidate
			break
		}
	}
	if exercise.ID == 0 {
		return AnswerResult{}, ErrExerciseNotInActivity
	}

	matched, err := Grade(exercise.Type, exercise.PayloadJSON, exercise.AnswerJSON, req.Answer)
	if err != nil {
		return AnswerResult{}, err
	}
	// AI runs outside the transaction: it can take seconds.
	result := s.aiGrader.Resolve(ctx, userID, exercise, req.Answer, matched)

	now := s.now()
	answer := models.UserExerciseAnswer{
		UserID:     userID,
		ExerciseID: exercise.ID,
		AnswerJSON: datatypes.JSON(req.Answer),
		IsCorrect:  result.IsCorrect,
		ScoreRatio: result.ScoreRatio,
		AnsweredAt: now,
	}

	err = s.repository.Transaction(func(tx Repository) error {
		progress, found, err := tx.FindActivityProgressForUpdate(userID, activity.ID)
		if err != nil {
			return err
		}
		// The first answer after a finished attempt starts a new attempt.
		if !found || progress.Status != models.ActivityInProgress || progress.Attempts == 0 {
			progress.UserID = userID
			progress.ActivityID = activity.ID
			progress.Status = models.ActivityInProgress
			progress.Attempts++
		}
		progress.LastAttemptAt = &now
		if err := tx.SaveActivityProgress(&progress); err != nil {
			return err
		}

		answer.ActivityAttempt = progress.Attempts
		if err := tx.UpsertAnswer(&answer); err != nil {
			return err
		}
		if err := tx.RecordGrammarAnswer(userID, exercise.SkillTag, result.ScoreRatio, result.IsCorrect, now); err != nil {
			return err
		}
		return markLessonStarted(tx, userID, activity.LessonID, now)
	})
	if err != nil {
		return AnswerResult{}, err
	}

	return AnswerResult{
		ExerciseID:    exercise.ID,
		Attempt:       answer.ActivityAttempt,
		IsCorrect:     result.IsCorrect,
		ScoreRatio:    result.ScoreRatio,
		Answer:        req.Answer,
		CorrectAnswer: result.CorrectAnswer,
		ExplanationVI: exercise.ExplanationVI,
		AnsweredAt:    now,
		GradedBy:      result.GradedBy,
		AIFeedback:    result.AIFeedback,
		AIStatus:      result.AIStatus,
	}, nil
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
	lessonProgress, err := s.repository.FindLessonProgress(userID, lessonIDs)
	if err != nil {
		return lessonContext{}, err
	}

	status := lessonStatus(lesson.ID, topicLessons, mapLessonProgress(lessonProgress), lesson.Topic.IsAlwaysUnlocked())
	if status == models.LessonStatusLocked {
		return lessonContext{}, ErrLessonLocked
	}

	activities, err := s.repository.FindLessonActivities(lesson.ID)
	if err != nil {
		return lessonContext{}, err
	}
	activityIDs := make([]uint, 0, len(activities))
	for _, activity := range activities {
		activityIDs = append(activityIDs, activity.ID)
	}
	progressRows, err := s.repository.FindActivityProgress(userID, activityIDs)
	if err != nil {
		return lessonContext{}, err
	}
	progress := make(map[uint]models.UserActivityProgress, len(progressRows))
	for _, row := range progressRows {
		progress[row.ActivityID] = row
	}

	return lessonContext{lesson: lesson, status: status, activities: activities, progress: progress}, nil
}

func (s Service) loadAccessibleActivity(userID uint, activityID uint) (models.LessonActivity, lessonContext, error) {
	activity, err := s.repository.FindActivity(activityID)
	if err != nil {
		return models.LessonActivity{}, lessonContext{}, err
	}

	ctx, err := s.loadLessonContext(userID, activity.LessonID)
	if err != nil {
		return models.LessonActivity{}, lessonContext{}, err
	}

	if activityLocks(ctx.activities, ctx.progress)[activity.ID] {
		return models.LessonActivity{}, lessonContext{}, ErrActivityLocked
	}

	return activity, ctx, nil
}

// activityLocks marks an activity as locked while any earlier required
// activity of the lesson has not been passed.
func activityLocks(activities []models.LessonActivity, progress map[uint]models.UserActivityProgress) map[uint]bool {
	locks := make(map[uint]bool, len(activities))
	blocked := false
	for _, activity := range activities {
		locks[activity.ID] = blocked
		row, ok := progress[activity.ID]
		if activity.IsRequired && !(ok && hasPassed(row)) {
			blocked = true
		}
	}
	return locks
}

// hasPassed is true once an activity has been passed, even while the learner
// is retrying it for a better score — a retry must not re-lock what follows.
func hasPassed(progress models.UserActivityProgress) bool {
	return progress.Status == models.ActivityPassed || progress.PassedAt != nil
}

// currentAttempt is the attempt whose answers should be shown: the open one,
// or 0 when the learner has not started or already finished the activity.
func currentAttempt(progress models.UserActivityProgress, hasProgress bool) int {
	if hasProgress && progress.Status == models.ActivityInProgress {
		return progress.Attempts
	}
	return 0
}

func lessonStatus(
	lessonID uint,
	topicLessons []models.Lesson,
	progressByLesson map[uint]models.UserLessonProgress,
	alwaysUnlocked bool,
) models.LessonStatus {
	if progress, ok := progressByLesson[lessonID]; ok {
		return progress.Status
	}
	if alwaysUnlocked {
		return models.LessonStatusUnlocked
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

func mapLessonProgress(rows []models.UserLessonProgress) map[uint]models.UserLessonProgress {
	result := make(map[uint]models.UserLessonProgress, len(rows))
	for _, row := range rows {
		result[row.LessonID] = row
	}
	return result
}

func toActivityItem(activity models.LessonActivity, progress models.UserActivityProgress, hasProgress bool, exerciseCount int, isLocked bool) ActivityItem {
	item := ActivityItem{
		ID:             activity.ID,
		Slug:           activity.Slug,
		Type:           activity.Type,
		Title:          activity.Title,
		InstructionsVI: activity.InstructionsVI,
		OrderIndex:     activity.OrderIndex,
		IsRequired:     activity.IsRequired,
		PassScore:      activity.PassScore,
		XPReward:       activity.XPReward,
		ExerciseCount:  exerciseCount,
		Status:         ActivityStatusNotStarted,
		IsLocked:       isLocked,
	}
	if hasProgress {
		item.Status = string(progress.Status)
		item.Score = progress.Score
		item.BestScore = progress.BestScore
		item.Attempts = progress.Attempts
	}
	return item
}

func toAnswerResult(exercise models.Exercise, answer models.UserExerciseAnswer) AnswerResult {
	result := AnswerResult{
		ExerciseID:    answer.ExerciseID,
		Attempt:       answer.ActivityAttempt,
		IsCorrect:     answer.IsCorrect,
		ScoreRatio:    answer.ScoreRatio,
		Answer:        rawJSON(answer.AnswerJSON),
		ExplanationVI: exercise.ExplanationVI,
		AnsweredAt:    answer.AnsweredAt,
	}
	// Re-grade only to recover the correct answer for display.
	if graded, err := Grade(exercise.Type, exercise.PayloadJSON, exercise.AnswerJSON, answer.AnswerJSON); err == nil {
		result.CorrectAnswer = graded.CorrectAnswer
	}
	return result
}

func rawJSON(value datatypes.JSON) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage("null")
	}
	return json.RawMessage(value)
}

// markLessonStarted moves an untouched lesson to IN_PROGRESS so the roadmap
// shows it as started.
func markLessonStarted(tx Repository, userID uint, lessonID uint, now time.Time) error {
	progress, found, err := tx.FindLessonProgressForUpdate(userID, lessonID)
	if err != nil {
		return err
	}
	if found && progress.Status != models.LessonStatusUnlocked {
		return nil
	}
	progress.Status = models.LessonStatusInProgress
	if progress.StartedAt == nil {
		progress.StartedAt = &now
	}
	progress.LastStudiedAt = &now
	return tx.SaveLessonProgress(&progress)
}

// phonemeCards returns the sound cards named in content.phonemes, in order.
func (s Service) phonemeCards(activity models.LessonActivity) ([]PhonemeCard, error) {
	var content struct {
		Phonemes []string `json:"phonemes"`
	}
	cards := []PhonemeCard{}
	if len(activity.ContentJSON) == 0 || json.Unmarshal(activity.ContentJSON, &content) != nil || len(content.Phonemes) == 0 {
		return cards, nil
	}
	found, err := s.repository.FindPhonemes(content.Phonemes)
	if err != nil {
		return nil, err
	}
	bySymbol := make(map[string]models.Phoneme, len(found))
	for _, phoneme := range found {
		bySymbol[phoneme.Symbol] = phoneme
	}
	for _, symbol := range content.Phonemes {
		if phoneme, ok := bySymbol[symbol]; ok {
			cards = append(cards, ToPhonemeCard(phoneme))
		}
	}
	return cards, nil
}
