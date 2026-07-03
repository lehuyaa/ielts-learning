package dashboard

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

	group := router.Group("/dashboard", middleware.Auth(jwtManager))
	group.GET("/summary", handler.GetSummary)
}
