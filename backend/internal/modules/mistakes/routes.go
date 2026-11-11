package mistakes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager) {
	handler := NewHandler(NewService(db))

	group := router.Group("/mistakes", middleware.Auth(jwtManager))
	group.GET("", handler.List)
	group.GET("/review", handler.ReviewQueue)
	group.PATCH("/:mistakeId", handler.Update)
	group.POST("/:mistakeId/review", handler.Review)
}
