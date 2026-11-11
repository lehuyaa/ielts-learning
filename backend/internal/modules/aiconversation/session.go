package aiconversation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
	activitymodule "ielts-learning/backend/internal/modules/activity"
	"ielts-learning/backend/internal/modules/aiusage"
	"ielts-learning/backend/internal/modules/grammarprogress"
	"ielts-learning/backend/internal/modules/mistakes"
	xpmodule "ielts-learning/backend/internal/modules/xp"
	"ielts-learning/backend/internal/shared/ai"
)

var (
	ErrSessionFinished = errors.New("session finished")
	ErrTurnLimit       = errors.New("turn limit reached")
	ErrHintsNotAllowed = errors.New("hints not allowed")
	ErrSessionEmpty    = errors.New("session has no learner turns")
	ErrAILimitReached  = errors.New("ai limit reached")
	ErrAIUnavailable   = errors.New("ai unavailable")
)

// sessionXP is awarded once per finished session whose task was completed.
const sessionXP = 10

const usageChat = "CHAT"

// StartSession starts a role-play of a scenario (optionally for a lesson's
// SPEAKING activity), or resumes the unfinished one.
func (s Service) StartSession(userID uint, req StartSessionRequest) (SessionResponse, error) {
	scenario, err := s.repository.FindScenarioForUser(userID, strings.TrimSpace(req.ScenarioSlug))
	if err != nil {
		return SessionResponse{}, err
	}
	var activity *models.LessonActivity
	if req.ActivityID != nil {
		found, err := s.repository.FindActivity(*req.ActivityID)
		if err != nil {
			return SessionResponse{}, err
		}
		if found.ScenarioID == nil || *found.ScenarioID != scenario.ID {
			return SessionResponse{}, ErrActivityMismatch
		}
		activity = &found
	}

	open, found, err := s.repository.FindOpenSession(userID, scenario.ID, req.ActivityID)
	if err != nil {
		return SessionResponse{}, err
	}
	if found && !req.Fresh {
		return s.sessionResponse(open, scenario, activity)
	}

	now := s.now()
	session := models.AIConversationSession{
		UserID:     userID,
		ScenarioID: scenario.ID,
		ActivityID: req.ActivityID,
		Status:     models.SessionInProgress,
		StartedAt:  now,
	}
	err = s.repository.DB().Transaction(func(tx *gorm.DB) error {
		repo := NewRepository(tx)
		if found {
			open.Status = models.SessionAbandoned
			if err := repo.SaveSession(&open); err != nil {
				return err
			}
		}
		if err := repo.CreateSession(&session); err != nil {
			return err
		}
		return repo.CreateMessages([]models.AIConversationMessage{{
			ScenarioID: scenario.ID, SessionID: &session.ID, Role: chatRoleAssistant, Content: openingLine(scenario),
		}})
	})
	if err != nil {
		return SessionResponse{}, err
	}
	return s.sessionResponse(session, scenario, activity)
}

func (s Service) GetSession(userID uint, sessionID uint) (SessionResponse, error) {
	session, err := s.repository.FindSession(userID, sessionID, false)
	if err != nil {
		return SessionResponse{}, err
	}
	scenario, err := s.repository.FindScenarioByID(session.ScenarioID)
	if err != nil {
		return SessionResponse{}, err
	}
	activity, err := s.sessionActivity(session)
	if err != nil {
		return SessionResponse{}, err
	}
	return s.sessionResponse(session, scenario, activity)
}

// SendSessionMessage sends one learner turn and returns the AI's reply,
// built with the level-aware prompt and the last 20 messages.
func (s Service) SendSessionMessage(ctx context.Context, userID uint, sessionID uint, text string) (SessionTurnResponse, error) {
	session, scenario, history, err := s.openSession(userID, sessionID)
	if err != nil {
		return SessionTurnResponse{}, err
	}
	limit := maxTurns(scenario)
	if session.TurnCount >= limit {
		return SessionTurnResponse{}, ErrTurnLimit
	}
	if err := s.takeQuota(userID); err != nil {
		return SessionTurnResponse{}, err
	}

	messages := []openAIMessage{{Role: chatRoleSystem, Content: buildSystemPrompt(scenario, session.TurnCount+1)}}
	for _, item := range lastMessages(history, maxHistoryMessages) {
		messages = append(messages, openAIMessage{Role: item.Role, Content: item.Content})
	}
	messages = append(messages, openAIMessage{Role: chatRoleUser, Content: text})
	reply, err := s.openAIClient.CreateChatCompletion(ctx, messages)
	if err != nil {
		return SessionTurnResponse{}, aiError(err)
	}
	reply = strings.TrimSpace(reply)

	saved := []models.AIConversationMessage{
		{ScenarioID: scenario.ID, SessionID: &session.ID, Role: chatRoleUser, Content: text},
		{ScenarioID: scenario.ID, SessionID: &session.ID, Role: chatRoleAssistant, Content: reply},
	}
	err = s.repository.DB().Transaction(func(tx *gorm.DB) error {
		repo := NewRepository(tx)
		current, err := repo.FindSession(userID, sessionID, true)
		if err != nil {
			return err
		}
		if current.Status != models.SessionInProgress {
			return ErrSessionFinished
		}
		if err := repo.CreateMessages(saved); err != nil {
			return err
		}
		current.TurnCount++
		session = current
		return repo.SaveSession(&current)
	})
	if err != nil {
		return SessionTurnResponse{}, err
	}

	answer := saved[1]
	return SessionTurnResponse{
		Reply:        SessionMessage{ID: answer.ID, Role: answer.Role, Content: answer.Content, CreatedAt: answer.CreatedAt},
		TurnCount:    session.TurnCount,
		MaxTurns:     limit,
		ShouldFinish: session.TurnCount >= limit,
	}, nil
}

var hintSchema = ai.Schema{
	Name: "speaking_hint",
	Schema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"hintVi":  map[string]any{"type": "string", "description": "one short Vietnamese sentence: what to say next"},
			"sayThis": map[string]any{"type": "string", "description": "one short English sentence the learner can say now"},
		},
		"required":             []string{"hintVi", "sayThis"},
		"additionalProperties": false,
	},
}

// Hint suggests what to say next (spec §6.5); not allowed in checkpoints.
func (s Service) Hint(ctx context.Context, userID uint, sessionID uint) (HintResponse, error) {
	session, scenario, history, err := s.openSession(userID, sessionID)
	if err != nil {
		return HintResponse{}, err
	}
	activity, err := s.sessionActivity(session)
	if err != nil {
		return HintResponse{}, err
	}
	if !hintsAllowed(activity) {
		return HintResponse{}, ErrHintsNotAllowed
	}
	if err := s.takeQuota(userID); err != nil {
		return HintResponse{}, err
	}

	level := scenarioLevel(scenario)
	system := fmt.Sprintf("You help a Vietnamese adult learner (CEFR %s) who is stuck in an English role-play. "+
		"Their task: %s. Using the conversation so far, give one short Vietnamese hint for what to say next and one simple English sentence at their level that they can say now.",
		level, defaultText(scenario.Task, scenario.Title))
	transcript := transcriptOf(lastMessages(history, 6))
	var hint HintResponse
	if err := s.openAIClient.CompleteJSON(ctx, []ai.Message{{Role: "system", Content: system}, {Role: "user", Content: transcript}}, hintSchema, &hint); err != nil {
		return HintResponse{}, aiError(err)
	}

	err = s.repository.DB().Model(&models.AIConversationSession{}).Where("id = ?", session.ID).
		Update("hint_count", gorm.Expr("hint_count + 1")).Error
	if err != nil {
		return HintResponse{}, fmt.Errorf("count hint: %w", err)
	}
	hint.HintCount = session.HintCount + 1
	return hint, nil
}

func scoreField() map[string]any { return map[string]any{"type": "integer", "description": "1-10"} }

var reportSchema = ai.Schema{
	Name: "speaking_report",
	Schema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"taskCompleted": map[string]any{"type": "boolean"},
			"scores": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"grammar": scoreField(), "vocabulary": scoreField(), "fluency": scoreField(), "taskAchievement": scoreField(),
				},
				"required":             []string{"grammar", "vocabulary", "fluency", "taskAchievement"},
				"additionalProperties": false,
			},
			"mistakes": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"said": map[string]any{"type": "string"}, "correct": map[string]any{"type": "string"},
						"skillTag": map[string]any{"type": "string"}, "explanationVi": map[string]any{"type": "string"},
					},
					"required":             []string{"said", "correct", "skillTag", "explanationVi"},
					"additionalProperties": false,
				},
			},
			"goodPhrases": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"newWords":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"tipVi":       map[string]any{"type": "string"},
		},
		"required":             []string{"taskCompleted", "scores", "mistakes", "goodPhrases", "newWords", "tipVi"},
		"additionalProperties": false,
	},
}

// Finish closes a session with the AI feedback report (spec §6.2): scores,
// mistakes (logged and fed to grammar mastery), good phrases, new words; XP
// when the task was completed.
func (s Service) Finish(ctx context.Context, userID uint, sessionID uint) (SessionResponse, error) {
	session, scenario, history, err := s.openSession(userID, sessionID)
	if err != nil {
		return SessionResponse{}, err
	}
	if session.TurnCount == 0 {
		return SessionResponse{}, ErrSessionEmpty
	}
	if err := s.takeQuota(userID); err != nil {
		return SessionResponse{}, err
	}

	level := scenarioLevel(scenario)
	system := fmt.Sprintf(`You assess a role-play by a Vietnamese adult learning English at CEFR %s.
Role of the AI: %s. Learner's task: %s. Target patterns: %s.
Judge only the learner's lines (marked "Learner").
- taskCompleted: did the learner do the task?
- scores 1-10 for grammar, vocabulary, fluency (complete sentences, few one-word answers), taskAchievement; be fair for the level.
- mistakes: real errors only, "said" quoted from the learner, "correct" the fixed sentence, skillTag "grammar:<slug>" (g1-to-be, g2-nouns-articles, g3-present-simple, g4-present-continuous, g5-there-is-prepositions, g6-adjectives-comparison, g7-can-must-should, g8-past-simple, g9-future) or "vocab:word-choice", "spelling"; explanationVi one short Vietnamese sentence. At most 8.
- goodPhrases: up to 5 good sentences the learner said. newWords: up to 6 useful words from the conversation the learner may not know.
- tipVi: one Vietnamese sentence with the most useful advice.`,
		level, defaultText(scenario.Role, "a conversation partner"), defaultText(scenario.Task, scenario.Title), strings.Join(targetPatterns(scenario), " | "))
	var report SessionReport
	if err := s.openAIClient.CompleteJSON(ctx, []ai.Message{{Role: "system", Content: system}, {Role: "user", Content: transcriptOf(history)}}, reportSchema, &report); err != nil {
		return SessionResponse{}, aiError(err)
	}
	report = finishReport(report)

	now := s.now()
	activity, err := s.sessionActivity(session)
	if err != nil {
		return SessionResponse{}, err
	}
	err = s.repository.DB().Transaction(func(tx *gorm.DB) error {
		current, err := NewRepository(tx).FindSession(userID, sessionID, true)
		if err != nil {
			return err
		}
		if current.Status != models.SessionInProgress {
			return ErrSessionFinished
		}

		xp := 0
		if report.TaskCompleted {
			sourceID := current.ID
			result, err := s.xp.AwardXP(tx, xpmodule.AwardInput{
				UserID: userID, SourceType: xpmodule.EventSpeakingSession, SourceID: &sourceID, XP: sessionXP,
				Description: "Speaking: " + scenario.Title, AwardedAt: now, PreventDuplicate: true,
			})
			if err != nil {
				return fmt.Errorf("award speaking xp: %w", err)
			}
			if result.EventCreated {
				xp = sessionXP
			}
		}
		report.XPAwarded = xp

		reportJSON, _ := json.Marshal(report)
		score := report.Overall
		current.Status = models.SessionFinished
		current.Score = &score
		current.TaskDone = report.TaskCompleted
		current.ReportJSON = datatypes.JSON(reportJSON)
		current.FinishedAt = &now
		if err := NewRepository(tx).SaveSession(&current); err != nil {
			return err
		}
		session = current

		wrongGrammar := map[string]bool{}
		for _, mistake := range report.Mistakes {
			if err := mistakes.Upsert(tx, models.LearnerMistake{
				UserID: userID, Source: models.MistakeFromSpeaking, SkillTag: defaultText(mistake.SkillTag, "speaking"),
				WrongText: mistake.Said, CorrectText: mistake.Correct, ExplanationVI: mistake.ExplanationVI, LastSeenAt: now,
			}); err != nil {
				return err
			}
			if strings.HasPrefix(mistake.SkillTag, "grammar:") && !wrongGrammar[mistake.SkillTag] {
				wrongGrammar[mistake.SkillTag] = true
				if err := grammarprogress.Record(tx, userID, mistake.SkillTag, 0, false, now); err != nil {
					return err
				}
			}
		}

		minutes := int(math.Min(30, math.Ceil(now.Sub(current.StartedAt).Minutes())))
		return s.daily.RecordActivityTx(tx, userID, activitymodule.Event{XPEarnedDelta: xp, ActiveMinutesDelta: minutes, OccurredAt: now})
	})
	if err != nil {
		return SessionResponse{}, err
	}
	return s.sessionResponse(session, scenario, activity)
}

// finishReport clamps scores and computes overall = 10 × average.
func finishReport(report SessionReport) SessionReport {
	clamp := func(value int) int { return max(1, min(10, value)) }
	report.Scores = SpeakingScores{
		Grammar: clamp(report.Scores.Grammar), Vocabulary: clamp(report.Scores.Vocabulary),
		Fluency: clamp(report.Scores.Fluency), TaskAchievement: clamp(report.Scores.TaskAchievement),
	}
	sum := report.Scores.Grammar + report.Scores.Vocabulary + report.Scores.Fluency + report.Scores.TaskAchievement
	report.Overall = int(math.Round(float64(sum) * 10 / 4))
	if report.Mistakes == nil {
		report.Mistakes = []SpeakingMistake{}
	}
	if report.GoodPhrases == nil {
		report.GoodPhrases = []string{}
	}
	if report.NewWords == nil {
		report.NewWords = []string{}
	}
	return report
}

// openSession loads an unfinished session with its scenario and messages.
func (s Service) openSession(userID uint, sessionID uint) (models.AIConversationSession, models.AIConversationScenario, []models.AIConversationMessage, error) {
	session, err := s.repository.FindSession(userID, sessionID, false)
	if err != nil {
		return session, models.AIConversationScenario{}, nil, err
	}
	if session.Status != models.SessionInProgress {
		return session, models.AIConversationScenario{}, nil, ErrSessionFinished
	}
	scenario, err := s.repository.FindScenarioByID(session.ScenarioID)
	if err != nil {
		return session, scenario, nil, err
	}
	history, err := s.repository.FindSessionMessages(session.ID)
	return session, scenario, history, err
}

func (s Service) sessionActivity(session models.AIConversationSession) (*models.LessonActivity, error) {
	if session.ActivityID == nil {
		return nil, nil
	}
	activity, err := s.repository.FindActivity(*session.ActivityID)
	if errors.Is(err, ErrActivityMismatch) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s Service) sessionResponse(session models.AIConversationSession, scenario models.AIConversationScenario, activity *models.LessonActivity) (SessionResponse, error) {
	messages, err := s.repository.FindSessionMessages(session.ID)
	if err != nil {
		return SessionResponse{}, err
	}
	return SessionResponse{
		ID:            session.ID,
		Status:        string(session.Status),
		ActivityID:    session.ActivityID,
		TurnCount:     session.TurnCount,
		MaxTurns:      maxTurns(scenario),
		HintCount:     session.HintCount,
		HintsAllowed:  hintsAllowed(activity),
		Score:         session.Score,
		TaskCompleted: session.TaskDone,
		Report:        decodeReport(session.ReportJSON),
		Scenario:      toScenarioBrief(scenario),
		Messages:      toSessionMessages(messages),
		StartedAt:     session.StartedAt,
		FinishedAt:    session.FinishedAt,
	}, nil
}

func (s Service) takeQuota(userID uint) error {
	if s.openAIClient == nil || !s.openAIClient.Configured() {
		return ErrAIChatNotConfigured
	}
	allowed, err := aiusage.Take(s.repository.DB(), userID, usageChat, s.chatLimit, s.now())
	if err != nil {
		return err
	}
	if !allowed {
		return ErrAILimitReached
	}
	return nil
}

func aiError(err error) error {
	if errors.Is(err, ErrAIChatNotConfigured) {
		return err
	}
	return fmt.Errorf("%w: %v", ErrAIUnavailable, err)
}

// transcriptOf renders messages as "AI:" / "Learner:" lines.
func transcriptOf(messages []models.AIConversationMessage) string {
	var b strings.Builder
	for _, message := range messages {
		speaker := "AI"
		if message.Role == chatRoleUser {
			speaker = "Learner"
		}
		fmt.Fprintf(&b, "%s: %s\n", speaker, message.Content)
	}
	return b.String()
}

func defaultText(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
