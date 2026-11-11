package practice

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/shared/ai"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager, client *ai.Client, dailyLimit int) {
	handler := NewHandler(NewService(db, client, dailyLimit))

	group := router.Group("/practice", middleware.Auth(jwtManager))
	group.POST("/generate", handler.Generate)
	group.GET("/sets/:setId", handler.Get)
	group.POST("/sets/:setId/answers", handler.Answer)
}
