package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/config"
	"ielts-learning/backend/internal/middleware"
)

func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to configure trusted proxies: %v", err)
	}

	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.FrontendURL))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
