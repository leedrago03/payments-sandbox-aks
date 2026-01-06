package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    ServerPort string
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    HMACKey    string
    DBSSLMode  string
}

func Load() (*Config, error) {
    godotenv.Load()
    
    return &Config{
        ServerPort: getEnv("PORT", "3006"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnvRequired("DB_PASSWORD"),
        DBName:     getEnv("DB_NAME", "audit"),
        HMACKey:    getEnvRequired("AUDIT_HMAC_KEY"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvRequired(key string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    panic("Missing required environment variable: " + key)
}
