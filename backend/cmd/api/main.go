package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ielts-learning/backend/docs"
	"ielts-learning/backend/internal/config"
	"ielts-learning/backend/internal/database"
	"ielts-learning/backend/internal/middleware"
	aiconversationmodule "ielts-learning/backend/internal/modules/aiconversation"
	authmodule "ielts-learning/backend/internal/modules/auth"
	contentimportmodule "ielts-learning/backend/internal/modules/contentimport"
	coursemodule "ielts-learning/backend/internal/modules/course"
	dashboardmodule "ielts-learning/backend/internal/modules/dashboard"
	flashcardmodule "ielts-learning/backend/internal/modules/flashcard"
	foundationmodule "ielts-learning/backend/internal/modules/foundation"
	grammarmodule "ielts-learning/backend/internal/modules/grammar"
	lessonmodule "ielts-learning/backend/internal/modules/lesson"
	lessonactivitymodule "ielts-learning/backend/internal/modules/lessonactivity"
	mistakesmodule "ielts-learning/backend/internal/modules/mistakes"
	notificationmodule "ielts-learning/backend/internal/modules/notification"
	passagemodule "ielts-learning/backend/internal/modules/passage"
	phonememodule "ielts-learning/backend/internal/modules/phoneme"
	practicemodule "ielts-learning/backend/internal/modules/practice"
	profilemodule "ielts-learning/backend/internal/modules/profile"
	pronunciationmodule "ielts-learning/backend/internal/modules/pronunciation"
	quizmodule "ielts-learning/backend/internal/modules/quiz"
	roadmapmodule "ielts-learning/backend/internal/modules/roadmap"
	topicmodule "ielts-learning/backend/internal/modules/topic"
	vocabularymodule "ielts-learning/backend/internal/modules/vocabulary"
	sharedai "ielts-learning/backend/internal/shared/ai"
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
	// .env is optional (e.g. in production env vars are set directly), so a
	// missing file is not a fatal error.
	_ = godotenv.Load()

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
	if err := os.MkdirAll(cfg.MediaDir, 0o755); err != nil {
		log.Fatalf("failed to create media folder: %v", err)
	}
	router.Static("/media", cfg.MediaDir)

	api := router.Group("/api/v1")
	authmodule.RegisterRoutes(api, db, jwtManager)
	dashboardmodule.RegisterRoutes(api, db, jwtManager)
	roadmapmodule.RegisterRoutes(api, db, jwtManager)
	coursemodule.RegisterRoutes(api, db, jwtManager)
	topicmodule.RegisterRoutes(api, db, jwtManager)
	lessonmodule.RegisterRoutes(api, db, jwtManager)
	aiClient := sharedai.New(cfg.OpenAIAPIKey, cfg.OpenAIModel, cfg.OpenAIBaseURL).WithTranscribeModel(cfg.OpenAITranscribeModel)
	// Vision extraction of many pages takes longer than a chat turn.
	contentAIClient := sharedai.New(cfg.OpenAIAPIKey, cfg.OpenAIModel, cfg.OpenAIBaseURL).
		WithTranscribeModel(cfg.OpenAITranscribeModel).WithTimeout(3 * time.Minute)
	aiGrader := lessonactivitymodule.NewAIGrader(db, aiClient, cfg.OpenAIModel, cfg.AIShortAnswerDailyLimit).
		WithWritingLimit(cfg.AIWritingDailyLimit)
	lessonactivitymodule.RegisterRoutes(api, db, jwtManager, aiGrader)
	grammarmodule.RegisterRoutes(api, db, jwtManager, aiGrader)
	phonememodule.RegisterRoutes(api, db, jwtManager)
	pronunciationmodule.RegisterRoutes(api, db, jwtManager,
		pronunciationmodule.NewAzureClient(cfg.AzureSpeechKey, cfg.AzureSpeechRegion, cfg.AzureSpeechEndpoint),
		cfg.PronunciationDailyLimit)
	passagemodule.RegisterRoutes(api, db, jwtManager)
	profilemodule.RegisterRoutes(api, db, jwtManager)
	notificationmodule.RegisterRoutes(api, db, jwtManager)
	vocabularymodule.RegisterRoutes(api, db, jwtManager)
	contentimportmodule.RegisterRoutes(api, db, jwtManager, contentAIClient, contentimportmodule.MediaStore{Dir: cfg.MediaDir})
	flashcardmodule.RegisterRoutes(api, db, jwtManager)
	quizmodule.RegisterRoutes(api, db, jwtManager)
	mistakesmodule.RegisterRoutes(api, db, jwtManager)
	foundationmodule.RegisterRoutes(api, db, jwtManager)
	practicemodule.RegisterRoutes(api, db, jwtManager, aiClient, cfg.AIPracticeDailyLimit)
	aiconversationmodule.RegisterRoutes(api, db, jwtManager, aiClient, cfg.AIChatDailyLimit)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
