package aiconversation

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/shared/ai"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager, client *ai.Client, chatDailyLimit int) {
	repository := NewRepository(db)
	service := NewService(repository, client).WithChatLimit(chatDailyLimit)
	handler := NewHandler(service)

	group := router.Group("/ai-conversations", middleware.Auth(jwtManager))
	group.POST("/scenarios", handler.CreateScenario)
	group.GET("/scenarios", handler.ListScenarios)
	group.GET("/scenarios/:slug/messages", handler.ListMessages)
	group.POST("/scenarios/:slug/messages", handler.SendChatMessage)
	group.POST("/sessions", handler.StartSession)
	group.GET("/sessions/:sessionId", handler.GetSession)
	group.POST("/sessions/:sessionId/messages", handler.SendSessionMessage)
	group.POST("/sessions/:sessionId/hint", handler.Hint)
	group.POST("/sessions/:sessionId/finish", handler.Finish)
}
