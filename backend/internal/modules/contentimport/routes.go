package contentimport

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager) {
	handler := NewHandler(NewService(NewRepository(db)))

	adminGroup := router.Group("/admin/content", middleware.Auth(jwtManager), middleware.RequireAdmin())
	adminGroup.POST("/import/preview", handler.Preview)
	adminGroup.POST("/import", handler.Import)
}
