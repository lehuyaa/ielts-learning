package aiconversation

import (
	"encoding/json"
	"time"

	"ielts-learning/backend/internal/models"
)

type StartSessionRequest struct {
	ScenarioSlug string `json:"scenarioSlug"`
	// ActivityID is set when the session is started from a lesson's SPEAKING activity.
	ActivityID *uint `json:"activityId"`
	// Fresh abandons an unfinished session and starts a new one.
	Fresh bool `json:"fresh"`
}

type SessionMessageRequest struct {
	Message string `json:"message"`
}

// ScenarioBrief is what the learner sees before and during a role-play.
type ScenarioBrief struct {
	ID             uint     `json:"id"`
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	CEFRLevel      string   `json:"cefrLevel"`
	Role           string   `json:"role"`
	Task           string   `json:"task"`
	TargetPatterns []string `json:"targetPatterns"`
	MaxTurns       int      `json:"maxTurns"`
	OpeningLine    string   `json:"openingLine"`
	IsSystem       bool     `json:"isSystem"`
	PassageID      *uint    `json:"passageId"`
	Group          string   `json:"group"`
}

type SessionMessage struct {
	ID        uint      `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type SpeakingScores struct {
	Grammar         int `json:"grammar"`
	Vocabulary      int `json:"vocabulary"`
	Fluency         int `json:"fluency"`
	TaskAchievement int `json:"taskAchievement"`
}

type SpeakingMistake struct {
	Said          string `json:"said"`
	Correct       string `json:"correct"`
	SkillTag      string `json:"skillTag"`
	ExplanationVI string `json:"explanationVi"`
}

// SessionReport is the feedback report of a finished session (spec §6.2).
type SessionReport struct {
	TaskCompleted bool              `json:"taskCompleted"`
	Scores        SpeakingScores    `json:"scores"`
	Overall       int               `json:"overall"`
	Mistakes      []SpeakingMistake `json:"mistakes"`
	GoodPhrases   []string          `json:"goodPhrases"`
	NewWords      []string          `json:"newWords"`
	TipVI         string            `json:"tipVi"`
	XPAwarded     int               `json:"xpAwarded"`
}

type SessionResponse struct {
	ID            uint             `json:"id"`
	Status        string           `json:"status"`
	ActivityID    *uint            `json:"activityId"`
	TurnCount     int              `json:"turnCount"`
	MaxTurns      int              `json:"maxTurns"`
	HintCount     int              `json:"hintCount"`
	HintsAllowed  bool             `json:"hintsAllowed"`
	Score         *int             `json:"score"`
	TaskCompleted bool             `json:"taskCompleted"`
	Report        *SessionReport   `json:"report"`
	Scenario      ScenarioBrief    `json:"scenario"`
	Messages      []SessionMessage `json:"messages"`
	StartedAt     time.Time        `json:"startedAt"`
	FinishedAt    *time.Time       `json:"finishedAt"`
}

type SessionTurnResponse struct {
	Reply     SessionMessage `json:"reply"`
	TurnCount int            `json:"turnCount"`
	MaxTurns  int            `json:"maxTurns"`
	// ShouldFinish is true once the learner has used every turn.
	ShouldFinish bool `json:"shouldFinish"`
}

type HintResponse struct {
	HintVI    string `json:"hintVi"`
	SayThis   string `json:"sayThis"`
	HintCount int    `json:"hintCount"`
}

func toScenarioBrief(scenario models.AIConversationScenario) ScenarioBrief {
	patterns := targetPatterns(scenario)
	if patterns == nil {
		patterns = []string{}
	}
	return ScenarioBrief{
		ID:             scenario.ID,
		Slug:           scenario.Slug,
		Title:          scenario.Title,
		Description:    scenario.Description,
		CEFRLevel:      scenarioLevel(scenario),
		Role:           scenario.Role,
		Task:           scenario.Task,
		TargetPatterns: patterns,
		MaxTurns:       maxTurns(scenario),
		OpeningLine:    openingLine(scenario),
		IsSystem:       scenario.IsSystem,
		PassageID:      scenario.PassageID,
		Group:          scenario.Group,
	}
}

func openingLine(scenario models.AIConversationScenario) string {
	if scenario.OpeningLine != "" {
		return scenario.OpeningLine
	}
	if scenario.SituationContext != "" {
		return "Let's begin. " + scenario.SituationContext
	}
	return "Hi! Let's practise \"" + scenario.Title + "\". " + scenario.Description
}

func toSessionMessages(messages []models.AIConversationMessage) []SessionMessage {
	items := make([]SessionMessage, 0, len(messages))
	for _, message := range messages {
		items = append(items, SessionMessage{ID: message.ID, Role: message.Role, Content: message.Content, CreatedAt: message.CreatedAt})
	}
	return items
}

func decodeReport(raw []byte) *SessionReport {
	if len(raw) == 0 {
		return nil
	}
	var report SessionReport
	if json.Unmarshal(raw, &report) != nil {
		return nil
	}
	return &report
}
