package contentimport

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/shared/response"
)

const maxBundleSizeBytes = 5 << 20 // 5MB

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

// Preview godoc
// @Summary Preview a content bundle import
// @Description Validate a lesson content bundle (JSON, multipart field "file") and report what would be created, updated or skipped. Does not write to the database.
// @Tags Admin
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Content bundle (.json, max 5MB)"
// @Success 200 {object} response.SuccessResponse{data=ImportResult}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/content/import/preview [post]
func (h Handler) Preview(c *gin.Context) {
	h.run(c, false)
}

// Import godoc
// @Summary Import a content bundle
// @Description Validate a lesson content bundle and write every valid topic, lesson, activity and exercise in one transaction. Invalid items are skipped.
// @Tags Admin
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Content bundle (.json, max 5MB)"
// @Success 200 {object} response.SuccessResponse{data=ImportResult}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/content/import [post]
func (h Handler) Import(c *gin.Context) {
	h.run(c, true)
}

func (h Handler) run(c *gin.Context, commit bool) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", `A JSON file is required under the "file" field`)
		return
	}
	if fileHeader.Size > maxBundleSizeBytes {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "File is too large (max 5MB)")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Unable to read uploaded file")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxBundleSizeBytes+1))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Unable to read uploaded file")
		return
	}

	result, err := h.service.Import(data, commit)
	if err != nil {
		if errors.Is(err, ErrInvalidBundle) {
			response.Error(c, http.StatusBadRequest, "INVALID_BUNDLE", err.Error())
			return
		}
		log.Printf("contentimport: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
		return
	}

	response.OK(c, result)
}
