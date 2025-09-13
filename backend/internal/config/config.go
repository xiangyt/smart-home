package config

import (
	"os"
)

type Config struct {
	Environment        string
	Port               string
	DatabaseURL        string
	HomeAssistantURL   string
	HomeAssistantToken string
	LogLevel           string
}

func Load() *Config {
	return &Config{
		Environment:        getEnv("ENVIRONMENT", "development"),
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "./smart_home.db"),
		HomeAssistantURL:   getEnv("HOME_ASSISTANT_URL", "http://localhost:8123"),
		HomeAssistantToken: getEnv("HOME_ASSISTANT_TOKEN", ""),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
