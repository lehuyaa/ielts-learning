package lessonactivity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"gorm.io/datatypes"

	"ielts-learning/backend/internal/models"
)

var (
	ErrNotWritingActivity   = errors.New("not a writing activity")
	ErrWritingTooShort      = errors.New("writing too short")
	ErrWritingTooLong       = errors.New("writing too long")
	ErrWritingNotFound      = errors.New("writing submission not found")
	ErrWritingNotRewritable = errors.New("writing submission does not need a rewrite")
	ErrWritingNotCompleted  = errors.New("writing not completed")
)

const maxWritingChars = 2000

type SubmitWritingRequest struct {
	ActivityID uint   `json:"activityId"`
	Text       string `json:"text"`
}

type RewriteWritingRequest struct {
	Text string `json:"text"`
}

// WritingSubmissionResponse is a graded text with its report and rewrite.
type WritingSubmissionResponse struct {
	ID           uint          `json:"id"`
	ActivityID   uint          `json:"activityId"`
	Text         string        `json:"text"`
	Status       string        `json:"status"`
	Score        int           `json:"score"`
	Passed       bool          `json:"passed"`
	PassScore    int           `json:"passScore"`
	Report       WritingReport `json:"report"`
	RewriteText  *string       `json:"rewriteText"`
	RewriteCount int           `json:"rewriteCount"`
	// RewriteReport lists the mistakes still found in the last rewrite.
	RewriteReport *WritingReport `json:"rewriteReport"`
	ModelAnswer   string         `json:"modelAnswer"`
	CompletedAt   *time.Time     `json:"completedAt"`
	CreatedAt     time.Time      `json:"createdAt"`
}

var sentenceEnd = regexp.MustCompile(`[.!?]+(\s+|$)`)

// countSentences counts sentences that contain at least one letter.
func countSentences(text string) int {
	count := 0
	for _, part := range sentenceEnd.Split(strings.TrimSpace(text), -1) {
		if strings.IndexFunc(part, unicode.IsLetter) >= 0 {
			count++
		}
	}
	return count
}

func writingContent(activity models.LessonActivity) WritingContent {
	var content WritingContent
	_ = json.Unmarshal(activity.ContentJSON, &content)
	if content.MinSentences <= 0 {
		content.MinSentences = 3
	}
	if content.MaxSentences < content.MinSentences {
		content.MaxSentences = content.MinSentences + 4
	}
	return content
}

// SubmitWriting grades a text for a WRITING activity (spec §6.3), records its
// mistakes and starts the activity attempt.
func (s Service) SubmitWriting(ctx context.Context, userID uint, req SubmitWritingRequest) (WritingSubmissionResponse, error) {
	activity, _, err := s.loadAccessibleActivity(userID, req.ActivityID)
	if err != nil {
		return WritingSubmissionResponse{}, err
	}
	if activity.Type != models.ActivityWriting {
		return WritingSubmissionResponse{}, ErrNotWritingActivity
	}
	content := writingContent(activity)
	text := strings.TrimSpace(req.Text)
	if err := checkLength(text, content); err != nil {
		return WritingSubmissionResponse{}, err
	}

	report, err := s.aiGrader.GradeWriting(ctx, userID, content, s.repository.LevelOfActivity(activity.ID), text)
	if err != nil {
		return WritingSubmissionResponse{}, err
	}

	now := s.now()
	reportJSON, _ := json.Marshal(report)
	submission := models.WritingSubmission{
		UserID:     userID,
		ActivityID: activity.ID,
		Text:       text,
		Status:     models.WritingRewriteRequired,
		Score:      report.Overall,
		ReportJSON: datatypes.JSON(reportJSON),
	}
	if !report.RewriteRequired {
		submission.Status = models.WritingCompleted
		submission.CompletedAt = &now
	}

	err = s.repository.Transaction(func(tx Repository) error {
		if err := tx.SaveWritingSubmission(&submission); err != nil {
			return err
		}
		if err := startAttempt(tx, userID, activity.ID, now); err != nil {
			return err
		}
		if err := recordWritingMistakes(tx, userID, report, now); err != nil {
			return err
		}
		return markLessonStarted(tx, userID, activity.LessonID, now)
	})
	if err != nil {
		return WritingSubmissionResponse{}, err
	}
	return toWritingResponse(submission, activity, content), nil
}

// RewriteWriting checks a corrected text. The submission completes when no
// mistakes remain; a rewrite equal to the AI's corrected text needs no AI call.
func (s Service) RewriteWriting(ctx context.Context, userID uint, submissionID uint, req RewriteWritingRequest) (WritingSubmissionResponse, error) {
	current, found, err := s.repository.FindWritingSubmissionForUpdate(userID, submissionID)
	if err != nil {
		return WritingSubmissionResponse{}, err
	}
	if !found {
		return WritingSubmissionResponse{}, ErrWritingNotFound
	}
	if current.Status != models.WritingRewriteRequired {
		return WritingSubmissionResponse{}, ErrWritingNotRewritable
	}
	activity, _, err := s.loadAccessibleActivity(userID, current.ActivityID)
	if err != nil {
		return WritingSubmissionResponse{}, err
	}
	content := writingContent(activity)
	text := strings.TrimSpace(req.Text)
	if err := checkLength(text, content); err != nil {
		return WritingSubmissionResponse{}, err
	}

	var original WritingReport
	_ = json.Unmarshal(current.ReportJSON, &original)

	remaining := WritingReport{Mistakes: []WritingMistake{}}
	if normalizeText(text) != normalizeText(original.CorrectedText) {
		remaining, err = s.aiGrader.CheckRewrite(ctx, userID, content, s.repository.LevelOfActivity(activity.ID), original, text)
		if err != nil {
			return WritingSubmissionResponse{}, err
		}
	}

	now := s.now()
	err = s.repository.Transaction(func(tx Repository) error {
		submission, found, err := tx.FindWritingSubmissionForUpdate(userID, submissionID)
		if err != nil {
			return err
		}
		if !found || submission.Status != models.WritingRewriteRequired {
			return ErrWritingNotRewritable
		}
		submission.RewriteText = &text
		submission.RewriteCount++
		remainingJSON, _ := json.Marshal(remaining)
		submission.RewriteReportJSON = datatypes.JSON(remainingJSON)
		if len(remaining.Mistakes) == 0 {
			submission.Status = models.WritingCompleted
			submission.CompletedAt = &now
		}
		if err := tx.SaveWritingSubmission(&submission); err != nil {
			return err
		}
		current = submission
		return nil
	})
	if err != nil {
		return WritingSubmissionResponse{}, err
	}
	return toWritingResponse(current, activity, content), nil
}

// latestWriting is shown when the activity is opened again.
func (s Service) latestWriting(userID uint, activity models.LessonActivity) (*WritingSubmissionResponse, error) {
	if activity.Type != models.ActivityWriting {
		return nil, nil
	}
	submission, found, err := s.repository.LatestWritingSubmission(userID, activity.ID, false)
	if err != nil || !found {
		return nil, err
	}
	response := toWritingResponse(submission, activity, writingContent(activity))
	return &response, nil
}

// writingScore is the score a WRITING activity is completed with: the latest
// submission whose rewrite is done.
func (s Service) writingScore(userID uint, activityID uint) (int, error) {
	submission, found, err := s.repository.LatestWritingSubmission(userID, activityID, true)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, ErrWritingNotCompleted
	}
	return submission.Score, nil
}

func checkLength(text string, content WritingContent) error {
	if len(text) > maxWritingChars {
		return fmt.Errorf("%w: max %d characters", ErrWritingTooLong, maxWritingChars)
	}
	if count := countSentences(text); count < content.MinSentences {
		return fmt.Errorf("%w: write at least %d sentences (now %d)", ErrWritingTooShort, content.MinSentences, count)
	}
	return nil
}

// startAttempt opens an attempt for the activity so it can be completed.
func startAttempt(tx Repository, userID uint, activityID uint, now time.Time) error {
	progress, found, err := tx.FindActivityProgressForUpdate(userID, activityID)
	if err != nil {
		return err
	}
	if !found || progress.Status != models.ActivityInProgress || progress.Attempts == 0 {
		progress.UserID = userID
		progress.ActivityID = activityID
		progress.Status = models.ActivityInProgress
		progress.Attempts++
	}
	progress.LastAttemptAt = &now
	return tx.SaveActivityProgress(&progress)
}

// recordWritingMistakes logs mistakes (Task 73) and feeds grammar mastery:
// a mistake tagged grammar:<slug> counts as a wrong answer, a target grammar
// point used correctly as a right one.
func recordWritingMistakes(tx Repository, userID uint, report WritingReport, now time.Time) error {
	wrongGrammar := map[string]bool{}
	for _, mistake := range report.Mistakes {
		tag := mistake.SkillTag
		if tag == "" {
			tag = "writing"
		}
		if err := tx.UpsertMistake(models.LearnerMistake{
			UserID:        userID,
			Source:        models.MistakeFromWriting,
			SkillTag:      truncate(tag, 80),
			WrongText:     truncate(mistake.Wrong, 500),
			CorrectText:   truncate(mistake.Correct, 500),
			ExplanationVI: mistake.ExplanationVI,
			Occurrences:   1,
			Status:        models.MistakeOpen,
			LastSeenAt:    now,
		}); err != nil {
			return err
		}
		if strings.HasPrefix(tag, "grammar:") && !wrongGrammar[tag] {
			wrongGrammar[tag] = true
			if err := tx.RecordGrammarAnswer(userID, tag, 0, false, now); err != nil {
				return err
			}
		}
	}
	for _, slug := range report.UsedTargets.Grammar {
		tag := "grammar:" + strings.TrimPrefix(slug, "grammar:")
		if !wrongGrammar[tag] {
			if err := tx.RecordGrammarAnswer(userID, tag, 1, true, now); err != nil {
				return err
			}
		}
	}
	return nil
}

func toWritingResponse(submission models.WritingSubmission, activity models.LessonActivity, content WritingContent) WritingSubmissionResponse {
	var report WritingReport
	_ = json.Unmarshal(submission.ReportJSON, &report)
	if report.Mistakes == nil {
		report.Mistakes = []WritingMistake{}
	}
	response := WritingSubmissionResponse{
		ID:           submission.ID,
		ActivityID:   submission.ActivityID,
		Text:         submission.Text,
		Status:       string(submission.Status),
		Score:        submission.Score,
		Passed:       submission.Score >= activity.PassScore,
		PassScore:    activity.PassScore,
		Report:       report,
		RewriteText:  submission.RewriteText,
		RewriteCount: submission.RewriteCount,
		ModelAnswer:  content.ModelAnswer,
		CompletedAt:  submission.CompletedAt,
		CreatedAt:    submission.CreatedAt,
	}
	if len(submission.RewriteReportJSON) > 0 {
		var remaining WritingReport
		if json.Unmarshal(submission.RewriteReportJSON, &remaining) == nil {
			if remaining.Mistakes == nil {
				remaining.Mistakes = []WritingMistake{}
			}
			response.RewriteReport = &remaining
		}
	}
	return response
}

func truncate(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}
