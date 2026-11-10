package contentimport

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/shared/response"
)

const (
	maxAudioBytes = 25 << 20 // Whisper's limit
	maxPageBytes  = 20 << 20
)

// Transcribe godoc
// @Summary Transcribe a recording into a passage draft
// @Description Upload an audio file (multipart "file", max 25MB, mp3/m4a/wav/webm/ogg). It is stored under /media/passages and sent to Whisper; the reply is a passage draft (text + sentence segments with times) to edit and save. Nothing is written to the database.
// @Tags Admin
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Audio file"
// @Param title formData string false "Passage title"
// @Param slug formData string false "Passage slug (default: from the title)"
// @Param genre formData string false "DIALOGUE, MONOLOGUE, ..."
// @Param source formData string false "Book / unit / page"
// @Success 200 {object} response.SuccessResponse{data=PassageDraft}
// @Failure 400 {object} response.ErrorResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /admin/passages/transcribe [post]
func (h Handler) Transcribe(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "AUDIO_INVALID", `An audio file is required under the "file" field`)
		return
	}
	data, ok := readUpload(c, fileHeader, maxAudioBytes, "AUDIO_INVALID")
	if !ok {
		return
	}
	draft, err := h.service.Transcribe(c.Request.Context(), fileHeader.Filename, data,
		c.PostForm("title"), c.PostForm("slug"), c.PostForm("genre"), c.PostForm("source"))
	if err != nil {
		writeMediaError(c, err)
		return
	}
	response.OK(c, draft)
}

// SavePassage godoc
// @Summary Save a passage
// @Description Validate a passage (e.g. an edited transcription draft) and create or update it by slug.
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PassageInput true "Passage"
// @Success 200 {object} response.SuccessResponse{data=lessonactivity.PassageCard}
// @Failure 400 {object} response.ErrorResponse
// @Router /admin/passages [post]
func (h Handler) SavePassage(c *gin.Context) {
	var input PassageInput
	decoder := json.NewDecoder(io.LimitReader(c.Request.Body, maxBundleSizeBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		response.ValidationError(c, map[string]string{"request": "Invalid passage JSON: " + err.Error()})
		return
	}
	card, _, err := h.service.SavePassage(input, input.CEFRLevel)
	if err != nil {
		writeMediaError(c, err)
		return
	}
	response.OK(c, card)
}

// Extract godoc
// @Summary Extract lesson content from pages
// @Description Upload up to 10 page images or one PDF (multipart "files"). A vision model extracts grammar explanations and exercises into a content bundle draft (activities are DRAFT unless publish=true) to review in the import preview. Nothing is written to the database.
// @Tags Admin
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param files formData file true "Page images or one PDF"
// @Param courseSlug formData string false "Default foundation-english"
// @Param cefrLevel formData string false "Default A1"
// @Param topicSlug formData string false "Topic slug for the extracted lesson"
// @Param topicTitle formData string false "Topic title"
// @Param notes formData string false "Extra instructions for the AI"
// @Param publish formData bool false "Publish activities right away"
// @Success 200 {object} response.SuccessResponse{data=ExtractResult}
// @Failure 400 {object} response.ErrorResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /admin/content/extract [post]
func (h Handler) Extract(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil || len(form.File["files"]) == 0 {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", `Upload page images or a PDF under the "files" field`)
		return
	}
	request := ExtractRequest{
		CourseSlug: c.PostForm("courseSlug"),
		CEFRLevel:  c.PostForm("cefrLevel"),
		TopicSlug:  c.PostForm("topicSlug"),
		TopicTitle: c.PostForm("topicTitle"),
		Notes:      c.PostForm("notes"),
		Publish:    c.PostForm("publish") == "true",
	}
	for _, header := range form.File["files"] {
		data, ok := readUpload(c, header, maxPageBytes, "INVALID_INPUT")
		if !ok {
			return
		}
		request.Files = append(request.Files, UploadedFile{Name: header.Filename, Data: data})
	}
	result, err := h.service.Extract(c.Request.Context(), request)
	if err != nil {
		writeMediaError(c, err)
		return
	}
	response.OK(c, result)
}

func readUpload(c *gin.Context, header *multipart.FileHeader, limit int64, code string) ([]byte, bool) {
	if header.Size > limit {
		response.Error(c, http.StatusBadRequest, code, header.Filename+" is too large")
		return nil, false
	}
	file, err := header.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, code, "Unable to read "+header.Filename)
		return nil, false
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil || int64(len(data)) > limit {
		response.Error(c, http.StatusBadRequest, code, "Unable to read "+header.Filename)
		return nil, false
	}
	return data, true
}

func writeMediaError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrAIUnavailable):
		log.Printf("contentimport ai: %v", err)
		message := "AI is not available: set OPENAI_API_KEY, or try again later"
		if detail := strings.TrimPrefix(err.Error(), ErrAIUnavailable.Error()+": "); detail != err.Error() {
			message = "AI request failed: " + detail
		}
		response.Error(c, http.StatusServiceUnavailable, "AI_UNAVAILABLE", message)
	case errors.Is(err, ErrAudioInvalid):
		response.Error(c, http.StatusBadRequest, "AUDIO_INVALID", err.Error())
	case errors.Is(err, ErrPassageInput), errors.Is(err, ErrExtractInput):
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
	default:
		log.Printf("contentimport: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
