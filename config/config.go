package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	JWTExpireHour int
}

func Load() *Config {
	// Load from environment variables, fallback to defaults for development.
	expire, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOUR", "24"))
	if err != nil {
		expire = 24
	}
	return &Config{
		AppPort:       getEnv("APP_PORT", "3000"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "campuscare"),
		JWTSecret:     getEnv("JWT_SECRET", "change_me_secret"),
		JWTExpireHour: expire,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c *Config) DSN() string {
	// pgx connection string
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func (c *Config) JWTExpireDuration() time.Duration {
	return time.Hour * time.Duration(c.JWTExpireHour)
}
