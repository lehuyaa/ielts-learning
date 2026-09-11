package aiconversation

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

	group := router.Group("/ai-conversations", middleware.Auth(jwtManager))
	group.POST("/scenarios", handler.CreateScenario)
	group.GET("/scenarios", handler.ListScenarios)
}
