package pronunciation

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/shared/response"
)

// maxAudioBytes leaves room above 15 s of 16 kHz mono 16-bit PCM (~480 KB).
const maxAudioBytes = 2 << 20

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

// Status godoc
// @Summary Pronunciation scoring status
// @Description Whether Azure pronunciation scoring is configured, the longest recording and the daily limit.
// @Tags Pronunciation
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=StatusResponse}
// @Router /pronunciation/status [get]
func (h Handler) Status(c *gin.Context) {
	response.OK(c, h.service.Status())
}

// Assess godoc
// @Summary Score a recording
// @Description Scores a 16 kHz mono WAV (max 15 s) against referenceText with Azure pronunciation assessment: accuracy, fluency, completeness and per-word / per-phoneme scores (IPA). With exerciseId (a SPEAK exercise) the exercise sentence is used and passed is returned; submit {"attemptId"} as the exercise answer. Updates the learner's sound progress. The audio is not stored.
// @Tags Pronunciation
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param audio formData file true "WAV recording (16 kHz, mono, 16-bit)"
// @Param referenceText formData string false "Text the learner read (required without exerciseId)"
// @Param exerciseId formData int false "SPEAK exercise ID"
// @Param source formData string false "WORD | SHADOWING | READ_ALOUD | EXERCISE"
// @Success 200 {object} response.SuccessResponse{data=AssessResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /pronunciation/assess [post]
func (h Handler) Assess(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	fileHeader, err := c.FormFile("audio")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "AUDIO_INVALID", `A recording is required under the "audio" field`)
		return
	}
	if fileHeader.Size > maxAudioBytes {
		response.Error(c, http.StatusBadRequest, "AUDIO_INVALID", "Recording is too long (max 15 seconds)")
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "AUDIO_INVALID", "Could not read the recording")
		return
	}
	defer file.Close()
	audio, err := io.ReadAll(io.LimitReader(file, maxAudioBytes+1))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "AUDIO_INVALID", "Could not read the recording")
		return
	}

	input := AssessInput{Audio: audio, ReferenceText: c.PostForm("referenceText"), Source: c.PostForm("source")}
	if raw := strings.TrimSpace(c.PostForm("exerciseId")); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || id == 0 {
			response.ValidationError(c, map[string]string{"exerciseId": "exerciseId must be a positive number"})
			return
		}
		exerciseID := uint(id)
		input.ExerciseID = &exerciseID
	}

	result, err := h.service.Assess(c.Request.Context(), userID, input)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrUnavailable):
		response.Error(c, http.StatusServiceUnavailable, "PRONUNCIATION_UNAVAILABLE", "Pronunciation scoring is not set up (AZURE_SPEECH_KEY / AZURE_SPEECH_REGION)")
	case errors.Is(err, ErrInvalidAudio):
		response.Error(c, http.StatusBadRequest, "AUDIO_INVALID", "Recording must be a 16 kHz mono WAV of at most 15 seconds")
	case errors.Is(err, ErrInvalidReference):
		response.ValidationError(c, map[string]string{"referenceText": "referenceText is required (max 500 characters)"})
	case errors.Is(err, ErrExerciseNotFound):
		response.Error(c, http.StatusNotFound, "EXERCISE_NOT_FOUND", "Speaking exercise not found")
	case errors.Is(err, ErrLimitReached):
		response.Error(c, http.StatusTooManyRequests, "PRONUNCIATION_LIMIT_REACHED", "You have used today's pronunciation scoring. Try again tomorrow.")
	case errors.Is(err, ErrNoSpeech):
		response.Error(c, http.StatusUnprocessableEntity, "NO_SPEECH", "No speech was heard. Speak closer to the microphone and try again.")
	case errors.Is(err, ErrScoringFailed):
		response.Error(c, http.StatusServiceUnavailable, "PRONUNCIATION_FAILED", "Pronunciation scoring failed. Please try again.")
	default:
		log.Printf("pronunciation: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
