package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DatabaseURL string
    FrontendURL string
}

func LoadConfig() Config {

    // Load .env file.
    // Environment variables already set by the OS take precedence.
    if err := godotenv.Load(); err != nil {
        log.Println(".env file not found, using system environment variables")
    }

    cfg := Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
        FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
    }

    if cfg.DatabaseURL == "" {
        log.Fatal("DATABASE_URL is required")
    }

    return cfg
}

func getEnv(key, fallback string) string {

    value := os.Getenv(key)

    if value == "" {
        return fallback
    }

    return value
}