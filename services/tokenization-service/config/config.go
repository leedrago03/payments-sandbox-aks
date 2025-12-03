package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    ServerPort     string
    DBHost         string
    DBPort         string
    DBUser         string
    DBPassword     string
    DBName         string
    EncryptionKey  string
}

func Load() (*Config, error) {
    godotenv.Load()
    
    return &Config{
        ServerPort:    getEnv("PORT", "3003"),
        DBHost:        getEnv("DB_HOST", "localhost"),
        DBPort:        getEnv("DB_PORT", "5432"),
        DBUser:        getEnv("DB_USER", "postgres"),
        DBPassword:    getEnv("DB_PASSWORD", "postgres123"),
        DBName:        getEnv("DB_NAME", "tokenization"),
        EncryptionKey: getEnv("ENCRYPTION_KEY", "dev-32-byte-encryption-key-12345"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
