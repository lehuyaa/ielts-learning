package aiconversation

import (
	"encoding/json"
	"fmt"
	"strings"

	"ielts-learning/backend/internal/models"
)

// maxHistoryMessages is how many earlier messages are sent to the model (spec §6.7).
const maxHistoryMessages = 20

// levelRules are the per-level speaking rules of docs/foundation-spec.md §6.1.
var levelRules = map[string]string{
	"A1": "max 10 words per sentence, only A1 words, one question per turn, speak slowly and clearly.",
	"A2": "max 15 words per sentence, A1-A2 words, one or two questions per turn.",
	"B1": "natural sentences, everyday idioms allowed, ask follow-up \"why/how\" questions.",
	"B2": "natural speed and vocabulary, challenge the learner's opinions politely.",
}

// scenarioLevel is the scenario's CEFR level, derived from its difficulty
// for custom scenarios.
func scenarioLevel(scenario models.AIConversationScenario) string {
	if scenario.CEFRLevel != nil && levelRules[*scenario.CEFRLevel] != "" {
		return *scenario.CEFRLevel
	}
	switch scenario.Level {
	case models.DifficultyIntermediate:
		return "B1"
	case models.DifficultyAdvanced:
		return "B2"
	default:
		return "A2"
	}
}

func targetPatterns(scenario models.AIConversationScenario) []string {
	var patterns []string
	if len(scenario.TargetPatternsJSON) > 0 {
		_ = json.Unmarshal(scenario.TargetPatternsJSON, &patterns)
	}
	return patterns
}

func maxTurns(scenario models.AIConversationScenario) int {
	if scenario.MaxTurns > 0 {
		return scenario.MaxTurns
	}
	return 8
}

// buildSystemPrompt is the role-play system prompt (spec §6.1). turn is the
// learner's next turn number, so the AI can close the conversation on time.
func buildSystemPrompt(scenario models.AIConversationScenario, turn int) string {
	level := scenarioLevel(scenario)
	role := strings.TrimSpace(scenario.Role)
	if role == "" {
		role = "a friendly English conversation partner"
	}
	situation := strings.TrimRight(strings.TrimSpace(scenario.Description), ".")
	if situation == "" {
		situation = scenario.Title
	}
	task := strings.TrimRight(strings.TrimSpace(scenario.Task), ".")
	if task == "" {
		task = "have a natural conversation about: " + scenario.Title
	}
	patterns := targetPatterns(scenario)
	limit := maxTurns(scenario)

	var b strings.Builder
	fmt.Fprintf(&b, "You are %s. Situation: %s.\n", role, situation)
	fmt.Fprintf(&b, "The learner is a Vietnamese adult at CEFR %s. Their task: %s.\n", level, task)
	if len(patterns) > 0 {
		fmt.Fprintf(&b, "Useful patterns for this lesson: %s.\n", strings.Join(patterns, " | "))
	}
	fmt.Fprintf(&b, "\nRules for %s: %s\n", level, levelRules[level])
	b.WriteString(`
Always:
- Stay in role. Never switch to Vietnamese unless asked for a hint.
- If the learner makes a mistake, do not stop the conversation; repeat their idea correctly in your reply (recast).
- If the learner is silent, off-topic or writes Vietnamese, give one short English sentence they can copy.
- Reply with only what you say: no stage directions, no lists, no translations.
`)
	fmt.Fprintf(&b, "- End naturally after %d learner turns or when the task is complete, with a closing line. This is learner turn %d of %d.\n", limit, turn, limit)
	if context := strings.TrimSpace(scenario.SituationContext); context != "" {
		fmt.Fprintf(&b, "\nMore context from the learner: %s\n", context)
	}
	return b.String()
}

// lastMessages keeps the most recent n messages, starting with a user turn
// when possible so the model sees whole exchanges.
func lastMessages(history []models.AIConversationMessage, n int) []models.AIConversationMessage {
	if len(history) <= n {
		return history
	}
	recent := history[len(history)-n:]
	if len(recent) > 0 && recent[0].Role == chatRoleAssistant {
		recent = recent[1:]
	}
	return recent
}
