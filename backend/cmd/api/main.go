package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ielts-learning/backend/docs"
	"ielts-learning/backend/internal/config"
	"ielts-learning/backend/internal/database"
	"ielts-learning/backend/internal/middleware"
	authmodule "ielts-learning/backend/internal/modules/auth"
	flashcardmodule "ielts-learning/backend/internal/modules/flashcard"
	lessonmodule "ielts-learning/backend/internal/modules/lesson"
	quizmodule "ielts-learning/backend/internal/modules/quiz"
	roadmapmodule "ielts-learning/backend/internal/modules/roadmap"
	topicmodule "ielts-learning/backend/internal/modules/topic"
	vocabularymodule "ielts-learning/backend/internal/modules/vocabulary"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

// @title IELTS Learning API
// @version v1
// @description Backend API for IELTS Vocabulary Learning Platform
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by a JWT access token.
func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	jwtManager := sharedjwt.NewManager(cfg.JWTSecret, cfg.JWTAccessTTLMinutes)

	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to configure trusted proxies: %v", err)
	}

	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.FrontendURL))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	authmodule.RegisterRoutes(api, db, jwtManager)
	roadmapmodule.RegisterRoutes(api, db, jwtManager)
	topicmodule.RegisterRoutes(api, db, jwtManager)
	lessonmodule.RegisterRoutes(api, db, jwtManager)
	vocabularymodule.RegisterRoutes(api, db, jwtManager)
	flashcardmodule.RegisterRoutes(api, db, jwtManager)
	quizmodule.RegisterRoutes(api, db, jwtManager)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
