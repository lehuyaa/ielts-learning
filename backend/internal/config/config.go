package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv              string
	AppPort             string
	FrontendURL         string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	JWTSecret           string
	JWTAccessTTLMinutes string
	OpenAIAPIKey        string
	OpenAIModel         string
	OpenAIBaseURL       string
	// OpenAITranscribeModel is the speech-to-text model for audio import (Task 64).
	OpenAITranscribeModel string
	// MediaDir stores uploaded files (imported recordings), served at /media.
	MediaDir string
	// AIShortAnswerDailyLimit caps AI short-answer gradings per user per day (spec §6.7).
	AIShortAnswerDailyLimit int
}

func Load() Config {
	return Config{
		AppEnv:                  getEnv("APP_ENV", "development"),
		AppPort:                 getEnv("APP_PORT", "8081"),
		FrontendURL:             getEnv("FRONTEND_URL", "http://localhost:5175"),
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  getEnv("DB_PORT", "3308"),
		DBUser:                  getEnv("DB_USER", "root"),
		DBPassword:              getEnv("DB_PASSWORD", "password"),
		DBName:                  getEnv("DB_NAME", "ielts_vocab"),
		JWTSecret:               getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTAccessTTLMinutes:     getEnv("JWT_ACCESS_TTL_MINUTES", "1440"),
		OpenAIAPIKey:            getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:             getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		OpenAIBaseURL:           getEnv("OPENAI_BASE_URL", ""),
		AIShortAnswerDailyLimit: getEnvInt("AI_SHORT_ANSWER_DAILY_LIMIT", 300),
		OpenAITranscribeModel:   getEnv("OPENAI_TRANSCRIBE_MODEL", "whisper-1"),
		MediaDir:                getEnv("MEDIA_DIR", "storage/media"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}
