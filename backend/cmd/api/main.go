package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/config"
	"ielts-learning/backend/internal/database"
	"ielts-learning/backend/internal/middleware"
	authmodule "ielts-learning/backend/internal/modules/auth"
	roadmapmodule "ielts-learning/backend/internal/modules/roadmap"
	topicmodule "ielts-learning/backend/internal/modules/topic"
	sharedjwt "ielts-learning/backend/internal/shared/jwt"
)

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

	api := router.Group("/api/v1")
	authmodule.RegisterRoutes(api, db, jwtManager)
	roadmapmodule.RegisterRoutes(api, db, jwtManager)
	topicmodule.RegisterRoutes(api, db, jwtManager)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
