package grammar

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager) {
	handler := NewHandler(NewService(NewRepository(db)))

	group := router.Group("/grammar-points", middleware.Auth(jwtManager))
	group.GET("", handler.List)
	group.GET("/:slug", handler.Get)
	group.GET("/:slug/practice", handler.Practice)
	group.POST("/:slug/practice/answers", handler.Answer)
}
