package grammar

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/modules/lessonactivity"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, jwtManager sharedjwt.Manager, aiGrader *lessonactivity.AIGrader) {
	handler := NewHandler(NewService(NewRepository(db), aiGrader))

	group := router.Group("/grammar-points", middleware.Auth(jwtManager))
	group.GET("", handler.List)
	group.GET("/:slug", handler.Get)
	group.GET("/:slug/practice", handler.Practice)
	group.POST("/:slug/practice/answers", handler.Answer)
}
