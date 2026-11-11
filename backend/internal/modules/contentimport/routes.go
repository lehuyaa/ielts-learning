package contentimport

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/shared/ai"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager, aiClient *ai.Client, media MediaStore) {
	service := NewService(NewRepository(db))
	if aiClient != nil {
		service = service.WithAI(aiClient, media)
	}
	handler := NewHandler(service)

	adminGroup := router.Group("/admin/content", middleware.Auth(jwtManager), middleware.RequireAdmin())
	adminGroup.POST("/import/preview", handler.Preview)
	adminGroup.POST("/import", handler.Import)
	adminGroup.POST("/extract", handler.Extract)

	passages := router.Group("/admin/passages", middleware.Auth(jwtManager), middleware.RequireAdmin())
	passages.POST("/transcribe", handler.Transcribe)
	passages.POST("", handler.SavePassage)
}
