package config

import (
	"os"
)

type Config struct {
	DBHost         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	SMSAPIKey      string
	SMSAPISecret   string
	EmailHost      string
	EmailPort      string
	EmailUser      string
	EmailPassword  string
	RedisURL       string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "authdb"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		SMSAPIKey:     getEnv("SMS_API_KEY", ""),
		SMSAPISecret:  getEnv("SMS_API_SECRET", ""),
		EmailHost:     getEnv("EMAIL_HOST", "smtp.gmail.com"),
		EmailPort:     getEnv("EMAIL_PORT", "587"),
		EmailUser:     getEnv("EMAIL_USER", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
