package passage

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager) {
	handler := NewHandler(NewService(NewRepository(db)))
	auth := middleware.Auth(jwtManager)

	router.GET("/passages/:passageId", auth, handler.Get)
	router.GET("/words/lookup", auth, handler.Lookup)
	router.POST("/words/:vocabularyId/flashcard", auth, handler.AddFlashcard)
}
