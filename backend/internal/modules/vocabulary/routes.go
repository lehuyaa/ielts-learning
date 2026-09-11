package vocabulary

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

	group := router.Group("/vocabularies", middleware.Auth(jwtManager))
	group.GET("", handler.List)
	group.GET("/:vocabularyId", handler.Get)

	adminGroup := router.Group("/admin/vocabularies", middleware.Auth(jwtManager), middleware.RequireAdmin())
	adminGroup.POST("/import/preview", handler.ImportPreview)
	adminGroup.POST("/import", handler.Import)
}
