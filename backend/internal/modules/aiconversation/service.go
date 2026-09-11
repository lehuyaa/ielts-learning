package aiconversation

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"ielts-learning/backend/internal/models"
)

var ErrSlugGenerationFailed = errors.New("could not generate a unique scenario slug")

var slugInvalidChars = regexp.MustCompile(`[^a-z0-9]+`)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
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
