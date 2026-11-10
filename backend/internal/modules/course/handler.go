package course

import (
	"errors"
	"net/http"

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

// List godoc
// @Summary List courses
// @Description Return published courses (IELTS and Foundation English) with the user's completed lessons and active course.
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=ListResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /courses [get]
func (h Handler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.List(userID)
	if err != nil {
		writeCourseError(c, err)
		return
	}

	response.OK(c, result)
}

// SetActiveCourse godoc
// @Summary Set active course
// @Description Choose which course the roadmap shows by default.
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SetActiveCourseRequest true "Active course request"
// @Success 200 {object} response.SuccessResponse{data=SetActiveCourseResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /me/active-course [patch]
func (h Handler) SetActiveCourse(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	var req SetActiveCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"request": "Invalid JSON request body",
		})
		return
	}
	if req.CourseID == 0 {
		response.ValidationError(c, map[string]string{
			"courseId": "courseId is required",
		})
		return
	}

	result, err := h.service.SetActiveCourse(userID, req)
	if err != nil {
		writeCourseError(c, err)
		return
	}

	response.OK(c, result)
}

func writeCourseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrCourseNotFound):
		response.Error(c, http.StatusNotFound, "COURSE_NOT_FOUND", "Course was not found")
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
