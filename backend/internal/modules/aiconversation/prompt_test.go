package aiconversation

import (
	"strings"
	"testing"

	"ielts-learning/backend/internal/models"
)

func TestBuildSystemPromptA1(t *testing.T) {
	level := "A1"
	scenario := models.AIConversationScenario{
		Title: "Greetings", Description: "Morning at the office.", Role: "a colleague called Lan",
		Task: "greet Lan", CEFRLevel: &level, MaxTurns: 6,
		TargetPatternsJSON: []byte(`["Good morning!","How are you?"]`),
		SituationContext:   "Lan is tired today.",
	}
	prompt := buildSystemPrompt(scenario, 3)
	for _, want := range []string{
		"You are a colleague called Lan. Situation: Morning at the office.\n",
		"CEFR A1. Their task: greet Lan.",
		"Good morning! | How are you?",
		"max 10 words per sentence",
		"after 6 learner turns",
		"learner turn 3 of 6",
		"More context from the learner: Lan is tired today.",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt is missing %q:\n%s", want, prompt)
		}
	}
}

func TestCustomScenarioLevel(t *testing.T) {
	if got := scenarioLevel(models.AIConversationScenario{Level: models.DifficultyBeginner}); got != "A2" {
		t.Errorf("beginner custom scenario = %s", got)
	}
	if got := scenarioLevel(models.AIConversationScenario{Level: models.DifficultyAdvanced}); got != "B2" {
		t.Errorf("advanced custom scenario = %s", got)
	}
}

func TestLastMessagesStartsWithLearner(t *testing.T) {
	var history []models.AIConversationMessage
	for i := 0; i < 25; i++ {
		role := chatRoleAssistant
		if i%2 == 1 {
			role = chatRoleUser
		}
		history = append(history, models.AIConversationMessage{Role: role})
	}
	recent := lastMessages(history, 20)
	if len(recent) > 20 || recent[0].Role != chatRoleUser {
		t.Fatalf("got %d messages starting with %s", len(recent), recent[0].Role)
	}
}

func TestFinishReportOverall(t *testing.T) {
	report := finishReport(SessionReport{Scores: SpeakingScores{Grammar: 6, Vocabulary: 7, Fluency: 0, TaskAchievement: 12}})
	if report.Scores.Fluency != 1 || report.Scores.TaskAchievement != 10 || report.Overall != 60 {
		t.Fatalf("report %+v", report)
	}
	if report.Mistakes == nil || report.NewWords == nil || report.GoodPhrases == nil {
		t.Fatal("lists must not be nil")
	}
}
