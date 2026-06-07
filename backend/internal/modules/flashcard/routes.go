package flashcard

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

	lessonGroup := router.Group("/lessons", middleware.Auth(jwtManager))
	lessonGroup.GET("/:lessonId/flashcards", handler.GetLessonFlashcards)

	reviewGroup := router.Group("/reviews", middleware.Auth(jwtManager))
	reviewGroup.GET("/due", handler.GetDueReviews)

	flashcardGroup := router.Group("/flashcards", middleware.Auth(jwtManager))
	flashcardGroup.POST("/review", handler.Review)
}
