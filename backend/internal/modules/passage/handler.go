package passage

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

// Get godoc
// @Summary Get passage
// @Description Return a listening/reading passage: text, audio, sentence segments and glossary.
// @Tags Passages
// @Produce json
// @Security BearerAuth
// @Param passageId path int true "Passage ID"
// @Success 200 {object} response.SuccessResponse{data=lessonactivity.PassageCard}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /passages/{passageId} [get]
func (h Handler) Get(c *gin.Context) {
	id, ok := parseID(c, "passageId")
	if !ok {
		return
	}
	result, err := h.service.Get(id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

// Lookup godoc
// @Summary Look up a tapped word
// @Description Tap-to-translate: meaning, IPA and audio of a word from a text, using the passage glossary first and then the vocabulary (base forms of plurals, -ed, -ing).
// @Tags Passages
// @Produce json
// @Security BearerAuth
// @Param word query string true "Word as written in the text"
// @Param passageId query int false "Passage whose glossary is checked first"
// @Success 200 {object} response.SuccessResponse{data=WordLookupResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /words/lookup [get]
func (h Handler) Lookup(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	word := strings.TrimSpace(c.Query("word"))
	if word == "" || len(word) > 60 {
		response.ValidationError(c, map[string]string{"word": "word is required (max 60 characters)"})
		return
	}
	var passageID uint64
	if raw := c.Query("passageId"); raw != "" {
		var err error
		passageID, err = strconv.ParseUint(raw, 10, 64)
		if err != nil {
			response.ValidationError(c, map[string]string{"passageId": "passageId must be a number"})
			return
		}
	}
	result, err := h.service.Lookup(userID, word, uint(passageID))
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

// AddFlashcard godoc
// @Summary Add a word to flashcards
// @Description Put a vocabulary word into the learner's reviews (status NEW, due now). Adding a word that is already there changes nothing.
// @Tags Passages
// @Produce json
// @Security BearerAuth
// @Param vocabularyId path int true "Vocabulary ID"
// @Success 200 {object} response.SuccessResponse{data=AddFlashcardResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /words/{vocabularyId}/flashcard [post]
func (h Handler) AddFlashcard(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	id, ok := parseID(c, "vocabularyId")
	if !ok {
		return
	}
	result, err := h.service.AddFlashcard(userID, id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

func parseID(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		response.ValidationError(c, map[string]string{name: name + " must be a positive number"})
		return 0, false
	}
	return uint(value), true
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrPassageNotFound):
		response.Error(c, http.StatusNotFound, "PASSAGE_NOT_FOUND", "Passage was not found")
	case errors.Is(err, ErrVocabularyNotFound):
		response.Error(c, http.StatusNotFound, "VOCABULARY_NOT_FOUND", "Word was not found")
	default:
		log.Printf("passage: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
