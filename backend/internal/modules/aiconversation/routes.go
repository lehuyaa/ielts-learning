package aiconversation

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/shared/ai"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager, client *ai.Client) {
	repository := NewRepository(db)
	service := NewService(repository, client)
	handler := NewHandler(service)

	group := router.Group("/ai-conversations", middleware.Auth(jwtManager))
	group.POST("/scenarios", handler.CreateScenario)
	group.GET("/scenarios", handler.ListScenarios)
	group.GET("/scenarios/:slug/messages", handler.ListMessages)
	group.POST("/scenarios/:slug/messages", handler.SendChatMessage)
}
