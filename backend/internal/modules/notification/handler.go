package notification

import (
	"errors"
	"net/http"
	"strconv"

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

// ListNotifications godoc
// @Summary List notifications
// @Description Return the authenticated user's notifications ordered by newest first.
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} response.SuccessResponse{data=ListNotificationsResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /notifications [get]
func (h Handler) ListNotifications(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.service.ListNotifications(userID, page, limit)
	if err != nil {
		writeNotificationError(c, err)
		return
	}

	response.OK(c, result)
}

// MarkNotificationAsRead godoc
// @Summary Mark notification as read
// @Description Mark one of the authenticated user's notifications as read.
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Notification ID"
// @Success 200 {object} response.SuccessResponse{data=MarkAsReadResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /notifications/{id}/read [patch]
func (h Handler) MarkNotificationAsRead(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	notificationIDValue, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Notification was not found")
		return
	}

	result, err := h.service.MarkAsRead(userID, uint(notificationIDValue))
	if err != nil {
		writeNotificationError(c, err)
		return
	}

	response.OK(c, result)
}

func writeNotificationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotificationNotFound):
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Notification was not found")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
