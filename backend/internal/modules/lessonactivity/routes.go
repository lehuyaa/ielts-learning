package lessonactivity

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager) {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	auth := middleware.Auth(jwtManager)
	router.GET("/lessons/:lessonId/activities", auth, handler.ListActivities)
	router.GET("/activities/:activityId", auth, handler.GetActivity)
	router.POST("/activities/:activityId/answers", auth, handler.SubmitAnswer)
}
