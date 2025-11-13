package config

import (
	"os"
)

type Config struct {
	Port            string
	OpenSearchURL   string
	MinIOEndpoint   string
	MinIOAccessKey  string
	MinIOSecretKey  string
	JWTSecret       string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		OpenSearchURL:   getEnv("OPENSEARCH_URL", "http://localhost:9200"),
		MinIOEndpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin"),
		JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
