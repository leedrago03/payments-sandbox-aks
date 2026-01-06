package config

import (
    "fmt"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    ServerPort   string
    DBHost       string
    DBPort       string
    DBUser       string
    DBPassword   string
    DBName       string
    Environment  string
    DBSSLMode    string
}

func Load() (*Config, error) {
    godotenv.Load()
    
    return &Config{
        ServerPort:  getEnv("PORT", "3002"),
        DBHost:      getEnv("DB_HOST", "localhost"),
        DBPort:      getEnv("DB_PORT", "5432"),
        DBUser:      getEnv("DB_USER", "postgres"),
        DBPassword:  getEnvRequired("DB_PASSWORD"),
        DBName:      getEnv("DB_NAME", "merchants"),
        Environment: getEnv("ENVIRONMENT", "development"),
        DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
    }, nil
}

func (c *Config) GetDBConnectionString() string {
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
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
