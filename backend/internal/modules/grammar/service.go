package grammar

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/grammarprogress"
	"ielts-learning/backend/internal/modules/lessonactivity"
	"ielts-learning/backend/internal/modules/mistakes"
)

var ErrExerciseNotAvailable = errors.New("exercise not available for practice")

type Service struct {
	repository Repository
	aiGrader   *lessonactivity.AIGrader
	now        func() time.Time
}

// NewService builds the service. aiGrader may be nil (no AI fallback grading).
func NewService(repository Repository, aiGrader *lessonactivity.AIGrader) Service {
	return Service{repository: repository, aiGrader: aiGrader, now: time.Now}
}

func (s Service) List(userID uint, cefrLevel string) (ListResponse, error) {
	points, err := s.repository.FindGrammarPoints(cefrLevel)
	if err != nil {
		return ListResponse{}, err
	}
	progress, err := s.repository.FindProgress(userID)
	if err != nil {
		return ListResponse{}, err
	}
	counts, err := s.repository.CountAvailableExercisesByTag(userID)
	if err != nil {
		return ListResponse{}, err
	}

	items := make([]GrammarPointItem, 0, len(points))
	for _, point := range points {
		row, ok := progress[point.ID]
		items = append(items, GrammarPointItem{
			ID:            point.ID,
			Slug:          point.Slug,
			CEFRLevel:     point.CEFRLevel,
			Title:         point.Title,
			TitleVI:       point.TitleVI,
			OrderIndex:    point.OrderIndex,
			ExerciseCount: counts[point.SkillTag()],
			Progress:      toProgress(row, ok),
		})
	}
	return ListResponse{Items: items}, nil
}

func (s Service) Get(userID uint, slug string) (DetailResponse, error) {
	point, err := s.repository.FindGrammarPoint(slug)
	if err != nil {
		return DetailResponse{}, err
	}
	row, ok, err := s.repository.FindProgressFor(userID, point.ID)
	if err != nil {
		return DetailResponse{}, err
	}
	counts, err := s.repository.CountAvailableExercisesByTag(userID)
	if err != nil {
		return DetailResponse{}, err
	}

	return DetailResponse{
		GrammarPoint:  lessonactivity.ToGrammarPointCard(point),
		Progress:      toProgress(row, ok),
		ExerciseCount: counts[point.SkillTag()],
	}, nil
}

// Practice picks up to 10 exercises for a grammar point from lessons the
// learner has started: last answered wrong first, then never answered, then
// the ones answered correctly longest ago.
func (s Service) Practice(userID uint, slug string) (PracticeResponse, error) {
	point, err := s.repository.FindGrammarPoint(slug)
	if err != nil {
		return PracticeResponse{}, err
	}
	exercises, err := s.repository.FindAvailableExercises(userID, point.SkillTag())
	if err != nil {
		return PracticeResponse{}, err
	}
	ids := make([]uint, 0, len(exercises))
	for _, exercise := range exercises {
		ids = append(ids, exercise.ID)
	}
	latest, err := s.repository.LatestAnswers(userID, ids)
	if err != nil {
		return PracticeResponse{}, err
	}

	selected := orderForPractice(exercises, latest)
	if len(selected) > practiceSize {
		selected = selected[:practiceSize]
	}

	result := PracticeResponse{Exercises: make([]PracticeExercise, 0, len(selected))}
	for index, exercise := range selected {
		result.Exercises = append(result.Exercises, PracticeExercise{
			ID:         exercise.ID,
			Type:       exercise.Type,
			Prompt:     exercise.Prompt,
			Payload:    json.RawMessage(exercise.PayloadJSON),
			Points:     exercise.Points,
			OrderIndex: index + 1,
			LastResult: lastResult(latest, exercise.ID),
		})
	}
	return result, nil
}

func (s Service) Answer(ctx context.Context, userID uint, slug string, req PracticeAnswerRequest) (PracticeAnswerResponse, error) {
	point, err := s.repository.FindGrammarPoint(slug)
	if err != nil {
		return PracticeAnswerResponse{}, err
	}
	exercises, err := s.repository.FindAvailableExercises(userID, point.SkillTag())
	if err != nil {
		return PracticeAnswerResponse{}, err
	}
	var exercise models.Exercise
	for _, candidate := range exercises {
		if candidate.ID == req.ExerciseID {
			exercise = candidate
			break
		}
	}
	if exercise.ID == 0 {
		return PracticeAnswerResponse{}, ErrExerciseNotAvailable
	}

	matched, err := lessonactivity.Grade(exercise.Type, exercise.PayloadJSON, exercise.AnswerJSON, req.Answer)
	if err != nil {
		return PracticeAnswerResponse{}, err
	}
	graded := s.aiGrader.Resolve(ctx, userID, exercise, req.Answer, matched)

	now := s.now()
	err = s.repository.DB().Transaction(func(tx *gorm.DB) error {
		txRepository := NewRepository(tx)
		if err := txRepository.SavePracticeAnswer(&models.UserExerciseAnswer{
			UserID:     userID,
			ExerciseID: exercise.ID,
			AnswerJSON: datatypes.JSON(req.Answer),
			IsCorrect:  graded.IsCorrect,
			ScoreRatio: graded.ScoreRatio,
			AnsweredAt: now,
		}); err != nil {
			return err
		}
		if err := grammarprogress.Record(tx, userID, exercise.SkillTag, graded.ScoreRatio, graded.IsCorrect, now); err != nil {
			return err
		}
		if mistake, ok := lessonactivity.ExerciseMistake(userID, exercise, string(models.ActivityGrammar), req.Answer, graded, now); ok {
			return mistakes.Upsert(tx, mistake)
		}
		return nil
	})
	if err != nil {
		return PracticeAnswerResponse{}, err
	}

	row, ok, err := s.repository.FindProgressFor(userID, point.ID)
	if err != nil {
		return PracticeAnswerResponse{}, err
	}

	return PracticeAnswerResponse{
		ExerciseID:    exercise.ID,
		IsCorrect:     graded.IsCorrect,
		ScoreRatio:    graded.ScoreRatio,
		Answer:        req.Answer,
		CorrectAnswer: graded.CorrectAnswer,
		ExplanationVI: exercise.ExplanationVI,
		AnsweredAt:    now,
		GradedBy:      graded.GradedBy,
		AIFeedback:    graded.AIFeedback,
		AIStatus:      graded.AIStatus,
		Progress:      toProgress(row, ok),
	}, nil
}

func orderForPractice(exercises []models.Exercise, latest map[uint]models.UserExerciseAnswer) []models.Exercise {
	rank := func(exercise models.Exercise) int {
		answer, ok := latest[exercise.ID]
		switch {
		case ok && !answer.IsCorrect:
			return 0
		case !ok:
			return 1
		default:
			return 2
		}
	}

	ordered := append([]models.Exercise{}, exercises...)
	sort.SliceStable(ordered, func(i, j int) bool {
		ri, rj := rank(ordered[i]), rank(ordered[j])
		if ri != rj {
			return ri < rj
		}
		ai, aOK := latest[ordered[i].ID]
		bj, bOK := latest[ordered[j].ID]
		if !aOK || !bOK {
			return false
		}
		if ri == 0 {
			return ai.AnsweredAt.After(bj.AnsweredAt) // most recent mistakes first
		}
		return ai.AnsweredAt.Before(bj.AnsweredAt) // correct ones: longest ago first
	})
	return ordered
}

func lastResult(latest map[uint]models.UserExerciseAnswer, exerciseID uint) string {
	answer, ok := latest[exerciseID]
	switch {
	case !ok:
		return "NEW"
	case answer.IsCorrect:
		return "CORRECT"
	default:
		return "WRONG"
	}
}

func toProgress(row models.UserGrammarProgress, ok bool) ProgressResponse {
	if !ok {
		return ProgressResponse{}
	}
	return ProgressResponse{
		Mastery:         row.Mastery,
		CorrectCount:    row.CorrectCount,
		WrongCount:      row.WrongCount,
		Practised:       true,
		LastPracticedAt: row.LastPracticed,
	}
}
