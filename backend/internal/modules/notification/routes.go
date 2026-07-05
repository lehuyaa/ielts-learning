package notification

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

	group := router.Group("/notifications", middleware.Auth(jwtManager))
	group.GET("", handler.ListNotifications)
	group.PATCH("/:id/read", handler.MarkNotificationAsRead)
}
