package aiconversation

import (
	"fmt"
	"strings"

	"ielts-learning/backend/internal/models"
)

type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string {
	return "validation failed"
}

var allowedLevels = map[models.DifficultyLevel]bool{
	models.DifficultyBeginner:     true,
	models.DifficultyIntermediate: true,
	models.DifficultyAdvanced:     true,
}

func NormalizeCreateScenarioRequest(req CreateScenarioRequest) CreateScenarioRequest {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.SituationContext = strings.TrimSpace(req.SituationContext)
	req.Duration = strings.TrimSpace(req.Duration)

	return req
}

func ValidateCreateScenarioRequest(req CreateScenarioRequest) (CreateScenarioRequest, error) {
	req = NormalizeCreateScenarioRequest(req)
	fields := make(map[string]string)

	if req.Title == "" {
		fields["title"] = "Title is required"
	} else if len(req.Title) < 3 {
		fields["title"] = "Title must be at least 3 characters"
	} else if len(req.Title) > 120 {
		fields["title"] = "Title must be at most 120 characters"
	}

	if req.Description == "" {
		fields["description"] = "Description is required"
	} else if len(req.Description) > 500 {
		fields["description"] = "Description must be at most 500 characters"
	}

	if len(req.SituationContext) > 2000 {
		fields["situationContext"] = "Situation context must be at most 2000 characters"
	}

	if len(req.Duration) > 50 {
		fields["duration"] = "Duration must be at most 50 characters"
	}

	if req.Level == "" {
		req.Level = models.DifficultyBeginner
	} else if !allowedLevels[req.Level] {
		fields["level"] = "Level must be BEGINNER, INTERMEDIATE, or ADVANCED"
	}

	if len(fields) > 0 {
		return req, ValidationError{Fields: fields}
	}

	return req, nil
}

const (
	maxChatMessageLength      = 4000
	maxChatSystemPromptLength = 4000
	maxChatHistoryItems       = 20
)

var allowedChatRoles = map[string]bool{
	"user":      true,
	"assistant": true,
}

func NormalizeChatRequest(req ChatRequest) ChatRequest {
	req.Message = strings.TrimSpace(req.Message)
	req.SystemPrompt = strings.TrimSpace(req.SystemPrompt)
	for i, item := range req.History {
		item.Role = strings.ToLower(strings.TrimSpace(item.Role))
		item.Content = strings.TrimSpace(item.Content)
		req.History[i] = item
	}

	return req
}

func ValidateChatRequest(req ChatRequest) (ChatRequest, error) {
	req = NormalizeChatRequest(req)
	fields := make(map[string]string)

	if req.Message == "" {
		fields["message"] = "Message is required"
	} else if len(req.Message) > maxChatMessageLength {
		fields["message"] = fmt.Sprintf("Message must be at most %d characters", maxChatMessageLength)
	}

	if len(req.SystemPrompt) > maxChatSystemPromptLength {
		fields["systemPrompt"] = fmt.Sprintf("System prompt must be at most %d characters", maxChatSystemPromptLength)
	}

	if len(req.History) > maxChatHistoryItems {
		fields["history"] = fmt.Sprintf("History can include at most %d prior messages", maxChatHistoryItems)
	} else {
		for _, item := range req.History {
			if !allowedChatRoles[item.Role] {
				fields["history"] = "Each history item's role must be \"user\" or \"assistant\""
				break
			}
			if item.Content == "" || len(item.Content) > maxChatMessageLength {
				fields["history"] = fmt.Sprintf("Each history item's content must be 1-%d characters", maxChatMessageLength)
				break
			}
		}
	}

	if len(fields) > 0 {
		return req, ValidationError{Fields: fields}
	}

	return req, nil
}
