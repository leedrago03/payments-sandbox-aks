package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    ServerPort     string
    SuccessRate    string
    TimeoutRate    string
}

func Load() (*Config, error) {
    godotenv.Load()
    
    return &Config{
        ServerPort:  getEnv("PORT", "3004"),
        SuccessRate: getEnv("SUCCESS_RATE", "80"),
        TimeoutRate: getEnv("TIMEOUT_RATE", "5"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
