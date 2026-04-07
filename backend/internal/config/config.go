package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	GeminiAPIKey   string
	DatabasePath   string
	LogPath        string
	Environment    string
	AllowedOrigins string
	AdminToken     string
	BackdoorPort   string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:           getEnv("PORT", "8080"),
		GeminiAPIKey:   getEnv("GEMINI_API_KEY", ""),
		DatabasePath:   getEnv("DATABASE_PATH", "./data/hotel.db"),
		LogPath:        getEnv("LOG_PATH", "./logs"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		AdminToken:     getEnv("ADMIN_TOKEN", ""),
		BackdoorPort:   getEnv("BACKDOOR_PORT", "9999"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
