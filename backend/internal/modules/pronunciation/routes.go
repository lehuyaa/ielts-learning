package pronunciation

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager, scorer Scorer, dailyLimit int) {
	handler := NewHandler(NewService(db, scorer, dailyLimit))

	group := router.Group("/pronunciation", middleware.Auth(jwtManager))
	group.GET("/status", handler.Status)
	group.POST("/assess", handler.Assess)
}
