package lessonactivity

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"ielts-learning/backend/internal/models"
	activitymodule "ielts-learning/backend/internal/modules/activity"
	xpmodule "ielts-learning/backend/internal/modules/xp"
)

// ErrActivityNotStarted means there is no open attempt to finish: the learner
// must answer at least one exercise first.
var ErrActivityNotStarted = errors.New("activity not started")

// ErrSpeakingNotCompleted means no role-play session of the activity is finished yet.
var ErrSpeakingNotCompleted = errors.New("speaking not completed")

// CompleteActivity closes the current attempt: scores it, passes or fails the
// activity, awards XP on the first pass, and completes the lesson once every
// required activity has been passed (docs/foundation-spec.md §4).
func (s Service) CompleteActivity(userID uint, activityID uint) (CompleteActivityResponse, error) {
	activity, ctx, err := s.loadAccessibleActivity(userID, activityID)
	if err != nil {
		return CompleteActivityResponse{}, err
	}

	exercises, err := s.repository.FindExercises(activity.ID)
	if err != nil {
		return CompleteActivityResponse{}, err
	}

	progress, hasProgress := ctx.progress[activity.ID]
	attempt := currentAttempt(progress, hasProgress)
	var tally attemptTally
	speakingDone := true
	if activity.Type == models.ActivitySpeaking && activity.ScenarioID != nil {
		// A role-play is scored by its latest finished session; it passes only
		// when the AI judged the task completed (spec §4).
		session, found, err := s.repository.LatestFinishedSession(userID, activity.ID)
		if err != nil {
			return CompleteActivityResponse{}, err
		}
		if !found || session.Score == nil {
			return CompleteActivityResponse{}, ErrSpeakingNotCompleted
		}
		tally = attemptTally{Score: *session.Score}
		speakingDone = session.TaskDone
	} else if activity.Type == models.ActivityWriting {
		// A WRITING activity is scored by the AI report of its latest
		// submission whose rewrite is done (spec §4).
		score, err := s.writingScore(userID, activity.ID)
		if err != nil {
			return CompleteActivityResponse{}, err
		}
		tally = attemptTally{Score: score}
	} else {
		if len(exercises) > 0 && attempt == 0 {
			return CompleteActivityResponse{}, ErrActivityNotStarted
		}
		exerciseIDs := make([]uint, 0, len(exercises))
		for _, exercise := range exercises {
			exerciseIDs = append(exerciseIDs, exercise.ID)
		}
		answers, err := s.repository.FindAnswers(userID, exerciseIDs, attempt)
		if err != nil {
			return CompleteActivityResponse{}, err
		}
		tally = scoreAttempt(exercises, answers)
	}
	passed := tally.Score >= activity.PassScore && speakingDone

	now := s.now()
	result := CompleteActivityResponse{
		Score:         tally.Score,
		Passed:        passed,
		CorrectCount:  tally.Correct,
		AnsweredCount: tally.Answered,
		TotalCount:    len(exercises),
	}

	err = s.repository.Transaction(func(tx Repository) error {
		row, found, err := tx.FindActivityProgressForUpdate(userID, activity.ID)
		if err != nil {
			return err
		}
		if !found {
			// Activities without exercises (e.g. word cards only) finish directly.
			row = models.UserActivityProgress{UserID: userID, ActivityID: activity.ID, Attempts: 1}
		}

		score := tally.Score
		row.Score = &score
		if row.BestScore == nil || score > *row.BestScore {
			row.BestScore = &score
		}
		row.LastAttemptAt = &now
		// A failed retry never takes away an earlier pass.
		row.Status = models.ActivityFailed
		if row.PassedAt != nil {
			row.Status = models.ActivityPassed
		}
		firstPass := false
		if passed {
			row.Status = models.ActivityPassed
			if row.PassedAt == nil {
				row.PassedAt = &now
				firstPass = true
			}
		}
		if err := tx.SaveActivityProgress(&row); err != nil {
			return err
		}
		ctx.progress[activity.ID] = row

		xpEarned := 0
		if firstPass && activity.XPReward > 0 {
			sourceID := activity.ID
			created, err := tx.AwardXP(xpmodule.AwardInput{
				UserID:           userID,
				SourceType:       xpmodule.EventActivityCompleted,
				SourceID:         &sourceID,
				XP:               activity.XPReward,
				Description:      fmt.Sprintf("Completed %s", activity.Title),
				AwardedAt:        now,
				PreventDuplicate: true,
			})
			if err != nil {
				return fmt.Errorf("award activity xp: %w", err)
			}
			if created {
				result.XPAwarded = activity.XPReward
				xpEarned += activity.XPReward
			}
		}

		lessonResult, lessonXP, err := s.updateLessonProgress(tx, userID, ctx, now)
		if err != nil {
			return err
		}
		result.Lesson = lessonResult
		xpEarned += lessonXP

		if err := tx.RecordDailyActivity(userID, activitymodule.Event{
			LessonsDoneDelta: boolToInt(lessonResult.JustCompleted),
			XPEarnedDelta:    xpEarned,
			OccurredAt:       now,
		}); err != nil {
			return fmt.Errorf("record daily activity: %w", err)
		}

		if lessonResult.JustCompleted {
			if err := tx.CheckAchievements(userID); err != nil {
				return fmt.Errorf("check achievements: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return CompleteActivityResponse{}, err
	}

	locks := activityLocks(ctx.activities, ctx.progress)
	result.Activity = toActivityItem(activity, ctx.progress[activity.ID], true, len(exercises), locks[activity.ID])
	result.NextActivityID = nextActivityID(ctx.activities, ctx.progress, locks)
	return result, nil
}

// updateLessonProgress marks the lesson started, and completed once every
// required activity has been passed. It returns the lesson XP awarded.
func (s Service) updateLessonProgress(tx Repository, userID uint, ctx lessonContext, now time.Time) (LessonResult, int, error) {
	lessonProgress, _, err := tx.FindLessonProgressForUpdate(userID, ctx.lesson.ID)
	if err != nil {
		return LessonResult{}, 0, err
	}

	passedRequired, totalRequired, scoreSum := 0, 0, 0
	for _, activity := range ctx.activities {
		if !activity.IsRequired {
			continue
		}
		totalRequired++
		row, ok := ctx.progress[activity.ID]
		if ok && hasPassed(row) {
			passedRequired++
			if row.BestScore != nil {
				scoreSum += *row.BestScore
			}
		}
	}

	if lessonProgress.StartedAt == nil {
		lessonProgress.StartedAt = &now
	}
	lessonProgress.LastStudiedAt = &now

	result := LessonResult{ID: ctx.lesson.ID, PassedRequired: passedRequired, TotalRequired: totalRequired}
	lessonXP := 0
	allPassed := totalRequired > 0 && passedRequired == totalRequired

	switch {
	case allPassed && lessonProgress.Status != models.LessonStatusCompleted:
		score := int(math.Round(float64(scoreSum) / float64(totalRequired)))
		lessonProgress.Status = models.LessonStatusCompleted
		lessonProgress.Score = &score
		if lessonProgress.BestScore == nil || score > *lessonProgress.BestScore {
			lessonProgress.BestScore = &score
		}
		lessonProgress.CompletedAt = &now
		result.JustCompleted = true

		if ctx.lesson.XPReward > 0 {
			sourceID := ctx.lesson.ID
			created, err := tx.AwardXP(xpmodule.AwardInput{
				UserID:           userID,
				SourceType:       xpmodule.EventLessonCompleted,
				SourceID:         &sourceID,
				XP:               ctx.lesson.XPReward,
				Description:      fmt.Sprintf("Completed lesson %s", ctx.lesson.Title),
				AwardedAt:        now,
				PreventDuplicate: true,
			})
			if err != nil {
				return LessonResult{}, 0, fmt.Errorf("award lesson xp: %w", err)
			}
			if created {
				lessonXP = ctx.lesson.XPReward
				if lessonXP > lessonProgress.BestXP {
					lessonProgress.BestXP = lessonXP
				}
			}
		}
	case lessonProgress.Status != models.LessonStatusCompleted:
		lessonProgress.Status = models.LessonStatusInProgress
	}

	if err := tx.SaveLessonProgress(&lessonProgress); err != nil {
		return LessonResult{}, 0, err
	}
	if result.JustCompleted {
		if err := tx.UnlockNextLesson(userID, ctx.lesson, now); err != nil {
			return LessonResult{}, 0, err
		}
	}

	result.Status = lessonProgress.Status
	result.Score = lessonProgress.Score
	result.XPAwarded = lessonXP
	return result, lessonXP, nil
}

type attemptTally struct {
	Score    int
	Correct  int
	Answered int
}

// scoreAttempt computes the activity score: points-weighted average of each
// exercise's score ratio; unanswered exercises count as 0.
func scoreAttempt(exercises []models.Exercise, answers []models.UserExerciseAnswer) attemptTally {
	if len(exercises) == 0 {
		return attemptTally{Score: 100}
	}

	byExercise := make(map[uint]models.UserExerciseAnswer, len(answers))
	for _, answer := range answers {
		byExercise[answer.ExerciseID] = answer
	}

	tally := attemptTally{}
	totalPoints, earned := 0, 0.0
	for _, exercise := range exercises {
		points := exercise.Points
		if points <= 0 {
			points = 1
		}
		totalPoints += points
		answer, ok := byExercise[exercise.ID]
		if !ok {
			continue
		}
		tally.Answered++
		if answer.IsCorrect {
			tally.Correct++
		}
		earned += float64(points) * answer.ScoreRatio
	}

	tally.Score = int(math.Round(100 * earned / float64(totalPoints)))
	return tally
}

func nextActivityID(activities []models.LessonActivity, progress map[uint]models.UserActivityProgress, locks map[uint]bool) *uint {
	for _, activity := range activities {
		row, ok := progress[activity.ID]
		if (!ok || !hasPassed(row)) && !locks[activity.ID] {
			id := activity.ID
			return &id
		}
	}
	return nil
}

// vocabularyCards loads a VOCABULARY activity's words in content order.
func (s Service) vocabularyCards(activity models.LessonActivity) ([]VocabularyCard, error) {
	if activity.Type != models.ActivityVocabulary || len(activity.ContentJSON) == 0 {
		return []VocabularyCard{}, nil
	}

	var content struct {
		VocabularySlugs []string `json:"vocabularySlugs"`
	}
	if err := json.Unmarshal(activity.ContentJSON, &content); err != nil {
		return nil, fmt.Errorf("%w: activity %d content: %v", ErrInvalidAnswerKey, activity.ID, err)
	}

	vocabularies, err := s.repository.FindVocabulariesBySlugs(content.VocabularySlugs)
	if err != nil {
		return nil, err
	}
	bySlug := make(map[string]models.Vocabulary, len(vocabularies))
	for _, vocabulary := range vocabularies {
		bySlug[vocabulary.Slug] = vocabulary
	}

	cards := make([]VocabularyCard, 0, len(content.VocabularySlugs))
	for _, slug := range content.VocabularySlugs {
		vocabulary, ok := bySlug[slug]
		if !ok {
			continue
		}
		cards = append(cards, VocabularyCard{
			ID:               vocabulary.ID,
			Word:             vocabulary.Word,
			Slug:             vocabulary.Slug,
			IPA:              vocabulary.IPA,
			AudioURL:         vocabulary.AudioURL,
			PartOfSpeech:     vocabulary.PartOfSpeech,
			MeaningVI:        vocabulary.MeaningVI,
			MeaningEN:        vocabulary.MeaningEN,
			ExampleSentence:  vocabulary.ExampleSentence,
			ExampleMeaningVI: vocabulary.ExampleMeaningVI,
		})
	}
	return cards, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
