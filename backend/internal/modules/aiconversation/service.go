package aiconversation

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"ielts-learning/backend/internal/models"
)

var ErrSlugGenerationFailed = errors.New("could not generate a unique scenario slug")
var ErrScenarioNotFound = errors.New("scenario not found")

const (
	chatRoleUser      = "user"
	chatRoleAssistant = "assistant"
	chatRoleSystem    = "system"
)

var slugInvalidChars = regexp.MustCompile(`[^a-z0-9]+`)

type Service struct {
	repository   Repository
	openAIClient *openAIClient
}

func NewService(repository Repository, openAIClient *openAIClient) Service {
	return Service{repository: repository, openAIClient: openAIClient}
}

func (s Service) CreateScenario(userID uint, req CreateScenarioRequest) (ScenarioResponse, error) {
	slug, err := s.generateUniqueSlug(userID, req.Title)
	if err != nil {
		return ScenarioResponse{}, err
	}

	scenario := models.AIConversationScenario{
		UserID:           userID,
		Slug:             slug,
		Title:            req.Title,
		Description:      req.Description,
		SituationContext: req.SituationContext,
		Level:            req.Level,
		Duration:         req.Duration,
	}

	if err := s.repository.Create(&scenario); err != nil {
		return ScenarioResponse{}, err
	}

	return toScenarioResponse(scenario), nil
}

func (s Service) ListScenarios(userID uint) (ListScenariosResponse, error) {
	scenarios, err := s.repository.FindByUser(userID)
	if err != nil {
		return ListScenariosResponse{}, err
	}

	return ListScenariosResponse{Items: toScenarioResponses(scenarios)}, nil
}

// SendChatMessage looks up the caller's scenario, replays its persisted
// history (plus the scenario's SituationContext as the system prompt) to the
// AI, then persists both the user's message and the AI's reply so the
// conversation survives across sessions.
func (s Service) SendChatMessage(ctx context.Context, userID uint, scenarioSlug string, userMessage string) (ChatResponse, error) {
	scenario, err := s.repository.FindByUserAndSlug(userID, scenarioSlug)
	if err != nil {
		return ChatResponse{}, err
	}

	history, err := s.repository.FindMessagesByScenario(scenario.ID)
	if err != nil {
		return ChatResponse{}, err
	}

	messages := make([]openAIMessage, 0, len(history)+2)
	if scenario.SituationContext != "" {
		messages = append(messages, openAIMessage{Role: chatRoleSystem, Content: scenario.SituationContext})
	}
	for _, item := range history {
		messages = append(messages, openAIMessage{Role: item.Role, Content: item.Content})
	}
	messages = append(messages, openAIMessage{Role: chatRoleUser, Content: userMessage})

	reply, err := s.openAIClient.CreateChatCompletion(ctx, messages)
	if err != nil {
		return ChatResponse{}, err
	}
	reply = strings.TrimSpace(reply)

	if err := s.repository.CreateMessage(&models.AIConversationMessage{
		ScenarioID: scenario.ID,
		Role:       chatRoleUser,
		Content:    userMessage,
	}); err != nil {
		return ChatResponse{}, err
	}

	if err := s.repository.CreateMessage(&models.AIConversationMessage{
		ScenarioID: scenario.ID,
		Role:       chatRoleAssistant,
		Content:    reply,
	}); err != nil {
		return ChatResponse{}, err
	}

	return ChatResponse{Reply: reply}, nil
}

func (s Service) ListMessages(userID uint, scenarioSlug string) (ListMessagesResponse, error) {
	scenario, err := s.repository.FindByUserAndSlug(userID, scenarioSlug)
	if err != nil {
		return ListMessagesResponse{}, err
	}

	messages, err := s.repository.FindMessagesByScenario(scenario.ID)
	if err != nil {
		return ListMessagesResponse{}, err
	}

	return ListMessagesResponse{Items: toMessageResponses(messages)}, nil
}

func (s Service) generateUniqueSlug(userID uint, title string) (string, error) {
	base := slugify(title)
	slug := base

	const maxAttempts = 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		exists, err := s.repository.ExistsBySlug(userID, slug)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}

		slug = fmt.Sprintf("%s-%d", base, randomSuffix())
	}

	return "", ErrSlugGenerationFailed
}

func slugify(value string) string {
	lowered := strings.ToLower(strings.TrimSpace(value))
	slug := slugInvalidChars.ReplaceAllString(lowered, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "scenario"
	}

	return slug
}

func randomSuffix() int {
	return rand.Intn(900000) + 100000
}
