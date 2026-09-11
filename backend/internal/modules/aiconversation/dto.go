package aiconversation

import (
	"time"

	"ielts-learning/backend/internal/models"
)

type CreateScenarioRequest struct {
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	SituationContext string                 `json:"situationContext"`
	Level            models.DifficultyLevel `json:"level"`
	Duration         string                 `json:"duration"`
}

type ListScenariosResponse struct {
	Items []ScenarioResponse `json:"items"`
}

type ScenarioResponse struct {
	ID               uint                   `json:"id"`
	Slug             string                 `json:"slug"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	SituationContext string                 `json:"situationContext"`
	Level            models.DifficultyLevel `json:"level"`
	Duration         string                 `json:"duration"`
	CreatedAt        time.Time              `json:"createdAt"`
}

func toScenarioResponse(scenario models.AIConversationScenario) ScenarioResponse {
	return ScenarioResponse{
		ID:               scenario.ID,
		Slug:             scenario.Slug,
		Title:            scenario.Title,
		Description:      scenario.Description,
		SituationContext: scenario.SituationContext,
		Level:            scenario.Level,
		Duration:         scenario.Duration,
		CreatedAt:        scenario.CreatedAt,
	}
}

func toScenarioResponses(scenarios []models.AIConversationScenario) []ScenarioResponse {
	items := make([]ScenarioResponse, 0, len(scenarios))
	for _, scenario := range scenarios {
		items = append(items, toScenarioResponse(scenario))
	}

	return items
}

type ChatHistoryItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message      string            `json:"message"`
	SystemPrompt string            `json:"systemPrompt"`
	History      []ChatHistoryItem `json:"history"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}
