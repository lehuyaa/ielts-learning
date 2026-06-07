package vocabulary

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

// List godoc
// @Summary List vocabularies
// @Description Return paginated vocabulary items with optional search, filters, and authenticated user progress.
// @Tags Vocabularies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search text"
// @Param difficulty query string false "Difficulty" Enums(BEGINNER, INTERMEDIATE, ADVANCED)
// @Param targetBand query number false "Target IELTS band"
// @Param status query string false "User vocabulary status" Enums(NEW, LEARNING, REVIEW, MASTERED)
// @Param page query int false "Page number"
// @Param limit query int false "Page size, max 100"
// @Success 200 {object} response.SuccessResponse{data=ListResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /vocabularies [get]
func (h Handler) List(c *gin.Context) {
	query, fields := parseListQuery(c)
	if len(fields) > 0 {
		response.ValidationError(c, fields)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.List(userID, query)
	if err != nil {
		writeVocabularyError(c, err)
		return
	}

	response.OK(c, result)
}

// Get godoc
// @Summary Get vocabulary detail
// @Description Return full dictionary-style vocabulary detail and authenticated user progress.
// @Tags Vocabularies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param vocabularyId path int true "Vocabulary ID"
// @Success 200 {object} response.SuccessResponse{data=DetailResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /vocabularies/{vocabularyId} [get]
func (h Handler) Get(c *gin.Context) {
	vocabularyID, err := parseVocabularyID(c.Param("vocabularyId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid vocabulary ID")
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.Get(userID, vocabularyID)
	if err != nil {
		writeVocabularyError(c, err)
		return
	}

	response.OK(c, result)
}

func parseListQuery(c *gin.Context) (ListQuery, map[string]string) {
	fields := make(map[string]string)
	query := ListQuery{
		Q:     strings.TrimSpace(c.Query("q")),
		Page:  1,
		Limit: 20,
	}

	if difficulty := strings.TrimSpace(c.Query("difficulty")); difficulty != "" {
		parsed, ok := parseDifficulty(difficulty)
		if !ok {
			fields["difficulty"] = "Difficulty must be BEGINNER, INTERMEDIATE, or ADVANCED"
		} else {
			query.Difficulty = &parsed
		}
	}

	if targetBand := strings.TrimSpace(c.Query("targetBand")); targetBand != "" {
		parsed, err := strconv.ParseFloat(targetBand, 64)
		if err != nil {
			fields["targetBand"] = "Target band must be a number"
		} else {
			query.TargetBand = &parsed
		}
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		parsed, ok := parseStatus(status)
		if !ok {
			fields["status"] = "Status must be NEW, LEARNING, REVIEW, or MASTERED"
		} else {
			query.Status = &parsed
		}
	}

	if page := strings.TrimSpace(c.Query("page")); page != "" {
		parsed, err := strconv.Atoi(page)
		if err != nil || parsed < 1 {
			fields["page"] = "Page must be at least 1"
		} else {
			query.Page = parsed
		}
	}

	if limit := strings.TrimSpace(c.Query("limit")); limit != "" {
		parsed, err := strconv.Atoi(limit)
		if err != nil || parsed < 1 {
			fields["limit"] = "Limit must be at least 1"
		} else if parsed > 100 {
			fields["limit"] = "Limit must be 100 or less"
		} else {
			query.Limit = parsed
		}
	}

	return query, fields
}

func parseDifficulty(value string) (models.DifficultyLevel, bool) {
	switch models.DifficultyLevel(strings.ToUpper(value)) {
	case models.DifficultyBeginner:
		return models.DifficultyBeginner, true
	case models.DifficultyIntermediate:
		return models.DifficultyIntermediate, true
	case models.DifficultyAdvanced:
		return models.DifficultyAdvanced, true
	default:
		return "", false
	}
}

func parseStatus(value string) (models.VocabularyStatus, bool) {
	switch models.VocabularyStatus(strings.ToUpper(value)) {
	case models.VocabularyStatusNew:
		return models.VocabularyStatusNew, true
	case models.VocabularyStatusLearning:
		return models.VocabularyStatusLearning, true
	case models.VocabularyStatusReview:
		return models.VocabularyStatusReview, true
	case models.VocabularyStatusMastered:
		return models.VocabularyStatusMastered, true
	default:
		return "", false
	}
}

func parseVocabularyID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid vocabulary id")
	}

	return uint(id), nil
}

func writeVocabularyError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrVocabularyNotFound):
		response.Error(c, http.StatusNotFound, "VOCABULARY_NOT_FOUND", "Vocabulary was not found")
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
